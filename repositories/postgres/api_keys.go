package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/ggarrido85/api-backend/models"
	"github.com/ggarrido85/api-backend/repositories/dbmodels"
)

func (db *Database) GetApiKeyByHash(ctx context.Context, hash []byte) (models.ApiKey, error) {
	query := `
		SELECT id, org_id, prefix, description, partner_id, role
		FROM api_keys
		WHERE key_hash = $1
		AND deleted_at IS NULL
	`

	var apiKey dbmodels.DBApiKey
	err := db.pool.QueryRow(ctx, query, hash).Scan(
		&apiKey.Id,
		&apiKey.OrganizationId,
		&apiKey.Prefix,
		&apiKey.Description,
		&apiKey.PartnerId,
		&apiKey.Role,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.ApiKey{}, models.NotFoundError
	}
	if err != nil {
		return models.ApiKey{}, fmt.Errorf("pool.QueryRow error: %w", err)
	}
	return dbmodels.AdaptApikey(apiKey)
}
