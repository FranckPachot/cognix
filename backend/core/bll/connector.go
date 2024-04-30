package bll

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/shopspring/decimal"
	"time"
)

type (
	ConnectorBL interface {
		GetAll(ctx context.Context, user *model.User) ([]*model.Connector, error)
		GetByID(ctx context.Context, user *model.User, id int64) (*model.Connector, error)
		Create(ctx context.Context, user *model.User, param *parameters.CreateConnectorParam) (*model.Connector, error)
		Update(ctx context.Context, id int64, user *model.User, param *parameters.UpdateConnectorParam) (*model.Connector, error)
	}
	connectorBL struct {
		connectorRepo  repository.ConnectorRepository
		credentialRepo repository.CredentialRepository
	}
)

func NewConnectorBL(connectorRepo repository.ConnectorRepository, credentialRepo repository.CredentialRepository) ConnectorBL {
	return &connectorBL{
		connectorRepo:  connectorRepo,
		credentialRepo: credentialRepo,
	}
}

func (c *connectorBL) Create(ctx context.Context, user *model.User, param *parameters.CreateConnectorParam) (*model.Connector, error) {

	conn := model.Connector{
		CredentialID:            param.CredentialID,
		Name:                    param.Name,
		Source:                  model.SourceType(param.Source),
		InputType:               param.InputType,
		ConnectorSpecificConfig: param.ConnectorSpecificConfig,
		RefreshFreq:             param.RefreshFreq,
		UserID:                  user.ID,
		TenantID:                user.TenantID,
		Shared:                  param.Shared,
		Disabled:                param.Disabled,
		CreatedDate:             time.Now().UTC(),
	}
	if param.CredentialID.Valid {
		cred, err := c.credentialRepo.GetByID(ctx, param.CredentialID.Decimal.IntPart(), user.TenantID, user.ID)
		if err != nil {
			return nil, err
		}
		if cred.Source != model.SourceType(param.Source) {
			return nil, utils.InvalidInput.New("wrong credential source")
		}
		conn.CredentialID = decimal.NewNullDecimal(cred.ID)
	}

	if err := c.connectorRepo.Create(ctx, &conn); err != nil {
		return nil, err
	}
	return &conn, nil
}

func (c *connectorBL) Update(ctx context.Context, id int64, user *model.User, param *parameters.UpdateConnectorParam) (*model.Connector, error) {
	conn, err := c.connectorRepo.GetByID(ctx, user.TenantID, user.ID, id)
	if err != nil {
		return nil, err
	}
	if param.CredentialID.Valid {
		cred, err := c.credentialRepo.GetByID(ctx, param.CredentialID.Decimal.IntPart(), user.TenantID, user.ID)
		if err != nil {
			return nil, err
		}
		if cred.Source != conn.Source {
			return nil, utils.InvalidInput.New("wrong credential source")
		}
	}
	conn.ConnectorSpecificConfig = param.ConnectorSpecificConfig
	conn.CredentialID = param.CredentialID
	conn.Name = param.Name
	conn.InputType = param.InputType
	conn.RefreshFreq = param.RefreshFreq
	conn.Shared = param.Shared
	conn.Disabled = param.Disabled
	conn.UpdatedDate = pg.NullTime{time.Now().UTC()}

	//	sql.NullTime{
	//Time:  time.Now().UTC(),
	//Valid: true,
	//	} //  null.TimeFrom(time.Now().UTC())
	if err = c.connectorRepo.Update(ctx, conn); err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *connectorBL) GetAll(ctx context.Context, user *model.User) ([]*model.Connector, error) {
	return c.connectorRepo.GetAll(ctx, user.TenantID, user.ID)
}

func (c *connectorBL) GetByID(ctx context.Context, user *model.User, id int64) (*model.Connector, error) {
	return c.connectorRepo.GetByID(ctx, user.TenantID, user.ID, id)
}
