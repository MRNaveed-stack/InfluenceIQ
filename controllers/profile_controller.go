package controllers

import (
	"InfluenceIQ/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// POST /api/profile/create
func CreateProfileHandler(c *gin.Context) {
	var input models.Profile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID := c.GetInt("user_id") // from JWT middleware
	input.UserID = userID

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := models.CreateProfile(ctx, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile created", "profile": input})
}

// GET /api/profile/me
func GetMyProfileHandler(c *gin.Context) {
	userID := c.GetInt("user_id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profile, err := models.GetProfileByUserID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// PUT /api/profile/update
func UpdateMyProfileHandler(c *gin.Context) {
	var input models.Profile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID := c.GetInt("user_id")
	input.UserID = userID

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := models.UpdateProfile(ctx, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
}

// DELETE /api/profile/delete
func DeleteMyProfileHandler(c *gin.Context) {
	userID := c.GetInt("user_id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := models.DeleteProfile(ctx, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile deleted"})
}
