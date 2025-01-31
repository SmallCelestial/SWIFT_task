package api

import (
	"SWIFT_task/internal/handler"
	"github.com/gin-gonic/gin"
	"log"
)

func SetupRouter(branchHandler *handler.BranchHandler) *gin.Engine {
	router := gin.Default()
	router.GET("/v1/swift-codes/:swift-code", branchHandler.GetBranchDetails)
	router.GET("/v1/swift-codes/country/:countryISO2code", branchHandler.GetBranchesByISO2code)
	router.POST("/v1/swift-codes/", branchHandler.AddSwiftCode)
	router.DELETE("/v1/swift-codes/:swift-code", branchHandler.RemoveSwiftCode)
	return router
}

func Run(branchHandler *handler.BranchHandler) {
	router := SetupRouter(branchHandler)
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
