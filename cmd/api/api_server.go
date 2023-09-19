package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	ciqLogger "github.com/RHEcosystemAppEng/cluster-iq/internal/logger"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	// version reflects the current version of the API
	version string
	// commit reflects the git short-hash of the compiled version
	commit string
	// TODO: comment rest of global vars
	inven  *inventory.Inventory
	router *gin.Engine
	apiURL string
	dbURL  string
	dbPass string
	logger *zap.Logger
	dbpool *pgxpool.Pool
)

const connStr = "postgresql://user:password@localhost:5432/clusteriq?sslmode=disable"

func init() {
	// Logging config
	logger = ciqLogger.NewLogger()

	// Getting config
	apiHost := os.Getenv("CIQ_API_HOST")
	apiPort := os.Getenv("CIQ_API_PORT")
	dbHost := os.Getenv("CIQ_DB_HOST")
	dbPort := os.Getenv("CIQ_DB_PORT")
	dbPass = os.Getenv("CIQ_DB_PASS")
	apiURL = fmt.Sprintf("%s:%s", apiHost, apiPort)
	dbURL = fmt.Sprintf("%s:%s", dbHost, dbPort)

	// Initializaion global vars
	inven = inventory.NewInventory()
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()

	// Configure GIN to use ZAP
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
}

func addHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
}

func main() {
	// Ignore Logger sync error
	defer func() { _ = logger.Sync() }()

	logger.Info("Starting ClusterIQ API", zap.String("version", version), zap.String("commit", commit))
	logger.Info("Connection properties", zap.String("api_url", apiURL), zap.String("db_url", dbURL))
	logger.Debug("Debug Mode active!")

	// Preparing API Endpoints
	baseGroup := router.Group("/api/v1")
	{
		instancesGroup := baseGroup.Group("/instances")
		{
			instancesGroup.GET("", HandlerGetInstances)
			instancesGroup.GET("/:instance_id", HandlerGetInstanceByID)
			instancesGroup.POST("", HandlerPostInstance)
			instancesGroup.DELETE("/:instance_id", HandlerDeleteInstance)
			instancesGroup.PATCH("/:instance_id", HandlerPatchInstance)
		}
	}

	// PGSQL connection
	var err error
	dbpool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		logger.Error("Can't connect to PSQL DB", zap.Error(err))
	}
	defer dbpool.Close()

	// Start API
	logger.Info("API Ready to serve")
	router.Run(apiURL)
}
