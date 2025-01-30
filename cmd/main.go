package main

import (
	"SWIFT_task/internal/db"
	"SWIFT_task/pkg"
	"log"
)

func main() {

	db.InitDB()

	allBranches, err := pkg.ParseSwiftCsvToBranches("data/swift-codes.csv")
	if err != nil {
		log.Fatalf("Błąd podczas parsowania CSV: %v", err)
	}

	db.SaveBranches(allBranches)
}
