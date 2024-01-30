package requester

import (
	"fmt"
	"testing"
	"time"
)

func TestPrevious8Days(t *testing.T) {
	now, err := time.Parse("02-01-2006", "24-02-2023")
	if err != nil {
		t.Fatalf("could not parse time to generate dates")
	}
	dates := GeneratePrevious8Days(now)
	base := 17
	i := 0
	for i < 8 {
		generated := fmt.Sprintf("%d-02-2023", (base + i))
		if dates[i] != generated {
			t.Fatalf("incorrect date %s != %s", dates[i], generated)
		}
		i++
	}
}
