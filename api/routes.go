package api

import (
	"SWIFT_task/internal/handler"
	"github.com/gin-gonic/gin"
	"log"
)

func SetupRouter(branchHandler *handler.BankHandler) *gin.Engine {
	router := gin.Default()
	router.GET("/v1/swift-codes/:swift-code", branchHandler.GetBankDetails)
	router.GET("/v1/swift-codes/country/:countryISO2code", branchHandler.GetBanksByISO2code)
	router.POST("/v1/swift-codes/", branchHandler.AddBank)
	router.DELETE("/v1/swift-codes/:swift-code", branchHandler.RemoveBank)
	return router
}

func Run(bankHandler *handler.BankHandler) {
	router := SetupRouter(bankHandler)
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
