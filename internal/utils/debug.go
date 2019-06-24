package utils

import (
	"fmt"
)

// PrintResult log requests results
func ShowWebsocketMessage(message []byte, id int) {
	str := string(message)
	var start, end, counter int
	for i, s := range str {
		if s == '"' {
			counter++
			if counter == 3 {
				start = i + 1
			} else if counter == 4 {
				end = i
			} else if counter > 4 {
				break
			}
		}
	}
	if start != end {
		print := str[start:end]
		//print = str
		fmt.Println("#", id, " get that:", print)
	}
}

func Debug(needPanic bool, text ...interface{}) {
	if needPanic {
		panic(text)
	}
	fmt.Println(text...)
}
