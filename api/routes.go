package api

import (
	"SWIFT_task/internal/handler"
	"github.com/gin-gonic/gin"
	"log"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/v1/swift-codes/:swift-code", handler.GetBranchDetails)
	return router
}

func Run() {
	router := SetupRouter()
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
