package engine

// OnlinePlayersJSON is a wrapper for sending OnlinePlayers by JSON
//easyjson:json
type OnlinePlayersJSON struct {
	Capacity    int32           `json:"capacity"`
	Players     []Player        `json:"players"`
	Connections ConnectionsJSON `json:"connections"`
	Flags       []Flag          `json:"flags"`
}

// JSON convert OnlinePlayers to OnlinePlayersJSON
func (op *OnlinePlayers) JSON() OnlinePlayersJSON {
	return OnlinePlayersJSON{
		Capacity:    op.m.Capacity(),
		Players:     op.m.CopyPlayers(),
		Connections: op.Connections.JSON(),
		Flags:       op.m.Flags(),
	}
}

// MarshalJSON - overriding the standard method json.Marshal
func (op *OnlinePlayers) MarshalJSON() ([]byte, error) {
	return op.JSON().MarshalJSON()
}

// UnmarshalJSON - overriding the standard method json.Unmarshal
func (op *OnlinePlayers) UnmarshalJSON(b []byte) error {
	temp := &OnlinePlayersJSON{}

	if err := temp.UnmarshalJSON(b); err != nil {
		return err
	}

	op.m.SetCapacity(temp.Capacity)
	op.m.SetPlayers(temp.Players)
	op.m.SetFlags(temp.Flags)

	return nil
}
