package game

import (
	"encoding/json"
	"sync"
)

type Flag struct {
	Cell Cell `json:"cell"`
	Set  bool `json:"set"`
}

// OnlinePlayers online players
type OnlinePlayers struct {
	capacityM *sync.RWMutex
	_capacity int

	playersM *sync.RWMutex
	_players []Player `json:"players"`

	flagsM      *sync.RWMutex
	_flags      []Flag `json:"flags"`
	flagsLeft   int
	Connections Connections `json:"connections"`
}

type OnlinePlayersJSON struct {
	Capacity    int             `json:"capacity"`
	Players     []Player        `json:"players"`
	Connections ConnectionsJSON `json:"connections"`
	Flags       []Flag          `json:"flags"`
}

func (op *OnlinePlayers) JSON() OnlinePlayersJSON {
	return OnlinePlayersJSON{
		Capacity:    op.Capacity(),
		Players:     op.RPlayers(),
		Connections: op.Connections.JSON(),
		Flags:       op.Flags(),
	}
}

func (op *OnlinePlayers) MarshalJSON() ([]byte, error) {
	return json.Marshal(op.JSON())
}

func (op *OnlinePlayers) UnmarshalJSON(b []byte) error {
	temp := &OnlinePlayersJSON{}

	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	op.SetCapacity(temp.Capacity)
	op.SetPlayers(temp.Players)
	op.SetFlags(temp.Flags)

	return nil
}

// NewConnections create instance of Connections
func newOnlinePlayers(size int, field Field) *OnlinePlayers {
	players := make([]Player, size)
	flags := field.RandomFlags(players)
	return &OnlinePlayers{
		capacityM: &sync.RWMutex{},
		_capacity: size,

		playersM: &sync.RWMutex{},
		_players: players,

		flagsM:      &sync.RWMutex{},
		_flags:      flags,
		flagsLeft:   size,
		Connections: *NewConnections(size),
	}

}
