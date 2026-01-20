package main

import (
	"fmt"
	"log"
	"os"

	"backend/config"
	_ "backend/docs"
	"backend/handlers"
	"backend/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           ITIDoorz Room Booking API
// @version         1.0
// @description     API Sistem Pemesanan Ruang Kelas dan Lab - ITI.
// @termsOfService  http://swagger.io/terms/

// @contact.name    Engineering Team
// @contact.email   admin@iti.ac.id

// @host            localhost:8080
// @BasePath        /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 2. Database Connection
	// parseTime=true is required to map MySQL DATETIME to Go time.Time
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", 
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Database unresponsive:", err)
	}
	
	// Menggunakan log.Println agar ada timestamp di terminal
	log.Println("Database connection established successfully.")

	// 3. Initialize Gin Engine
	r := gin.Default()

	// 4. CORS Configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true // Allow all for development
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept"} 
	
	r.Use(cors.New(corsConfig))

	// 5. Static Files Setup
	if _, err := os.Stat("storage"); os.IsNotExist(err) {
		_ = os.Mkdir("storage", 0755)
	}
	r.Static("/storage", "./storage")

	// 6. Swagger Documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 7. Dependency Injection (Repositories)
	roomRepo := repository.NewRoomRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	userRepo := repository.NewUserRepository(db)

	// 8. Initialize Handlers
	roomHandler := handlers.RoomHandler{Repo: roomRepo}
	bookingHandler := handlers.BookingHandler{Repo: bookingRepo}
	authHandler := handlers.AuthHandler{UserRepo: userRepo}
	userHandler := handlers.UserHandler{Repo: userRepo}
	reportHandler := handlers.ReportHandler{
		BookingRepo: bookingRepo,
		RoomRepo:    roomRepo,
	}

	// 9. API Routes
	api := r.Group("/api")
	{
		// Public Routes
		api.POST("/login", authHandler.Login)
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "success", "message": "Backend ITIDoorz is Live!"})
		})

		// Users
		api.GET("/users", userHandler.GetAllUsers)
		api.POST("/users", userHandler.CreateUser)
		api.PUT("/users/:id", userHandler.UpdateUser)
		api.DELETE("/users/:id", userHandler.DeleteUser)

		// Rooms
		api.GET("/rooms", roomHandler.GetRooms)
		api.POST("/rooms", roomHandler.CreateRoom)
		api.PUT("/rooms/:id", roomHandler.UpdateRoom)
		api.DELETE("/rooms/:id", roomHandler.DeleteRoom)

		// Bookings
		api.GET("/bookings", bookingHandler.GetBookings)
		api.POST("/bookings", bookingHandler.CreateBooking)
		api.POST("/bookings/:id/approve", bookingHandler.ApproveBooking)
		api.POST("/bookings/:id/reject", bookingHandler.RejectBooking)

		// Stats & Reports
		api.GET("/stats", bookingHandler.GetStats)
		api.GET("/reports/generate", reportHandler.GenerateMonthlyReport)
	}

	// 10. Start Server
	log.Println("Server running on port :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}