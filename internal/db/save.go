package db

import (
	"SWIFT_task/internal/model"
	"fmt"
	"gorm.io/gorm"
	"log"
)

func SaveBranches(branches map[string]*model.Bank, db *gorm.DB) {
	for _, branch := range branches {
		result := db.Create(branch)
		if result.Error != nil {
			log.Printf("Error saving branches %s", result.Error.Error())
		}
	}

	for _, branch := range branches {
		if !branch.IsHeadquarter {
			headquarter, ok := branches[branch.GetHeadQuarterSwiftCode()]
			if !ok {
				fmt.Printf("Bank: %+v is not headquarter, but also hasn't headquarter\n", branch)
			} else {
				assignBranches := &model.BankRelationship{
					HeadquarterSwiftCode: headquarter.SwiftCode,
					BranchSwiftCode:      branch.SwiftCode,
					Headquarter:          headquarter,
					Branch:               branch,
				}

				result := db.Create(assignBranches)
				if result.Error != nil {
					log.Printf("Error saving branches %s", result.Error.Error())
				}
			}
		}
	}
}
