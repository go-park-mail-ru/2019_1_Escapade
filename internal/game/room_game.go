package game

import (
	"fmt"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
)

// Winners determine who won the game
func (room *Room) Winners() (winners []int) {
	if room.done() {
		return nil
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	max := 0.

	players := room.Players.RPlayers()
	for _, player := range players {
		if player.Points > max && !player.Died {
			max = player.Points
		}
	}

	winners = make([]int, 0)
	for index, player := range players {
		if player.Points == max && !player.Died {
			max = player.Points
			winners = append(winners, index)
		}
	}

	return
}

// Winner check id in the winners slice
func (room *Room) Winner(winners []int, find int) bool {
	for i := range winners {
		if find == i {
			return true
		}
	}
	return false
}

// FlagFound is called, when somebody find cell flag
func (room *Room) FlagFound(founder Connection, found *Cell) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	which := 0
	for _, flag := range room.Players.Flags() {
		if flag.Cell.X == found.X && flag.Cell.Y == found.Y {
			which = flag.Cell.PlayerID
		}
	}

	if which == founder.ID() {
		return
	}

	room.Players.IncreasePlayerPoints(founder.Index(), 300)

	killConn, index := room.Players.Connections.SearchByID(which)
	fmt.Println(killConn.User.Name, "was found by", founder.User.Name)
	if index >= 0 {
		room.Kill(killConn, ActionFlagLost)
	}
}

// isAlive check if connection is player and he is not died
func (room *Room) isAlive(conn *Connection) bool {
	index := conn.Index()
	return index >= 0 && !room.Players.Player(index).Finished
}

// Kill make user die and check for finish battle
func (room *Room) Kill(conn *Connection, action int) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	fmt.Println("Status", room.Status)
	// cause all in pointers
	if room.Status < StatusFlagPlacing && room.Status > StatusRunning {
		return
	}

	if room.isAlive(conn) {
		room.SetFinished(conn)

		cell := room.Players.Flag(conn.Index())
		//room.Field.SetCellFlagTaken(&cell.Cell)

		cells := make([]Cell, 0)
		room.Field.saveCell(&cell.Cell, &cells)

		room.sendNewCells(room.All, cell.Cell)

		if room.Players.Capacity() <= room.killed()+1 {
			room.chanStatus <- StatusFinished
			//room.FinishGame(false)
		}
		pa := *room.addAction(conn.ID(), action)
		go room.sendAction(pa, room.All)
	}
	return
}

// GiveUp kill connection, that call it
func (room *Room) GiveUp(conn *Connection) {
	room.Kill(conn, ActionGiveUp)
}

// flagExists find players with such flag. This - flag owner
func (room *Room) flagExists(cell Cell, this *Connection) (found bool, conn *Connection) {
	var player int
	flags := room.Players.Flags()
	for index, flag := range flags {
		if (flag.Cell.X == cell.X) && (flag.Cell.Y == cell.Y) {
			if this == nil || index != this.Index() {
				found = true
				player = index
			}
			break
		}
	}
	if !found {
		return
	}
	conn = room.Players.Connections.SearchByIndex(player)
	return
}

// SetAndSendNewCell set and send cell to conn
func (room *Room) SetAndSendNewCell(conn Connection) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	found := true
	// create until it become unique
	var cell Cell
	for found {
		cell = room.Field.CreateRandomFlag(conn.ID())
		found, _ = room.flagExists(cell, nil)
	}
	if room.Players.SetFlag(conn, cell) {
		room.chanStatus <- StatusRunning
	}
	response := models.RandomFlagSet(cell)
	conn.SendInformation(response)
}

// SetFlag handle user want set flag
func (room *Room) SetFlag(conn *Connection, cell *Cell) bool {
	if room.done() {
		return false
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	// if user try set flag after game launch
	if room.Status != StatusFlagPlacing {
		response := models.FailFlagSet(cell, re.ErrorBattleAlreadyBegan())
		conn.SendInformation(response)
		return false
	}

	if !room.Field.IsInside(cell) {
		response := models.FailFlagSet(cell, re.ErrorCellOutside())
		conn.SendInformation(response)
		return false
	}

	if !room.isAlive(conn) {
		response := models.FailFlagSet(cell, re.ErrorPlayerFinished())
		conn.SendInformation(response)
		return false
	}

	if found, prevConn := room.flagExists(*cell, conn); found {
		pa := *room.addAction(conn.ID(), ActionFlagСonflict)
		room.sendAction(pa, room.All)
		go room.SetAndSendNewCell(*conn)
		go room.SetAndSendNewCell(*prevConn)
		return true
	}

	pa := *room.addAction(conn.ID(), ActionFlagSet)
	room.sendAction(pa, room.All)

	if room.Players.SetFlag(*conn, *cell) {
		room.chanStatus <- StatusRunning
	}
	return true
}

// setFlags set players flags to field
// call it if game has already begun
func (room *Room) setFlags() {
	flags := room.Players.Flags()
	for _, cell := range flags {
		room.Field.SetFlag(&cell.Cell)
	}
}

// FillField set flags and mines
func (room *Room) FillField() {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	room.Field.Zero()
	room.setFlags()
	room.Field.SetMines()

}

// addAction creates an action and passes it on appendAction()
func (room *Room) addAction(id int, action int) (pa *PlayerAction) {
	if room.done() {
		return
	}
	room.wGroup.Add(1)
	defer func() {
		room.wGroup.Done()
	}()

	pa = NewPlayerAction(id, action)
	room.appendAction(pa)
	return
}
