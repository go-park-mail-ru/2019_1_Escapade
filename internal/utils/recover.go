package utils

import "fmt"

//"fmt"

// CatchPanic stop pani to continue executing
func CatchPanic(place string) {
	// panic doesnt recover

	if r := recover(); r != nil {
		fmt.Println("Panic recovered in", place)
		fmt.Println("More", r)
	}
}
