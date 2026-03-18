package routes

import (
	"github.com/Shubham23jha/digital-post-office/internal/bootstrap"
	"github.com/Shubham23jha/digital-post-office/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, app *bootstrap.App) {

	api := r.Group("/api")

	// =========================
	// AUTH ROUTES
	// =========================

	auth := api.Group("/auth")
	{
		// Public APIs
		auth.POST("/signup", app.AuthHandler.Register)

		auth.POST("/login", app.AuthHandler.Login)

		auth.POST("/refresh", app.AuthHandler.Refresh)

		auth.POST("/logout", app.AuthHandler.Logout)

		// Protected APIs
		authProtected := auth.Group("/")
		authProtected.Use(
			middleware.AuthMiddleware(),
		)
		{
			authProtected.POST(
				"/logout-all",
				app.AuthHandler.LogoutAll,
			)
		}
	}

	// =========================
	// CAMPAIGN ROUTES
	// =========================
	campaigns := api.Group("/campaigns")
	{
		campaigns.GET("/", app.CampaignHandler.List)
		campaigns.POST("/", app.CampaignHandler.Create)
		campaigns.GET("/:id", app.CampaignHandler.Get)
	}
}
