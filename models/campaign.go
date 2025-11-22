package models

import (
	"InfluenceIQ/config"
	"context"
	"time"
)

// Campaign represents a brand campaign record in PostgreSQL
type Campaign struct {
	ID          int       `json:"id"`
	BrandID     int       `json:"brand_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category,omitempty"`
	Budget      float64   `json:"budget"`
	Deadline    time.Time `json:"deadline"`
	Status      string    `json:"status"` // active, closed, draft
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func CreateCampaign(ctx context.Context, c *Campaign) error {
	query := `
		INSERT INTO campaigns (
			brand_id, title, description, category, budget, deadline, status, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, COALESCE($7, 'active'), NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return config.DB.QueryRow(ctx, query,
		c.BrandID, c.Title, c.Description, c.Category, c.Budget, c.Deadline, c.Status,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func GetCampaignByID(ctx context.Context, id int) (*Campaign, error) {
	var c Campaign
	query := `
		SELECT id, brand_id, title, description, category, budget, deadline, status, created_at, updated_at
		FROM campaigns
		WHERE id = $1
	`
	err := config.DB.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.BrandID, &c.Title, &c.Description, &c.Category,
		&c.Budget, &c.Deadline, &c.Status, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func UpdateCampaign(ctx context.Context, c *Campaign) error {
	query := `
		UPDATE campaigns
		SET title = $1, description = $2, category = $3, budget = $4,
		    deadline = $5, status = $6, updated_at = NOW()
		WHERE id = $7 AND brand_id = $8
	`
	_, err := config.DB.Exec(ctx, query,
		c.Title, c.Description, c.Category, c.Budget, c.Deadline, c.Status, c.ID, c.BrandID,
	)
	return err
}

func DeleteCampaign(ctx context.Context, id int, brandID int) error {
	_, err := config.DB.Exec(ctx, `
		DELETE FROM campaigns WHERE id = $1 AND brand_id = $2
	`, id, brandID)
	return err
}
