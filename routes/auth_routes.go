package routes

import (
	"InfluenceIQ/controllers"
	"InfluenceIQ/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/signup", controllers.Signup)
		auth.POST("/login", controllers.Login)
	}

	// Protected Profile Routes
	r.GET("/profile/me", middleware.AuthMiddleware(), controllers.GetMyProfileHandler)
	r.PUT("/profile/update", middleware.AuthMiddleware(), controllers.UpdateMyProfileHandler)
	r.POST("/profile/create", middleware.AuthMiddleware(), controllers.CreateProfileHandler)
	r.DELETE("/profile/delete", middleware.AuthMiddleware(), controllers.DeleteMyProfileHandler)

	// Protected Campaign Routes
	campaign := r.Group("/campaign")
	campaign.Use(middleware.AuthMiddleware())
	{
		campaign.POST("/", controllers.CreateCampaign)
		campaign.GET("/", controllers.GetAllCampaigns)
		campaign.GET("/me", controllers.GetMyCampaigns)
		campaign.GET("/:id", controllers.GetCampaignByID)
		campaign.DELETE("/:id", controllers.DeleteCampaign)
	}

	// Protected Applications
	app := r.Group("/application")
	app.Use(middleware.AuthMiddleware())
	{
		app.POST("/apply/:id", controllers.ApplyToCampaign)
		app.GET("/my", controllers.GetMyApplications)
		app.GET("/campaign/:id", controllers.GetApplicationsForCampaign)
		app.PUT("/:id/status", controllers.UpdateApplicationStatus)
	}

	//  Public AI Endpoints (NO AUTH)
	ai := r.Group("/ai")
	{
		ai.POST("/campaign-idea", controllers.GenerateCampaignIdea)
		ai.POST("/recommend", controllers.RecommendInfluencers)
		ai.POST("/captions", controllers.GenerateCaptions)
	}
}
