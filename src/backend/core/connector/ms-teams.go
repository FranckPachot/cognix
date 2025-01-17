package connector

import (
	microsoftcore "cognix.ch/api/v2/core/connector/microsoft-core"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"jaytaylor.com/html2text"
	"strconv"
	"strings"
	"time"
)

const (
	msTeamsChannelsURL = "https://graph.microsoft.com/v1.0/teams/%s/channels"
	msTeamsMessagesURL = "https://graph.microsoft.com/v1.0/teams/%s/channels/%s/messages/microsoft.graph.delta()"
	//msTeamsMessagesURL = "https://graph.microsoft.com/v1.0/teams/%s/channels/%s/messages"
	msTeamRepliesURL = "https://graph.microsoft.com/v1.0/teams/%s/channels/%s/messages/%s/replies"
	msTeamsInfoURL   = "https://graph.microsoft.com/v1.0/teams"

	msTeamsFilesFolder   = "https://graph.microsoft.com/v1.0/teams/%s/channels/%s/filesFolder"
	msTeamsFolderContent = "https://graph.microsoft.com/v1.0/groups/%s/drive/items/%s/children"

	msTeamsChats           = "https://graph.microsoft.com/v1.0/chats?$top=50"
	msTeamsChatMessagesURL = "https://graph.microsoft.com/v1.0/chats/%s/messages?$top=50"

	msTeamsParamTeamID = "team_id"

	messageTemplate = `#%s
##%s
`

	messageTypeMessage            = "message"
	attachmentContentTypReference = "reference"
)

/*
https://graph.microsoft.com/v1.0/groups/94100e5f-a30f-433d-965e-bde4e817f62a/team/channels/19:65a0a68789ea4abe97c8eec4d6f43786@thread.tacv2
https://graph.microsoft.com/v1.0/teams/94100e5f-a30f-433d-965e-bde4e817f62a/channels
https://graph.microsoft.com/v1.0/teams/94100e5f-a30f-433d-965e-bde4e817f62a/channels/19:65a0a68789ea4abe97c8eec4d6f43786@thread.tacv2/messages
https://graph.microsoft.com/v1.0/teams/94100e5f-a30f-433d-965e-bde4e817f62a/channels/19:65a0a68789ea4abe97c8eec4d6f43786@thread.tacv2/messages/
https://graph.microsoft.com/v1.0/teams/94100e5f-a30f-433d-965e-bde4e817f62a/channels/19:65a0a68789ea4abe97c8eec4d6f43786@thread.tacv2/messages/1718016334912/replies

https://graph.microsoft.com/v1.0/teams/94100e5f-a30f-433d-965e-bde4e817f62a/channels/19:65a0a68789ea4abe97c8eec4d6f43786@thread.tacv2/messages/1718121958378/replies

https://graph.microsoft.com/v1.0/drives/01SZITRJYIBUNPFKHAYJCLISV4DXJPIJV4/items/

https://graph.microsoft.com/v1.0/drives/b!oxsuyS45_EKmyHYegUv4SmEjVp8sBIFPvH1TNMZJZqPviFyz50UFTqjI-nC6wDfJ

https://graph.microsoft.com/v1.0/groups/94100e5f-a30f-433d-965e-bde4e817f62a/drive/root/children

// get drive items for channel
https://graph.microsoft.com/v1.0/teams/94100e5f-a30f-433d-965e-bde4e817f62a/channels/19:65a0a68789ea4abe97c8eec4d6f43786@thread.tacv2/filesFolder
// get files from channel
https://graph.microsoft.com/v1.0/groups/94100e5f-a30f-433d-965e-bde4e817f62a/drive/items/01SZITRJYIBUNPFKHAYJCLISV4DXJPIJV4/children

	/groups/94100e5f-a30f-433d-965e-bde4e817f62a/drive/items/01SZITRJYIBUNPFKHAYJCLISV4DXJPIJV4
*/
type (
	MSTeams struct {
		Base
		param         *MSTeamParameters
		state         *MSTeamState
		client        *resty.Client
		fileSizeLimit int
		sessionID     uuid.NullUUID
	}
	MSTeamParameters struct {
		Team         string                      `json:"team"`
		Channels     model.StringSlice           `json:"channels"`
		AnalyzeChats bool                        `json:"analyze_chats"`
		Token        *oauth2.Token               `json:"token"`
		Files        *microsoftcore.MSDriveParam `json:"files"`
	}
	// MSTeamState store ms team state after each execute
	MSTeamState struct {
		Channels map[string]*MSTeamChannelState `json:"channels"`
		Chats    map[string]*MSTeamMessageState `json:"chats"`
	}

	MSTeamChannelState struct {
		// Link for request changes after last execution
		DeltaLink string                         `json:"delta_link"`
		Topics    map[string]*MSTeamMessageState `json:"topics"`
	}
	// MSTeamMessageState store
	MSTeamMessageState struct {
		LastCreatedDateTime time.Time `json:"last_created_date_time"`
	}
	MSTeamsResult struct {
		PrevLoadTime string
		Messages     []string
	}
)

func (c *MSTeams) Validate() error {
	return nil
}

func (c *MSTeams) PrepareTask(ctx context.Context, sessionID uuid.UUID, task Task) error {
	params := make(map[string]string)

	if c.param.Team != "" {
		teamID, err := c.getTeamID(ctx)
		if err != nil {
			zap.S().Errorf("Prepare task get teamID : %s ", err.Error())
			return err
		}
		params[msTeamsParamTeamID] = teamID
	}
	params[model.ParamSessionID] = sessionID.String()
	return task.RunConnector(ctx, &proto.ConnectorRequest{
		Id:     c.model.ID.IntPart(),
		Params: params,
	})
}

func (c *MSTeams) Execute(ctx context.Context, param map[string]string) chan *Response {
	c.resultCh = make(chan *Response)
	var fileSizeLimit int
	if size, ok := param[model.ParamFileLimit]; ok {
		fileSizeLimit, _ = strconv.Atoi(size)
	}
	if fileSizeLimit == 0 {
		fileSizeLimit = 1
	}
	c.fileSizeLimit = fileSizeLimit * model.GB
	paramSessionID, _ := param[model.ParamSessionID]
	if uuidSessionID, err := uuid.Parse(paramSessionID); err != nil {
		c.sessionID = uuid.NullUUID{uuid.New(), true}
	} else {
		c.sessionID = uuid.NullUUID{uuidSessionID, true}
	}

	for _, doc := range c.model.Docs {
		if doc.Signature == "" {
			// do not delete document with chat history.
			doc.IsExists = true
		}
	}
	go func() {
		defer close(c.resultCh)
		if err := c.execute(ctx, param); err != nil {
			zap.S().Errorf("execute %s ", err.Error())
		}
		return
	}()
	return c.resultCh
}
func (c *MSTeams) execute(ctx context.Context, param map[string]string) error {

	if c.param.AnalyzeChats {
		msDrive := microsoftcore.NewMSDrive(c.param.Files,
			c.model,
			c.sessionID, c.client,
			"", "",
			c.getFile,
		)
		if err := c.loadChats(ctx, msDrive, ""); err != nil {
			zap.S().Errorf("error loading chats : %s ", err.Error())
			//return fmt.Errorf("load chats : %s", err.Error())
		}
	}

	if teamID, ok := param[msTeamsParamTeamID]; ok {
		if err := c.loadChannels(ctx, teamID); err != nil {
			zap.S().Errorf("error loading channels : %s ", err.Error())
			//return fmt.Errorf("load channels : %s", err.Error())
		}
	}
	// save current state
	zap.S().Infof("save connector state.")
	if err := c.model.State.FromStruct(c.state); err == nil {
		return c.connectorRepo.Update(ctx, c.model)
	}
	return nil
}

func (c *MSTeams) loadChannels(ctx context.Context, teamID string) error {
	channelIDs, err := c.getChannel(ctx, teamID)
	if err != nil {
		return err
	}

	// loop by channels
	for _, channelID := range channelIDs {
		// prepare state for channel
		channelState, ok := c.state.Channels[channelID]
		if !ok {
			channelState = &MSTeamChannelState{
				DeltaLink: "",
				Topics:    make(map[string]*MSTeamMessageState),
			}
			c.state.Channels[channelID] = channelState
		}

		topics, err := c.getTopicsByChannel(ctx, teamID, channelID)
		if err != nil {
			return err
		}

		//  load topics
		for _, topic := range topics {
			// create unique id for store new messages in new document
			sourceID := fmt.Sprintf("%s-%s-%s", channelID, topic.Id, uuid.New().String())

			replies, err := c.getReplies(ctx, teamID, channelID, topic)
			if err != nil {
				return err
			}
			if len(replies.Messages) == 0 {
				continue
			}
			doc := &model.Document{
				SourceID:        sourceID,
				ConnectorID:     c.model.ID,
				URL:             "",
				ChunkingSession: c.sessionID,
				Analyzed:        false,
				CreationDate:    time.Now().UTC(),
				LastUpdate:      pg.NullTime{time.Now().UTC()},
				OriginalURL:     topic.WebUrl,
				IsExists:        true,
			}
			c.model.DocsMap[sourceID] = doc

			fileName := fmt.Sprintf("%s_%s.md",
				strings.ReplaceAll(uuid.New().String(), "-", ""),
				strings.ReplaceAll(topic.Subject, " ", ""))
			c.resultCh <- &Response{
				URL:        doc.URL,
				Name:       fileName,
				SourceID:   sourceID,
				DocumentID: doc.ID.IntPart(),
				MimeType:   "plain/text",
				FileType:   proto.FileType_MD,
				Signature:  "",
				Content: &Content{
					Bucket:        model.BucketName(c.model.User.EmbeddingModel.TenantID),
					URL:           "",
					AppendContent: true,
					Body:          []byte(strings.Join(replies.Messages, "\n")),
				},
				UpToData: false,
			}
		}

		if c.param.Files != nil {
			if err = c.loadFiles(ctx, teamID, channelID); err != nil {
				return err
			}
		}
	}
	return nil
}

// loadFiles scrap channel files
func (c *MSTeams) loadFiles(ctx context.Context, teamID, channelID string) error {
	var folderInfo microsoftcore.TeamFilesFolder
	if err := c.requestAndParse(ctx, fmt.Sprintf(msTeamsFilesFolder, teamID, channelID), &folderInfo); err != nil {
		return err
	}
	baseUrl := fmt.Sprintf(msTeamsFolderContent, teamID, folderInfo.Id)
	folderURL := fmt.Sprintf(msTeamsFolderContent, teamID, `%s`)
	msDrive := microsoftcore.NewMSDrive(c.param.Files,
		c.model,
		c.sessionID, c.client,
		baseUrl, folderURL,
		c.getFile,
	)
	return msDrive.Execute(ctx, c.fileSizeLimit)

}

// getChannel get channels from team
func (c *MSTeams) getChannel(ctx context.Context, teamID string) ([]string, error) {
	var channelResp microsoftcore.ChannelResponse
	if err := c.requestAndParse(ctx, fmt.Sprintf(msTeamsChannelsURL, teamID), &channelResp); err != nil {
		return nil, err
	}
	var channels []string
	for _, channel := range channelResp.Value {
		if len(c.param.Channels) == 0 ||
			c.param.Channels.InArray(channel.DisplayName) {
			channels = append(channels, channel.Id)
		}
	}
	if len(channels) == 0 {
		return nil, fmt.Errorf("channel not found")
	}
	return channels, nil
}

func (c *MSTeams) getReplies(ctx context.Context, teamID, channelID string, msg *microsoftcore.MessageBody) (*MSTeamsResult, error) {
	var repliesResp microsoftcore.MessageResponse
	err := c.requestAndParse(ctx, fmt.Sprintf(msTeamRepliesURL, teamID, channelID, msg.Id), &repliesResp)
	if err != nil {
		return nil, err
	}
	var result MSTeamsResult
	var messages []string

	state, ok := c.state.Channels[channelID].Topics[msg.Id]
	if !ok {
		state = &MSTeamMessageState{}
		c.state.Channels[channelID].Topics[msg.Id] = state

		if message := c.buildMDMessage(msg); message != "" {
			messages = append(messages, message)
		}
	} else {
		result.PrevLoadTime = state.LastCreatedDateTime.Format("2006-01-02-15-04-05")
	}
	lastTime := state.LastCreatedDateTime

	for _, repl := range repliesResp.Value {
		if state.LastCreatedDateTime.UTC().After(repl.CreatedDateTime.UTC()) ||
			state.LastCreatedDateTime.UTC().Equal(repl.CreatedDateTime.UTC()) {
			// ignore messages that were analyzed before
			continue
		}
		if repl.CreatedDateTime.UTC().After(lastTime.UTC()) {
			// store timestamp of last message
			lastTime = repl.CreatedDateTime
		}
		if message := c.buildMDMessage(repl); message != "" {
			messages = append(messages, message)
		}

	}
	result.Messages = messages
	state.LastCreatedDateTime = lastTime
	return &result, nil
}

func (c *MSTeams) getTopicsByChannel(ctx context.Context, teamID, channelID string) ([]*microsoftcore.MessageBody, error) {
	var messagesResp microsoftcore.MessageResponse
	// Get url from state. Load changes from previous scan.
	state := c.state.Channels[channelID]

	url := state.DeltaLink
	if url == "" {
		// Load all history if stored lin is empty
		url = fmt.Sprintf(msTeamsMessagesURL, teamID, channelID)
	}

	if err := c.requestAndParse(ctx, url, &messagesResp); err != nil {
		return nil, err
	}
	if len(messagesResp.Value) > 0 {
		if messagesResp.OdataNextLink != "" {
			state.DeltaLink = messagesResp.OdataNextLink
		}
		if messagesResp.OdataDeltaLink != "" {
			state.DeltaLink = messagesResp.OdataDeltaLink
		}
	}
	return messagesResp.Value, nil
}

// getTeamID get team id for current user
func (c *MSTeams) getTeamID(ctx context.Context) (string, error) {
	var team microsoftcore.TeamResponse

	if err := c.requestAndParse(ctx, msTeamsInfoURL, &team); err != nil {
		return "", err
	}
	if len(team.Value) == 0 {
		return "", fmt.Errorf("team not found")
	}
	for _, tm := range team.Value {
		if tm.DisplayName == c.param.Team {
			return tm.Id, nil
		}
	}
	return "", fmt.Errorf("team not found")
}

// requestAndParse request graph endpoint and parse result.
func (c *MSTeams) requestAndParse(ctx context.Context, url string, result interface{}) error {
	response, err := c.client.R().SetContext(ctx).Get(url)
	if err = utils.WrapRestyError(response, err); err != nil {
		return err
	}
	return json.Unmarshal(response.Body(), result)
}

// getFile callback for receive files
func (c *MSTeams) getFile(payload *microsoftcore.Response) {
	response := &Response{
		URL:        payload.URL,
		Name:       payload.Name,
		SourceID:   payload.SourceID,
		DocumentID: payload.DocumentID,
		MimeType:   payload.MimeType,
		FileType:   payload.FileType,
		Signature:  payload.Signature,
		Content: &Content{
			Bucket: model.BucketName(c.model.User.EmbeddingModel.TenantID),
			URL:    payload.URL,
		},
	}
	c.resultCh <- response
}

func (c *MSTeams) buildMDMessage(msg *microsoftcore.MessageBody) string {
	userName := msg.Subject
	if msg.From != nil && msg.From.User != nil {
		userName = msg.From.User.DisplayName
	}
	message := msg.Subject
	if msg.Body != nil {
		message = msg.Body.Content
		if msg.Body.ContentType == "html" {
			if m, err := html2text.FromString(message, html2text.Options{
				PrettyTables: true,
			}); err != nil {
				zap.S().Errorf("error building html message: %v", err)
			} else {
				message = m
			}
		}
	}
	if userName == "" && message == "" {
		return ""
	}
	return fmt.Sprintf(messageTemplate, userName, message)
}

func (c *MSTeams) loadChats(ctx context.Context, msDrive *microsoftcore.MSDrive, nextLink string) error {
	var response microsoftcore.MSTeamsChatResponse
	url := nextLink
	if url == "" {
		url = msTeamsChats
	}
	if err := c.requestAndParse(ctx, url, &response); err != nil {
		return nil
	}
	for _, chat := range response.Value {
		sourceID := fmt.Sprintf("chat:%s", chat.Id)
		state, ok := c.state.Chats[chat.Id]
		if !ok {
			state = &MSTeamMessageState{
				LastCreatedDateTime: time.Time{},
			}
			c.state.Chats[chat.Id] = state
		}

		result, err := c.loadChatMessages(ctx, msDrive, state, chat.Id, fmt.Sprintf(msTeamsChatMessagesURL, chat.Id))
		if err != nil {
			zap.S().Errorf("error loading chat messages: %s", err.Error())
			continue
		}
		if len(result) == 0 {
			continue
		}
		doc := &model.Document{
			SourceID:        sourceID,
			ConnectorID:     c.model.ID,
			URL:             "",
			ChunkingSession: c.sessionID,
			Analyzed:        false,
			CreationDate:    time.Now().UTC(),
			LastUpdate:      pg.NullTime{time.Now().UTC()},
			OriginalURL:     chat.WebUrl,
			IsExists:        true,
		}
		c.model.DocsMap[sourceID] = doc

		fileName := utils.StripFileName(fmt.Sprintf("%s_%s.md", uuid.New().String(), chat.Id))
		c.resultCh <- &Response{
			URL:        doc.URL,
			Name:       fileName,
			SourceID:   sourceID,
			DocumentID: doc.ID.IntPart(),
			MimeType:   "text/markdown",
			FileType:   proto.FileType_MD,
			Signature:  "",
			Content: &Content{
				Bucket:        model.BucketName(c.model.User.EmbeddingModel.TenantID),
				URL:           "",
				AppendContent: true,
				Body:          []byte(strings.Join(result, "\n")),
			},
			UpToData: false,
		}
	}
	if response.NexLink != "" {
		return c.loadChats(ctx, msDrive, response.NexLink)
	}
	return nil
}
func (c *MSTeams) loadChatMessages(ctx context.Context,
	msDrive *microsoftcore.MSDrive,
	state *MSTeamMessageState,
	chatID, url string) ([]string, error) {
	var response microsoftcore.MessageResponse
	if err := c.requestAndParse(ctx, url, &response); err != nil {
		return nil, err
	}
	lastTime := state.LastCreatedDateTime.UTC()

	var messages []string

	for _, msg := range response.Value {
		// do not scan system messages
		if msg.MessageType != messageTypeMessage {
			continue
		}
		if state.LastCreatedDateTime.UTC().After(msg.CreatedDateTime.UTC()) ||
			state.LastCreatedDateTime.UTC().Equal(msg.CreatedDateTime.UTC()) {
			// messages in desc order. not needed to process messages that were loaded before.
			return messages, nil
		}

		// renew newest message time
		if lastTime.UTC().Before(msg.CreatedDateTime.UTC()) {
			lastTime = msg.CreatedDateTime
		}
		if message := c.buildMDMessage(msg); message != "" {
			messages = append(messages, message)
		}
		for _, attachment := range msg.Attachments {
			if err := c.loadAttachment(ctx, msDrive, attachment); err != nil {
				zap.S().Errorf("error loading attachment: %v", err)
			}
		}
	}

	if response.OdataNextLink != "" {
		if nested, err := c.loadChatMessages(ctx, msDrive, state, chatID, response.OdataNextLink); err == nil {
			messages = append(messages, nested...)
		} else {
			zap.S().Errorf("error loading nested chat messages: %v", err)
		}

	}
	state.LastCreatedDateTime = lastTime
	return messages, nil
}

func (c *MSTeams) loadAttachment(ctx context.Context, msDrive *microsoftcore.MSDrive, attachment *microsoftcore.Attachment) error {

	if attachment.ContentType != attachmentContentTypReference {
		// do not scrap replies
		return nil
	}
	if err := msDrive.DownloadItem(ctx, attachment.Id, c.fileSizeLimit); err != nil {
		zap.S().Errorf("download file %s", err.Error())
	}
	return nil
}

// NewMSTeams creates new instance of MsTeams connector
func NewMSTeams(connector *model.Connector,
	connectorRepo repository.ConnectorRepository,
	oauthURL string) (Connector, error) {
	conn := MSTeams{
		Base: Base{
			connectorRepo: connectorRepo,
			oauthClient: resty.New().
				SetTimeout(time.Minute).
				SetBaseURL(oauthURL),
		},
		param: &MSTeamParameters{},
		state: &MSTeamState{},
	}
	conn.Base.Config(connector)

	if err := connector.ConnectorSpecificConfig.ToStruct(conn.param); err != nil {
		return nil, err
	}

	newToken, err := conn.refreshToken(conn.param.Token)
	if err != nil {
		return nil, err
	}
	if newToken != nil {
		conn.param.Token = newToken
	}
	if err = connector.State.ToStruct(conn.state); err != nil {
		zap.S().Infof("can not parse state %v", err)
	}
	if conn.state.Channels == nil {
		conn.state.Channels = make(map[string]*MSTeamChannelState)
	}
	if conn.state.Chats == nil {
		conn.state.Chats = make(map[string]*MSTeamMessageState)
	}
	conn.client = resty.New().
		SetTimeout(time.Minute).
		SetHeader(authorizationHeader, fmt.Sprintf("%s %s",
			conn.param.Token.TokenType,
			conn.param.Token.AccessToken))
	return &conn, nil
}
