package main

import (
	"SWIFT_task/pkg"
	"fmt"
	"log"
)

func main() {
	allBranches, err := pkg.ParseSwiftCsvToBranches("data/swift-codes.csv")
	if err != nil {
		log.Fatalf("Błąd podczas parsowania CSV: %v", err)
	}

	for swiftCode, branch := range allBranches {
		fmt.Printf("SWIFT Kod: %s\n", swiftCode)
		fmt.Printf("%+v\n", branch)
		fmt.Println("---------")
	}
}
