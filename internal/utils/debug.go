package utils

import "fmt"

// ShowWebsocketMessage record information transmitted over the websocket
func ShowWebsocketMessage(message []byte, id int32) {
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
		Debug(false, "#", id, " get that:", print)
	}
}

func PrintResult(catched error, number int, place string) {
	if catched != nil {
		Debug(false, "api/"+place+" failed(code:", number, "). Error message:"+catched.Error())
	} else {
		Debug(false, "api/"+place+" success(code:", number, ")")
	}
}

// Debug record information for logging
// if 'needPanic' is true then panic(text) will be called
func Debug(needPanic bool, text ...interface{}) {
	fmt.Print("    ")
	fmt.Println(text...)
	if needPanic {
		panic("panic called!")
	}
}
