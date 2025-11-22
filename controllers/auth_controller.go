package controllers

import (
	"InfluenceIQ/config"
	"InfluenceIQ/models"
	"InfluenceIQ/utils"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ---------- SIGNUP ----------
func Signup(c *gin.Context) {
	var input models.SignupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Check if email or username already exists
	var exists bool
	err := config.DB.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1 OR username=$2)`,
		input.Email, input.Username).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Email or username already in use"})
		return
	}

	// 2. Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 3. Insert new user
	var userID int
	query := `
		INSERT INTO users (username, email, password, full_name, role)
		VALUES ($1, $2, $3, $4, COALESCE($5, 'viewer'))
		RETURNING user_id
	`
	err = config.DB.QueryRow(ctx, query,
		strings.TrimSpace(input.Username),
		strings.TrimSpace(input.Email),
		string(hashed),
		strings.TrimSpace(input.FullName),
		"viewer",
	).Scan(&userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// 4. Generate JWT token
	token, err := utils.GenerateToken(userID, 0, "viewer")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Signup successful",
		"user": gin.H{
			"id":       userID,
			"username": input.Username,
			"email":    input.Email,
			"fullName": input.FullName,
			"role":     "viewer",
		},
		"token": token,
	})
}

type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

// ---------- LOGIN ----------
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var userID int
	var email, username, fullName, role, hashedPassword string

	query := `
		SELECT user_id, email, username, full_name, role, password
		FROM users
		WHERE email = $1 OR username = $2
	`
	err := config.DB.QueryRow(ctx, query, strings.TrimSpace(req.EmailOrUsername), strings.TrimSpace(req.EmailOrUsername)).
		Scan(&userID, &email, &username, &fullName, &role, &hashedPassword)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateToken(userID, 0, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":       userID,
			"username": username,
			"email":    email,
			"fullName": fullName,
			"role":     role,
		},
		"token": token,
	})
}
