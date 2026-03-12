package integration

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Shubham23jha/go-gin-postgres-clean/internal/bootstrap"
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/routes"
	"github.com/Shubham23jha/go-gin-postgres-clean/pkg/database"
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
