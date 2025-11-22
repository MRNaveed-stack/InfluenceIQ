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

// POST /api/campaigns/:id/apply
func ApplyToCampaign(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "unauthorized"})
		return
	}

	campaignID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid campaign id"})
		return
	}

	var req struct {
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Prevent duplicate applications
	if _, err := models.GetApplicationByCampaignAndInfluencer(ctx, campaignID, userID.(int)); err == nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "error": "already applied"})
		return
	}

	app := models.CampaignApplication{
		CampaignID:   campaignID,
		InfluencerID: userID.(int),
		Message:      req.Message,
		Status:       "pending",
	}

	if err := models.CreateApplication(ctx, &app); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to apply"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": app})
}

// GET /api/applications/mine
func GetMyApplications(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "unauthorized"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	apps, err := models.GetApplicationsByInfluencer(ctx, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to fetch applications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": apps})
}

// GET /api/campaigns/:id/applications
func GetApplicationsForCampaign(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "unauthorized"})
		return
	}

	campaignID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid campaign id"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Verify campaign ownership
	var brandID int
	err = config.DB.QueryRow(ctx, `SELECT brand_id FROM campaigns WHERE id = $1`, campaignID).Scan(&brandID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "campaign not found"})
		return
	}
	if brandID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "not your campaign"})
		return
	}

	apps, err := models.GetApplicationsByCampaign(ctx, campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to fetch applications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": apps})
}

// PUT /api/applications/:id/status
func UpdateApplicationStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "unauthorized"})
		return
	}

	appID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid application id"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=pending accepted rejected"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var campaignID, brandID int
	err = config.DB.QueryRow(ctx,
		`SELECT c.id, c.brand_id FROM campaigns c
		 JOIN campaign_applications a ON a.campaign_id = c.id
		 WHERE a.id = $1`, appID).Scan(&campaignID, &brandID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "application or campaign not found"})
		return
	}
	if brandID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "not your campaign"})
		return
	}

	if err := models.UpdateApplicationStatus(ctx, appID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "status updated"})
}
