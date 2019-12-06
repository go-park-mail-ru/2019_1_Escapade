package constants

import "fmt"

import "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"

// ErrorFieldWidth is called when field width is invalid
func ErrorFieldWidth(rs *models.RoomSettings) error {
	return fmt.Errorf("Field width is invalid:%d."+
		"Need more then %d and less then %d", rs.Width,
		FIELD.WidthMin, FIELD.WidthMax)
}

// ErrorFieldHeight is called when field height is invalid
func ErrorFieldHeight(rs *models.RoomSettings) error {
	return fmt.Errorf("Field height is invalid:%d."+
		"Need more then %d and less then %d", rs.Height,
		FIELD.HeightMin, FIELD.HeightMax)
}

// ErrorRoomName is called when room name is invalid
func ErrorRoomName(rs *models.RoomSettings) error {
	return fmt.Errorf("Name's '%s' length is invalid:%d."+
		"Need more then %d and less then %d", rs.Name,
		len(rs.Name), ROOM.NameMin, ROOM.NameMax)
}

// ErrorPlayers is called when amount of players is invalid
func ErrorPlayers(rs *models.RoomSettings) error {
	return fmt.Errorf("Players amount is invalid:%d."+
		"Need more then %d and less then %d. Also "+
		"amount should be less or equal then width*height(%d)", rs.Players,
		ROOM.PlayersMin, ROOM.PlayersMax, rs.Width*rs.Height)
}

// ErrorObservers is called when amount of observers is invalid
func ErrorObservers(rs *models.RoomSettings) error {
	return fmt.Errorf("Observers amount is invalid:%d."+
		"Need equal or more then 0 and less then %d",
		rs.Observers, ROOM.ObserversMax)
}

// ErrorTimeToPrepare is called when time to prepare is invalid
func ErrorTimeToPrepare(rs *models.RoomSettings) error {
	return fmt.Errorf("Time to prepare is invalid:%d."+
		"Need more then %d and less then %d", rs.TimeToPrepare,
		ROOM.TimeToPrepareMin, ROOM.TimeToPrepareMax)
}

// ErrorTimeToPlay is called when time to play is invalid
func ErrorTimeToPlay(rs *models.RoomSettings) error {
	return fmt.Errorf("Time to play is invalid:%d."+
		"Need more then %d and less then %d", rs.TimeToPlay,
		ROOM.TimeToPlayMin, ROOM.TimeToPlayMax)
}

// ErrorConstantsNotSet is called when constants havent set
func ErrorConstantsNotSet() error {
	return fmt.Errorf("Constants not set")
}

// ErrorMines is called when amount of mines is invalid
func ErrorMines(mines, max int32) error {
	return fmt.Errorf("Mines amount is invalid:%d."+
		"It should be more then 0 and less then"+
		"'width*height-players_amount'(%d)", mines, max)
}
