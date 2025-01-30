package main

import (
	"SWIFT_task/api"
	"SWIFT_task/internal/db"
)

func main() {
	db.InitDB()
	api.Run()
}
