package models

import (
	"InfluenceIQ/config"
	"context"
	"time"
)

// Profile represents user profile data.
type Profile struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	DisplayName    string    `json:"display_name"`
	AvatarURL      string    `json:"avatar_url,omitempty"`
	Bio            string    `json:"bio,omitempty"`
	AccountType    string    `json:"account_type"` // "influencer" or "brand"
	Category       string    `json:"category,omitempty"`
	FollowerCount  int       `json:"follower_count,omitempty"`
	EngagementRate float64   `json:"engagement_rate,omitempty"`
	CompanyName    string    `json:"company_name,omitempty"`
	Industry       string    `json:"industry,omitempty"`
	Website        string    `json:"website,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func CreateProfile(ctx context.Context, p *Profile) error {
	query := `
		INSERT INTO profiles (
			user_id, display_name, avatar_url, bio, account_type,
			category, follower_count, engagement_rate,
			company_name, industry, website, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return config.DB.QueryRow(ctx, query,
		p.UserID, p.DisplayName, p.AvatarURL, p.Bio, p.AccountType,
		p.Category, p.FollowerCount, p.EngagementRate,
		p.CompanyName, p.Industry, p.Website,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func GetProfileByUserID(ctx context.Context, userID int) (*Profile, error) {
	var p Profile
	query := `
		SELECT id, user_id, display_name, avatar_url, bio, account_type,
		       category, follower_count, engagement_rate,
		       company_name, industry, website, created_at, updated_at
		FROM profiles
		WHERE user_id = $1
	`
	err := config.DB.QueryRow(ctx, query, userID).Scan(
		&p.ID, &p.UserID, &p.DisplayName, &p.AvatarURL, &p.Bio, &p.AccountType,
		&p.Category, &p.FollowerCount, &p.EngagementRate,
		&p.CompanyName, &p.Industry, &p.Website, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func UpdateProfile(ctx context.Context, p *Profile) error {
	query := `
		UPDATE profiles
		SET display_name = $1, avatar_url = $2, bio = $3, account_type = $4,
			category = $5, follower_count = $6, engagement_rate = $7,
			company_name = $8, industry = $9, website = $10, updated_at = NOW()
		WHERE user_id = $11
	`
	_, err := config.DB.Exec(ctx, query,
		p.DisplayName, p.AvatarURL, p.Bio, p.AccountType,
		p.Category, p.FollowerCount, p.EngagementRate,
		p.CompanyName, p.Industry, p.Website, p.UserID,
	)
	return err
}

func DeleteProfile(ctx context.Context, userID int) error {
	_, err := config.DB.Exec(ctx, `DELETE FROM profiles WHERE user_id = $1`, userID)
	return err
}
