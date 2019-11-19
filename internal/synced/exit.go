package synced

// https://stackoverflow.com/questions/27629380/how-to-exit-a-go-program-honoring-deferred-calls

import "os"

// Exit passed to the func panic() to cause the program to
// terminate(os.Exit) with the specified code and save the execution
// of defered funcs. Before use, place defer HandleExit()
type Exit struct{ Code int }

// HandleExit handle panic() where is the Exit structure passed
func HandleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(Exit); ok == true {
			os.Exit(exit.Code)
		}
		panic(e) // not an Exit, bubble up
	}
}
