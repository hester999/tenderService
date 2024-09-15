package internal

import (
	"strconv"
)

func IntIsLast(url string) (int, bool) {
	if len(url) == 0 {
		return 0, false
	}

	res := url[len(url)-1:]

	digit, err := strconv.Atoi(res)
	if err != nil {
		return 0, false
	}

	return digit, true
}
