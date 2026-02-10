package server

import (
	"go-boilerplate/internal/health"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-boilerplate/internal/auth"
	"go-boilerplate/internal/users"
)

func NewRouter(db *pgxpool.Pool, jwtSecret string) *gin.Engine {
	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Health check (no auth)
	r.GET("/health", health.Handler(db))

	// User handler
	userHandler := &users.Handler{
		DB:        db,
		JWTSecret: jwtSecret,
	}

	// Public auth routes
	r.POST("/users", userHandler.Register)
	r.POST("/login", userHandler.Login)

	// Protected routes
	authGroup := r.Group("/")
	authGroup.Use(auth.Middleware(jwtSecret))
	{
		authGroup.GET("/me", userHandler.Me)
	}

	return r
}
