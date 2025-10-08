package router

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"

    "cbdc-backend/handlers"
)

// SetupRouter configures and returns the Gin engine
func SetupRouter() *gin.Engine {
    r := gin.Default()
    r.Use(cors.Default())

    api := r.Group("/api")
    {
        api.GET("/health", handlers.Health)
        api.POST("/upload", handlers.Upload)
        api.POST("/query", handlers.Query)
    }

    return r
}


