package internal

import (
	"regexp"
)

func GetTenderId(url string) string {

	re := regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)

	match := re.FindString(url)

	return match
}
