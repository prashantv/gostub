package gostub_test

import (
	"fmt"
	"time"

	"github.com/prashantv/gostub"
)

// Production code
var timeNow = time.Now

func GetDay() int {
	return timeNow().Day()
}

// Test code
func Example_stubFunctions() {
	stubs := gostub.Stub(&timeNow, func() time.Time {
		return time.Date(2015, 07, 02, 0, 0, 0, 0, time.UTC)
	})
	defer stubs.Reset()

	fmt.Println("Stubbed:", GetDay())
	// Output:
	// Stubbed: 2
}
