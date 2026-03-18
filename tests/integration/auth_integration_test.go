package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Shubham23jha/digital-post-office/internal/bootstrap"
	"github.com/Shubham23jha/digital-post-office/internal/models"
	"github.com/Shubham23jha/digital-post-office/internal/routes"
	"github.com/Shubham23jha/digital-post-office/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthFlow(t *testing.T) {
	// Initialize Gin engine
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Initialize App using wire-generated code
	log.Printf("database.DB in TestAuthFlow: %p", database.DB)
	if database.DB == nil {
		t.Fatal("database.DB is nil!")
	}
	app, err := bootstrap.InitializeApp(database.DB)
	assert.NoError(t, err)

	// Register routes
	routes.Register(r, app)

	// Test Signup
	t.Run("SignUp", func(t *testing.T) {
		signupPayload := map[string]string{
			"firstName":   "Test",
			"lastName":    "User",
			"email":       "test@example.com",
			"password":    "password123",
			"phoneNumber": "1234567890",
		}
		jsonPayload, _ := json.Marshal(signupPayload)

		req, _ := http.NewRequest("POST", "/api/auth/signup", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "User created")
	})

	// Test Login
	t.Run("Login", func(t *testing.T) {
		loginPayload := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
			"deviceID": "test-device",
		}
		jsonPayload, _ := json.Marshal(loginPayload)

		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "access_token")

		// Verify cookies (refresh_token)
		cookies := w.Result().Cookies()
		found := false
		for _, cookie := range cookies {
			if cookie.Name == "refresh_token" {
				found = true
				assert.NotEmpty(t, cookie.Value)
			}
		}
		assert.True(t, found, "refresh_token cookie should be set")
	})
}
func TestSessionLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	app, err := bootstrap.InitializeApp(database.DB)
	assert.NoError(t, err)

	routes.Register(router, app)

	email := "limit@example.com"
	password := "password123"

	// 1. Register
	signupPayload := models.RegisterRequest{
		FirstName:   "Limit",
		LastName:    "User",
		Email:       email,
		PhoneNumber: "9999999999",
		Password:    password,
	}
	body, _ := json.Marshal(signupPayload)
	reqSignup, _ := http.NewRequest("POST", "/api/auth/signup", bytes.NewBuffer(body))
	wSignup := httptest.NewRecorder()
	router.ServeHTTP(wSignup, reqSignup)
	assert.Equal(t, 200, wSignup.Code)

	// 2. Login 3 times with different device IDs
	for i := 1; i <= 3; i++ {
		loginPayload := models.LoginRequest{
			Email:    email,
			Password: password,
			DeviceID: fmt.Sprintf("device-%d", i),
		}
		loginBody, _ := json.Marshal(loginPayload)
		reqLogin, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(loginBody))
		wLogin := httptest.NewRecorder()
		router.ServeHTTP(wLogin, reqLogin)

		assert.Equal(t, 200, wLogin.Code, "Login %d should succeed", i)
		time.Sleep(time.Second) // Ensure different IssuedAt for tokens
	}

	// 3. 4th Login should fail with 401 and "device limit reached"
	finalLoginPayload := models.LoginRequest{
		Email:    email,
		Password: password,
		DeviceID: "device-4",
	}
	finalLoginBody, _ := json.Marshal(finalLoginPayload)
	reqFinal, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(finalLoginBody))
	wFinal := httptest.NewRecorder()
	router.ServeHTTP(wFinal, reqFinal)

	assert.Equal(t, 401, wFinal.Code)
	assert.Contains(t, wFinal.Body.String(), "device limit reached")
}
