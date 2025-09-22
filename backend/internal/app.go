package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/imlargo/go-api-template/api/docs"
	"github.com/imlargo/go-api-template/internal/cache"
	"github.com/imlargo/go-api-template/internal/config"
	"github.com/imlargo/go-api-template/internal/handlers"
	"github.com/imlargo/go-api-template/internal/metrics"
	"github.com/imlargo/go-api-template/internal/middleware"
	"github.com/imlargo/go-api-template/internal/services"
	"github.com/imlargo/go-api-template/internal/store"
	"github.com/imlargo/go-api-template/pkg/jwt"
	"github.com/imlargo/go-api-template/pkg/kv"
	"github.com/imlargo/go-api-template/pkg/push"
	"github.com/imlargo/go-api-template/pkg/ratelimiter"
	"github.com/imlargo/go-api-template/pkg/sse"
	"github.com/imlargo/go-api-template/pkg/storage"
	"github.com/imlargo/go-api-template/pkg/utils"
	"go.uber.org/zap"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Application struct {
	Config      config.AppConfig
	Store       *store.Store
	Storage     storage.FileStorage
	Metrics     metrics.MetricsService
	Cache       kv.KeyValueStore
	CacheKeys   *cache.CacheKeys
	RateLimiter ratelimiter.RateLimiter
	Logger      *zap.SugaredLogger
	Router      *gin.Engine
}

func (app *Application) Mount() {

	jwtAuth := jwt.NewJwt(jwt.Config{
		Secret:   app.Config.Auth.JwtSecret,
		Issuer:   app.Config.Auth.JwtIssuer,
		Audience: app.Config.Auth.JwtAudience,
	})

	// Adapters
	sseManager := sse.NewSSEManager()
	pushNotificationDispatcher := push.NewPushNotifier(app.Config.PushNotification.VAPIDPrivateKey, app.Config.PushNotification.VAPIDPublicKey)

	// Services
	serviceContainer := services.NewService(app.Store, app.Logger, &app.Config, app.CacheKeys, app.Cache)
	userService := services.NewUserService(serviceContainer)
	authService := services.NewAuthService(serviceContainer, userService, jwtAuth)
	fileService := services.NewFileService(serviceContainer, app.Storage)
	notificationService := services.NewNotificationService(serviceContainer, sseManager, pushNotificationDispatcher)

	// Platform services
	courseService := services.NewCourseService(serviceContainer)
	moduleService := services.NewModuleService(serviceContainer)
	contentService := services.NewContentService(serviceContainer)
	evaluationService := services.NewEvaluationService(serviceContainer)
	// questionService := services.NewQuestionService(serviceContainer)
	answerService := services.NewAnswerService(serviceContainer)
	enrollmentService := services.NewEnrollmentService(serviceContainer)
	evaluationAttemptService := services.NewEvaluationAttemptService(serviceContainer, answerService)
	// userProgressService := services.NewUserProgressService(serviceContainer, enrollmentService)

	// Handlers
	handlerContainer := handlers.NewHandler(app.Logger)
	authHandler := handlers.NewAuthHandler(handlerContainer, authService)
	notificationHandler := handlers.NewNotificationHandler(handlerContainer, notificationService)
	fileHandler := handlers.NewFileHandler(handlerContainer, fileService)

	// Platform handlers
	courseHandler := handlers.NewCourseHandler(handlerContainer, courseService)
	moduleHandler := handlers.NewModuleHandler(handlerContainer, moduleService)
	contentHandler := handlers.NewContentHandler(handlerContainer, contentService)
	evaluationHandler := handlers.NewEvaluationHandler(handlerContainer, evaluationService)
	enrollmentHandler := handlers.NewEnrollmentHandler(handlerContainer, enrollmentService)
	evaluationAttemptHandler := handlers.NewEvaluationAttemptHandler(handlerContainer, evaluationAttemptService)

	// Middlewares
	apiKeyMiddleware := middleware.ApiKeyMiddleware(app.Config.Auth.ApiKey)
	authMiddleware := middleware.AuthTokenMiddleware(jwtAuth)
	metricsMiddleware := middleware.NewMetricsMiddleware(app.Metrics)
	rateLimiterMiddleware := middleware.NewRateLimiterMiddleware(app.RateLimiter)
	corsMiddleware := middleware.NewCorsMiddleware(app.Config.Server.Host, []string{"http://localhost:5173"})

	// Metrics
	app.Router.GET("/internal/metrics", middleware.BearerApiKeyMiddleware(app.Config.Auth.ApiKey), gin.WrapH(promhttp.Handler()))

	// Register middlewares
	app.Router.Use(metricsMiddleware)
	app.Router.Use(corsMiddleware)
	if app.Config.RateLimiter.Enabled {
		app.Router.Use(rateLimiterMiddleware)
	}

	app.registerDocs()

	// Routes
	app.Router.POST("/auth/login", authHandler.Login)
	app.Router.POST("/auth/register", authHandler.Register)
	app.Router.GET("/auth/me", authMiddleware, authHandler.GetUserInfo)

	app.Router.GET("/api/v1/notifications/subscribe", notificationHandler.SubscribeSSE)

	v1 := app.Router.Group("/api/v1", authMiddleware)

	// Files
	v1.GET("/files/:id/download", fileHandler.DownloadFile)

	// Notifications
	v1.GET("/notifications", notificationHandler.GetUserNotifications)
	v1.POST("/notifications/read", notificationHandler.MarkNotificationsAsRead)

	v1.POST("/notifications/send", apiKeyMiddleware, notificationHandler.DispatchSSE)
	v1.POST("/notifications/unsubscribe", notificationHandler.UnsubscribeSSE)
	v1.GET("/notifications/subscriptions", notificationHandler.GetSSESubscriptions)
	v1.POST("/notifications/push/send", apiKeyMiddleware, notificationHandler.DispatchPush)
	v1.POST("/notifications/push/subscribe/:userID", notificationHandler.SubscribePush)
	v1.GET("/notifications/push/subscriptions/:id", notificationHandler.GetPushSubscription)
	v1.POST("/notifications/dispatch", notificationHandler.DispatchNotification)

	// Courses
	v1.POST("/courses", courseHandler.CreateCourse)
	v1.GET("/courses", courseHandler.GetAllCourses)
	v1.GET("/courses/:id", courseHandler.GetCourse)
	v1.PUT("/courses/:id", courseHandler.UpdateCourse)
	v1.DELETE("/courses/:id", courseHandler.DeleteCourse)

	// Modules
	v1.POST("/modules", moduleHandler.CreateModule)
	v1.GET("/modules/:id", moduleHandler.GetModule)
	v1.PUT("/modules/:id", moduleHandler.UpdateModule)
	v1.DELETE("/modules/:id", moduleHandler.DeleteModule)
	v1.GET("/courses/:id/modules", moduleHandler.GetModulesByCourse)
	v1.POST("/courses/:id/modules/reorder", moduleHandler.ReorderModules)

	// Content
	v1.POST("/content", contentHandler.CreateContent)
	v1.GET("/content/:id", contentHandler.GetContent)
	v1.PUT("/content/:id", contentHandler.UpdateContent)
	v1.DELETE("/content/:id", contentHandler.DeleteContent)
	v1.GET("/modules/:id/content", contentHandler.GetContentsByModule)

	// Evaluations
	v1.POST("/evaluations", evaluationHandler.CreateEvaluation)
	v1.GET("/evaluations/:id", evaluationHandler.GetEvaluation)
	v1.PUT("/evaluations/:id", evaluationHandler.UpdateEvaluation)
	v1.DELETE("/evaluations/:id", evaluationHandler.DeleteEvaluation)
	v1.GET("/modules/:id/evaluations", evaluationHandler.GetEvaluationsByModule)

	// Enrollments
	v1.POST("/enrollments", enrollmentHandler.CreateEnrollment)
	v1.GET("/enrollments/:id", enrollmentHandler.GetEnrollment)
	v1.DELETE("/enrollments/:id", enrollmentHandler.DeleteEnrollment)
	v1.GET("/users/:userId/enrollments", enrollmentHandler.GetUserEnrollments)
	v1.GET("/courses/:id/enrollments", enrollmentHandler.GetCourseEnrollments)
	v1.POST("/users/:userId/courses/:id/complete", enrollmentHandler.CompleteEnrollment)
	v1.PUT("/users/:userId/courses/:id/progress", enrollmentHandler.UpdateProgress)

	// Evaluation Attempts
	v1.POST("/evaluation-attempts/start", evaluationAttemptHandler.StartAttempt)
	v1.POST("/evaluation-attempts/:id/submit", evaluationAttemptHandler.SubmitAttempt)
	v1.GET("/evaluation-attempts/:id", evaluationAttemptHandler.GetAttempt)
	v1.GET("/users/:userId/evaluations/:evaluationId/attempts", evaluationAttemptHandler.GetUserAttempts)
	v1.GET("/users/:userId/evaluations/:evaluationId/can-attempt", evaluationAttemptHandler.CanUserAttempt)
	v1.POST("/evaluation-attempts/:id/score", evaluationAttemptHandler.ScoreAttempt)
}

func (app *Application) registerDocs() {
	host := app.Config.Server.Host
	if utils.IsLocalhostURL(host) {
		host += ":" + app.Config.Server.Port
	}

	if utils.IsHttpsURL(host) {
		docs.SwaggerInfo.Schemes = []string{"https"}
	} else {
		docs.SwaggerInfo.Schemes = []string{"http"}
	}

	docs.SwaggerInfo.Host = utils.CleanHostURL(host)
	docs.SwaggerInfo.BasePath = "/"

	schemaUrl := host
	schemaUrl += "/internal/docs/doc.json"

	urlSwaggerJson := ginSwagger.URL(schemaUrl)
	app.Router.GET("/internal/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, urlSwaggerJson))
}

func (app *Application) Run() {
	addr := utils.CleanHostURL(":" + app.Config.Server.Port)
	app.Router.Run(addr)
}
