package main

import (
    "log"
    "os"
    "spotsync-api/handler"
    "spotsync-api/middleware"
    "spotsync-api/models"
    "spotsync-api/repository"
    "spotsync-api/service"

    "github.com/go-playground/validator/v10"
    "github.com/joho/godotenv"
    "github.com/labstack/echo/v4"
    echoMiddleware "github.com/labstack/echo/v4/middleware"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type CustomValidator struct {
    validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
    return cv.validator.Struct(i)
}

func main() {
    // Load .env
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: .env file not found")
    }

    // Database connection
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        dsn = "host=localhost user=postgres password=postgres dbname=spotsync port=5432 sslmode=disable"
    }

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Auto migrate
    if err := db.AutoMigrate(&models.User{}, &models.ParkingZone{}, &models.Reservation{}); err != nil {
        log.Fatal("Failed to migrate database:", err)
    }

    // Initialize repositories
    userRepo := repository.NewUserRepository(db)
    zoneRepo := repository.NewZoneRepository(db)
    reservationRepo := repository.NewReservationRepository(db)

    // Initialize services
    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        jwtSecret = "your-secret-key-change-in-production"
        log.Println("Warning: JWT_SECRET not set, using default")
    }

    authService := service.NewAuthService(userRepo, jwtSecret)
    zoneService := service.NewZoneService(zoneRepo)
    reservationService := service.NewReservationService(db, reservationRepo, zoneRepo)

    // Initialize handlers
    authHandler := handler.NewAuthHandler(authService)
    zoneHandler := handler.NewZoneHandler(zoneService)
    reservationHandler := handler.NewReservationHandler(reservationService)

    // Initialize middleware
    authMiddleware := middleware.NewAuthMiddleware(authService)

    // Echo setup
    e := echo.New()

    // Custom validator
    e.Validator = &CustomValidator{validator: validator.New()}

    // Middleware
    e.Use(echoMiddleware.Logger())
    e.Use(echoMiddleware.Recover())
    e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
        AllowOrigins: []string{"*"},
        AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
        AllowHeaders: []string{echo.HeaderAuthorization, echo.HeaderContentType},
    }))

    // Health check
    e.GET("/health", func(c echo.Context) error {
        return c.JSON(200, map[string]string{"status": "ok"})
    })

    // API Routes
    api := e.Group("/api/v1")

    // Auth routes (Public)
    auth := api.Group("/auth")
    auth.POST("/register", authHandler.Register)
    auth.POST("/login", authHandler.Login)

    // Zone routes - Public GET, Admin POST
    zones := api.Group("/zones")
    zones.GET("", zoneHandler.GetAllZones)
    zones.GET("/:id", zoneHandler.GetZoneByID)
    
    // Admin-only zone creation - FIXED: Handler first, then middleware
    zones.POST("", zoneHandler.CreateZone, authMiddleware.Authenticate, authMiddleware.RequireAdmin)

    // Reservation routes (Authenticated)
    reservations := api.Group("/reservations")
    reservations.Use(authMiddleware.Authenticate) // Apply middleware to all routes in group
    reservations.POST("", reservationHandler.CreateReservation)
    reservations.GET("/my-reservations", reservationHandler.GetMyReservations)
    reservations.DELETE("/:id", reservationHandler.CancelReservation)

    // Admin-only reservations routes
    adminReservations := api.Group("/reservations")
    adminReservations.Use(authMiddleware.Authenticate, authMiddleware.RequireAdmin)
    adminReservations.GET("", reservationHandler.GetAllReservations)

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on :%s", port)
    if err := e.Start(":" + port); err != nil {
        log.Fatal(err)
    }
}