package internal

import "fmt"

func ValidateStatus(tenderStatus []string, validTenderStatus map[string]bool) error {
	for _, tenderStatus := range tenderStatus {
		if _, exists := validTenderStatus[tenderStatus]; !exists {
			return fmt.Errorf("Invalid tender status  provided. Available values: Created, Published, Closed")
		}
	}
	return nil
}
