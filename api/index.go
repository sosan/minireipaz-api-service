package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"minireipaz/pkg/config"
	"minireipaz/pkg/di"
	"minireipaz/pkg/honeycomb"
	"minireipaz/pkg/interfaces/middlewares"
	"minireipaz/pkg/interfaces/routes"
	"net/http"
)

var (
	app *gin.Engine
)

// Init initializes the application without starting the server.
func init() {
	InitApp()
}

// InitApp initializes the Gin application.
func InitApp() {
	log.Print("---- Initializing App ----")
	config.LoadEnvs(".")

	// Setup OpenTelemetry
	ctx := context.Background()
	tp, exp, err := honeycomb.SetupHoneyComb(ctx)
	if err != nil {
		log.Panicf("ERROR | Failed to initialize OpenTelemetry: %v", err)
	}

	// Ensure sub processes and telemetry are exported correctly.
	defer func() {
		_ = exp.Shutdown(ctx)
		_ = tp.Shutdown(ctx)
	}()

	// Initialize Gin app
	gin.SetMode(gin.ReleaseMode)
	app = gin.New()

	// Dependency injection and routes setup
	worflowcontroller, authService, userController, dashboardController, authController := di.InitDependencies()
	middlewares.Register(app, authService)
	routes.Register(app, worflowcontroller, userController, dashboardController, authController)
}

// Handler is the main function that Vercel calls to handle HTTP requests.
func Handler(w http.ResponseWriter, r *http.Request) {
	// If app is not initialized, initialize it
	if app == nil {
		InitApp()
	}
	// Use Gin to serve the HTTP request
	app.ServeHTTP(w, r)
}

func Dummy() {
	RunWebserver()
}

func RunWebserver() {
	addr := config.GetEnv("FRONTEND_ADDR", ":3020")
	err := app.Run(addr)
	if err != nil {
		log.Panicf("ERROR | Starting gin failed, %v", err)
	}
}
