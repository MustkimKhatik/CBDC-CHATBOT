package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// Health is a simple health check endpoint
func Health(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status":  "ok",
        "service": "cbdc-backend",
    })
}


