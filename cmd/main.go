package main

import (
	//"SWIFT_task/api"
	"SWIFT_task/internal/db"
	//"SWIFT_task/internal/handler"
	//"SWIFT_task/internal/repository"
	//"SWIFT_task/internal/service"
)

func main() {
	database := db.InitDB()
	db.SaveData("data/swift-codes.csv", database)

	//branchRepo := repository.NewBranchRepository(database)
	//branchService := service.NewBranchService(branchRepo)
	//branchHandler := handler.NewBranchHandler(branchService)
	//
	//api.Run(branchHandler)
}
