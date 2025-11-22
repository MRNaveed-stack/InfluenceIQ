package models

import (
	"InfluenceIQ/config"
	"context"
	"time"
)

type CampaignApplication struct {
	ID           int       `json:"id"`
	CampaignID   int       `json:"campaign_id"`
	InfluencerID int       `json:"influencer_id"`
	Status       string    `json:"status"`
	Message      string    `json:"message"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func CreateApplication(ctx context.Context, a *CampaignApplication) error {
	query := `
		INSERT INTO campaign_applications (campaign_id, influencer_id, status, message, created_at, updated_at)
		VALUES ($1, $2, COALESCE($3, 'pending'), $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return config.DB.QueryRow(ctx, query,
		a.CampaignID, a.InfluencerID, a.Status, a.Message,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func GetApplicationByCampaignAndInfluencer(ctx context.Context, campaignID, influencerID int) (*CampaignApplication, error) {
	var a CampaignApplication
	query := `
		SELECT id, campaign_id, influencer_id, status, message, created_at, updated_at
		FROM campaign_applications
		WHERE campaign_id = $1 AND influencer_id = $2
	`
	err := config.DB.QueryRow(ctx, query, campaignID, influencerID).Scan(
		&a.ID, &a.CampaignID, &a.InfluencerID, &a.Status, &a.Message, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func GetApplicationsByInfluencer(ctx context.Context, influencerID int) ([]CampaignApplication, error) {
	query := `
		SELECT id, campaign_id, influencer_id, status, message, created_at, updated_at
		FROM campaign_applications
		WHERE influencer_id = $1
		ORDER BY created_at DESC
	`
	rows, err := config.DB.Query(ctx, query, influencerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []CampaignApplication
	for rows.Next() {
		var a CampaignApplication
		if err := rows.Scan(
			&a.ID, &a.CampaignID, &a.InfluencerID, &a.Status, &a.Message, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		apps = append(apps, a)
	}
	return apps, nil
}

func GetApplicationsByCampaign(ctx context.Context, campaignID int) ([]CampaignApplication, error) {
	query := `
		SELECT id, campaign_id, influencer_id, status, message, created_at, updated_at
		FROM campaign_applications
		WHERE campaign_id = $1
		ORDER BY created_at DESC
	`
	rows, err := config.DB.Query(ctx, query, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []CampaignApplication
	for rows.Next() {
		var a CampaignApplication
		if err := rows.Scan(
			&a.ID, &a.CampaignID, &a.InfluencerID, &a.Status, &a.Message, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		apps = append(apps, a)
	}
	return apps, nil
}

func UpdateApplicationStatus(ctx context.Context, appID int, newStatus string) error {
	query := `
		UPDATE campaign_applications
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := config.DB.Exec(ctx, query, newStatus, appID)
	return err
}
