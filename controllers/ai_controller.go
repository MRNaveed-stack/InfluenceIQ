package controllers

import (
	"fmt"
	"net/http"

	"InfluenceIQ/services"

	"github.com/gin-gonic/gin"
)

func GenerateCampaignIdea(c *gin.Context) {
	var req struct {
		Product      string  `json:"product" binding:"required"`
		TargetMarket string  `json:"target_market"`
		Budget       float64 `json:"budget"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	prompt := fmt.Sprintf(
		"You are a creative marketing AI. Suggest a short campaign idea for the product \"%s\". Target audience: %s. Budget: %.2f USD. Return only 2–3 sentences, catchy and actionable.",
		req.Product, req.TargetMarket, req.Budget,
	)

	text, err := services.CallGemini(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "idea": text})
}

func RecommendInfluencers(c *gin.Context) {
	var req struct {
		Category string  `json:"category" binding:"required"`
		Budget   float64 `json:"budget"`
		Audience string  `json:"audience"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	prompt := fmt.Sprintf(
		"You are an influencer marketing strategist AI. Recommend 3 ideal influencer profiles for a campaign in the \"%s\" category. Consider a budget of %.2f USD and target audience: %s. Return influencers as bullet points with short justification.",
		req.Category, req.Budget, req.Audience,
	)

	text, err := services.CallGemini(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "recommendations": text})
}

func GenerateCaptions(c *gin.Context) {
	var req struct {
		Theme string `json:"theme" binding:"required"`
		Tone  string `json:"tone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	prompt := fmt.Sprintf(
		"Write 3 short social media captions about \"%s\" in a %s tone. Make them engaging and unique.",
		req.Theme, req.Tone,
	)

	text, err := services.CallGemini(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "captions": text})
}
