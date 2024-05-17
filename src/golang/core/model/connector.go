package model

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

const (
	CollectionTenant = "tenant_%s"
	CollectionUser   = "user_%s"

	StatusFailed  = "failed"
	StatusSuccess = "success"
)

type Connector struct {
	tableName               struct{}             `pg:"connectors,omitempty"`
	ID                      decimal.Decimal      `json:"id,omitempty"`
	CredentialID            decimal.NullDecimal  `json:"credential_id,omitempty"`
	Name                    string               `json:"name,omitempty"`
	Source                  SourceType           `json:"source,omitempty"`
	InputType               string               `json:"input_type,omitempty"`
	ConnectorSpecificConfig JSONMap              `json:"connector_specific_config,omitempty"`
	RefreshFreq             int                  `json:"refresh_freq,omitempty"`
	UserID                  uuid.UUID            `json:"user_id,omitempty"`
	TenantID                uuid.UUID            `json:"tenant_id,omitempty"`
	Shared                  bool                 `json:"shared,omitempty" pg:",use_zero"`
	Disabled                bool                 `json:"disabled,omitempty" pg:",use_zero"`
	LastSuccessfulIndexTime pg.NullTime          `json:"last_successful_index_time,omitempty" pg:",use_zero"`
	LastAttemptStatus       string               `json:"last_attempt_status,omitempty"`
	TotalDocsIndexed        int                  `json:"total_docs_indexed" pg:",use_zero"`
	CreatedDate             time.Time            `json:"created_date,omitempty"`
	UpdatedDate             pg.NullTime          `json:"updated_date,omitempty" pg:",use_zero"`
	DeletedDate             pg.NullTime          `json:"deleted_date,omitempty" pg:",use_zero"`
	Credential              *Credential          `json:"credential,omitempty" pg:"rel:has-one,fk:credential_id"`
	Docs                    []*Document          `json:"docs,omitempty" pg:"rel:has-many"`
	DocsMap                 map[string]*Document `json:"docs_map,omitempty" pg:"-"`
}

func (c *Connector) CollectionName() string {
	return CollectionName(c.Shared, c.UserID, c.TenantID)
}
func CollectionName(isShared bool, userID, tenantID uuid.UUID) string {
	if isShared {
		return strings.ReplaceAll(fmt.Sprintf(CollectionTenant, tenantID), "-", "")
	}
	return strings.ReplaceAll(fmt.Sprintf(CollectionUser, userID), "-", "")

}