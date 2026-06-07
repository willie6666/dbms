package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"vapor_auror_backend/controllers"
	"vapor_auror_backend/database"
	"vapor_auror_backend/models"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestUpdateGame(t *testing.T) {
	// Setup in-memory DB
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	database.DB = db

	// Migrate models
	err = db.AutoMigrate(&models.Game{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}
	err = db.AutoMigrate(&models.GameMedia{})
	if err != nil {
		t.Fatalf("Failed to migrate media: %v", err)
	}

	// Insert test data
	testGame := models.Game{
		DeveloperID: 1,
		Title:       "Test Game",
		Description: "Old Description",
		Price:       100.0,
	}
	db.Create(&testGame)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Mock middleware to set user_id and role
	router.Use(func(c *gin.Context) {
		c.Set("user_id", float64(1)) // developerID = 1
		c.Set("role", "DEVELOPER")
		c.Next()
	})

	router.PUT("/api/developer/games/:id", controllers.UpdateGame)

	// Create request
	updateData := controllers.UpdateGameInput{
		Price: 250.0,
		Desc:  "Updated Description from Test",
	}
	body, _ := json.Marshal(updateData)

	req, _ := http.NewRequest(http.MethodPut, "/api/developer/games/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Fetch from DB to verify
	var updatedGame models.Game
	db.First(&updatedGame, testGame.GameID)

	if updatedGame.Description != "Updated Description from Test" {
		t.Errorf("Expected description 'Updated Description from Test', got '%s'", updatedGame.Description)
	}
	if updatedGame.Price != 250.0 {
		t.Errorf("Expected price 250.0, got %f", updatedGame.Price)
	}
	
	t.Logf("Successfully updated game! title=%s, new_desc=%s, new_price=%f", updatedGame.Title, updatedGame.Description, updatedGame.Price)
}
