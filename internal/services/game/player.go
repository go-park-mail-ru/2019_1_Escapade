package game

import "fmt"

type Player struct {
	Name string
	ID   int
}

func NewPlayer(name string, id int) *Player {
	player := &Player{
		Name: name,
		ID:   id,
	}
	return player
}

func (p *Player) Command(command string) {
	fmt.Println("Command: '", command, "' received by player: ", p.Name)
}

func (p *Player) GetState() string {
	return "Game state for Player: " + p.Name
}

func (p *Player) GiveUp() {
	fmt.Println("Player gave up: ", p.Name)
}
