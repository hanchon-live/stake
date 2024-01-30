package requester

import (
	"time"
)

func GenerateLast8Days() []string {
	currentTime := time.Now()
	return GeneratePrevious8Days(currentTime)
}

func GeneratePrevious8Days(now time.Time) []string {
	ret := []string{}
	offset := -7
	now = now.AddDate(0, 0, offset)
	i := 1
	for i <= 8 {
		ret = append(ret, now.Format("02-01-2006"))
		now = now.AddDate(0, 0, 1)
		i++
	}
	return ret
}
