package health

import (
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"mi-app-backend/db/config"
)

// Handler handles health check requests
type Handler struct{}

// NewHealthHandler creates a new health check handler
func NewHealthHandler() *Handler {
	return &Handler{}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status      string    `json:"status"`
	Version     string    `json:"version"`
	Environment string    `json:"environment"`
	Timestamp   time.Time `json:"timestamp"`
	GUID        string    `json:"guid"`
	Info        InfoData  `json:"info"`
}

// InfoData represents additional health information
type InfoData struct {
	GoVersion  string `json:"goVersion"`
	DatabaseOK bool   `json:"databaseOK"`
	Uptime     string `json:"uptime"`
}

// Status godoc
// @Summary Get API health status
// @Description Returns the health status of the API and its dependencies
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *Handler) Status(c *gin.Context) {
	startTime := time.Now()

	// Test database connection
	dbOK := testDBConnection()

	// Get environment
	env := "development"
	if gin.Mode() == gin.ReleaseMode {
		env = "production"
	}

	// Create response
	response := HealthResponse{
		Status:      "ok",
		Version:     "1.0.0",
		Environment: env,
		Timestamp:   time.Now(),
		GUID:        uuid.New().String(),
		Info: InfoData{
			GoVersion:  runtime.Version(),
			DatabaseOK: dbOK,
			Uptime:     time.Since(startTime).String(),
		},
	}

	c.JSON(http.StatusOK, response)
}

// Check godoc
// @Summary Quick health check endpoint
// @Description Simple health check that returns 200 OK if the API is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health/check [get]
func (h *Handler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// testDBConnection tests the database connection
func testDBConnection() bool {
	client, err := config.NewSupabaseClient()
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return false
	}
	defer client.Close()

	_, err = client.ExecuteQuery()
	if err != nil {
		log.Printf("Error executing test query: %v", err)
		return false
	}

	return true
}
