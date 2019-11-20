package engine

import (
	"sync"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
)

// RoomInformationI contains the meta information about room
// Memento pattern
type RoomInformationI interface {
	setName(name string)
	Name() string

	setID(id string)
	ID() string

	Settings() *models.RoomSettings
	setSettings(rs *models.RoomSettings)

	RoomID() int32
}

// RoomInformation implement RoomInformationI
type RoomInformation struct {
	dbRoomID int32

	idM *sync.RWMutex
	_id string

	nameM *sync.RWMutex
	_name string

	settings *models.RoomSettings
}

// Init set Room and RoomNotifier pointers
func (room *RoomInformation) Init(rs *models.RoomSettings, id string,
	dbRoomID int32) {
	room.nameM = &sync.RWMutex{}
	room._name = rs.Name

	room.settings = rs

	room.idM = &sync.RWMutex{}
	room._id = id

	room.dbRoomID = dbRoomID

	room.setID(utils.RandomString(16))
}

func (room *RoomInformation) Settings() *models.RoomSettings {
	return room.settings
}

func (room *RoomInformation) setSettings(rs *models.RoomSettings) {
	room.settings = rs
}

func (room *RoomInformation) RoomID() int32 {
	return room.dbRoomID
}

////////////////////////// mutex

// Name return the name of room given by its creator
func (room *RoomInformation) Name() string {
	room.nameM.RLock()
	v := room._name
	room.nameM.RUnlock()
	return v
}

// ID return room's unique identificator
func (room *RoomInformation) ID() string {
	room.idM.RLock()
	v := room._id
	room.idM.RUnlock()
	return v
}

func (room *RoomInformation) setName(name string) {
	room.nameM.Lock()
	room._name = name
	room.nameM.Unlock()
}

func (room *RoomInformation) setID(id string) {
	room.idM.Lock()
	room._id = id
	room.idM.Unlock()
}
