package game

// People send to user, if he disconnect and 'forgot' everything
// about users or it is his first connect
type People struct {
	Players   *Connections
	Observers *Connections
}
