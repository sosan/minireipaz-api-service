package routes

import (
	"minireipaz/pkg/common"
	"minireipaz/pkg/interfaces/controllers"
	"minireipaz/pkg/interfaces/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(app *gin.Engine, workflowController *controllers.WorkflowController, userController *controllers.UserController, dashboardController *controllers.DashboardController, authController *controllers.AuthController) {
	app.NoRoute(ErrRouter)

	// Routes in groups
	api := app.Group("/api")
	{
		api.GET("/ping", common.Ping)

		workflows := api.Group("/workflows")
		{
			workflows.POST("", middlewares.ValidateWorkflow(), workflowController.CreateWorkflow)
			// workflows.GET("/:uuid", workflowController.GetWorkflow)
		}

		users := api.Group("/users")
		{
			users.POST("", middlewares.ValidateUser(), userController.SyncUseWrithIDProvider)
			users.GET("/:stub", userController.GetUserByStub)
		}

		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("/:id", middlewares.ValidateUserAuth(), dashboardController.GetUserDashboardByID)
		}

		auth := api.Group("/auth")
		{
			auth.GET("/verify/:id", authController.VerifyUserToken)
		}
	}
}

func ErrRouter(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"error": "Page not found",
	})
}
