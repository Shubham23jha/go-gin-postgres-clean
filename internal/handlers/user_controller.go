package handlers

import (
	"fmt"
	"net/http"

	"github.com/Shubham23jha/go-gin-postgres-clean/internal/models"
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/service"
	"github.com/Shubham23jha/go-gin-postgres-clean/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
)

type AuthHandler struct {
	userService    service.UserService
	sessionService service.SessionService
}

func NewAuthHandler(
	userService service.UserService,
	sessionService service.SessionService,
) *AuthHandler {
	return &AuthHandler{
		userService:    userService,
		sessionService: sessionService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.Register(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created"})
}

func (h *AuthHandler) Login(c *gin.Context) {

	var req models.LoginRequest

	// =========================
	// Bind Request
	// =========================
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// =========================
	// Device Info Extraction
	// =========================
	uaString := c.Request.UserAgent()
	ua := user_agent.New(uaString)

	browser, _ := ua.Browser()
	os := ua.OS()
	ip := c.ClientIP()

	if req.Email == "" && req.PhoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email or phone number is required",
		})
		return
	}

	var accessToken string
	var refreshToken string
	var err error

	// =========================
	// Login with Email
	// =========================
	if req.Email != "" {
		accessToken, refreshToken, err =
			h.userService.LoginWithEmail(
				req.Email,
				req.Password,
				req.DeviceID,
				os,
				browser,
				ip,
			)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	// =========================
	// Login with Phone
	// =========================
	if req.PhoneNumber != "" {
		accessToken, refreshToken, err =
			h.userService.LoginWithPhone(
				req.PhoneNumber,
				req.Password,
				req.DeviceID,
				os,
				browser,
				ip,
			)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	// =========================
	// Set Refresh Token Cookie
	// =========================
	c.SetCookie(
		"refresh_token",
		refreshToken,
		60*60*24*7, // 7 days
		"/",
		"",
		false, // secure → true in production (HTTPS)
		true,  // httpOnly
	)

	// =========================
	// Send Access Token
	// =========================
	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")

	if err != nil {
		c.JSON(401, gin.H{
			"error": "no active session",
		})
		return
	}

	// Deactivate session in DB
	err = h.sessionService.LogoutByToken(refreshToken)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to logout",
		})
		return
	}

	//Clear cookie
	c.SetCookie(
		"refresh_token",
		"",
		-1, // expire immediately
		"/",
		"",
		false,
		true,
	)

	c.JSON(200, gin.H{
		"message": "logged out successfully",
	})
}

func (h *AuthHandler) LogoutAll(c *gin.Context) {

	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	userID := userIDRaw.(uint)

	err := h.sessionService.LogoutAll(userID)

	if err != nil {
		c.JSON(500, gin.H{"error": "failed"})
		return
	}

	c.SetCookie("refresh_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.JSON(200, gin.H{
		"message": "logged out from all devices",
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {

	//Read refresh cookie
	refreshToken, err := c.Cookie("refresh_token")

	fmt.Println("Refresh Token:", refreshToken)

	if err != nil {
		c.JSON(401, gin.H{
			"error": "refresh token missing",
		})
		return
	}

	claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid"})
		return
	}
	_, err = h.sessionService.ValidateRefreshToken(refreshToken)

	if err != nil {
		c.JSON(401, gin.H{
			"error": err.Error(),
		})
		return
	}
	accessToken, err := utils.GenerateAccessToken(claims.UserID, claims.Email)

	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to generate access token",
		})
		return
	}
	c.JSON(200, gin.H{
		"access_token": accessToken,
	})
}
