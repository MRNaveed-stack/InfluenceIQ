package controllers

import (
	"InfluenceIQ/config"
	"InfluenceIQ/models"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// POST /api/campaigns
func CreateCampaign(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "unauthorized"})
		return
	}

	var req struct {
		Title       string    `json:"title" binding:"required"`
		Description string    `json:"description" binding:"required"`
		Category    string    `json:"category"`
		Budget      float64   `json:"budget"`
		Deadline    time.Time `json:"deadline"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	campaign := models.Campaign{
		BrandID:     userID.(int),
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Budget:      req.Budget,
		Deadline:    req.Deadline,
		Status:      "active",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := models.CreateCampaign(ctx, &campaign); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to create campaign"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": campaign})
}

// GET /api/campaigns
func GetAllCampaigns(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Optionally, you can make a helper GetAllCampaigns(ctx)
	query := `
		SELECT id, brand_id, title, description, category, budget, deadline, status, created_at, updated_at
		FROM campaigns
		ORDER BY created_at DESC
	`
	rows, err := config.DB.Query(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to fetch campaigns"})
		return
	}
	defer rows.Close()

	var campaigns []models.Campaign
	for rows.Next() {
		var camp models.Campaign
		if err := rows.Scan(
			&camp.ID, &camp.BrandID, &camp.Title, &camp.Description, &camp.Category,
			&camp.Budget, &camp.Deadline, &camp.Status, &camp.CreatedAt, &camp.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "scan error"})
			return
		}
		campaigns = append(campaigns, camp)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": campaigns})
}

// GET /api/campaigns/mine
func GetMyCampaigns(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "unauthorized"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	campaigns, err := models.GetCampaignByID(ctx, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to fetch campaigns"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": campaigns})
}

// GET /api/campaigns/:id
func GetCampaignByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid campaign id"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	campaign, err := models.GetCampaignByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "campaign not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": campaign})
}

// PUT /api/campaigns/:id
func UpdateCampaign(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid id"})
		return
	}

	var req models.Campaign
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	req.ID = id
	req.BrandID = userID.(int)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := models.UpdateCampaign(ctx, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to update campaign"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "campaign updated"})
}

// DELETE /api/campaigns/:id
func DeleteCampaign(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid id"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := models.DeleteCampaign(ctx, id, userID.(int)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to delete campaign"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "campaign deleted"})
}
