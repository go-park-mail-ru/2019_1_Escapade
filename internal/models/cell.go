package models

// send to user one cell
type Cell struct {
	X int
	Y int
	// IsOpen bool
	// IsMarked bool
	Value int
}

/*
0-8- amount of mines
9 - mine
10 - flag
11 - almost open
*/

// RoomJoining joining to room
type RoomJoining struct {
	Success bool
	Message string
}

func RoomJoiningFault(reason string) *RoomJoining {
	rj := &RoomJoining{
		Success: false,
		Message: reason,
	}
	return rj
}

func RoomJoiningSuccess() *RoomJoining {
	rj := &RoomJoining{
		Success: true,
	}
	return rj
}

type Status struct {
	Ready       bool
	PeopleFound int
}

func NewStatus(peopleFound int, ready bool) *Status {
	status := &Status{
		Ready:       ready,
		PeopleFound: peopleFound,
	}
	return status
}

type Points struct {
	points int
}

type Finish struct {
	finish bool
	win    bool
}

type GameInfo struct {
	// show, what of fields is filled
	who int
	//waiting for players. There messages send before game starts
	Status Status
	//cell information. Send during game for all players
	cell Cell
	//Player points. Individually for every player
	points Points
	// Results. Send at the end of game
	finish Finish
}

type Flag struct {
	ID   int
	Cell *Cell
}

func NewFlag(cell *Cell, id int) *Flag {
	flag := &Flag{
		ID:   id,
		Cell: cell,
	}
	return flag
}
