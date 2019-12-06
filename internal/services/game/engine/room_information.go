package engine

import (
	"sync"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
	room_ "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine/room"
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

	PlayingTime() time.Duration
	RecruitmentTime() time.Duration

	Date() time.Time
	SetDate(date time.Time)

	RoomID() int32
}

// RoomInformation implement RoomInformationI
type RoomInformation struct {
	dbRoomID int32

	idM *sync.RWMutex
	_id string

	nameM *sync.RWMutex
	_name string

	dateM *sync.RWMutex
	_date time.Time

	recruitmentTimeM *sync.RWMutex
	_recruitmentTime time.Duration

	playingTimeM *sync.RWMutex
	_playingTime time.Duration

	location *time.Location

	settingsM *sync.RWMutex
	_settings *models.RoomSettings
}

func (room *RoomInformation) init(settings *models.RoomSettings,
	id string, dbRoomID int32, location *time.Location) {

	room.dbRoomID = dbRoomID

	room.idM = &sync.RWMutex{}
	if id == "" {
		id = utils.RandomString(16)
	}
	room.setID(id)

	room.nameM = &sync.RWMutex{}
	room.setName(settings.Name)

	room.dateM = &sync.RWMutex{}
	room.SetDate(time.Now().In(location))

	room.recruitmentTimeM = &sync.RWMutex{}
	room._recruitmentTime = 0

	room.playingTimeM = &sync.RWMutex{}
	room._playingTime = 0

	room.location = location

	room.settingsM = &sync.RWMutex{}
	room.setSettings(settings)
}

func (room *RoomInformation) subscribe(builder RBuilderI) {
	var events EventsI
	builder.BuildEvents(&events)
	room.eventsSubscribe(events)
}

// Init set Room and RoomNotifier pointers
func (room *RoomInformation) Init(rs *models.RoomSettings, id string,
	dbRoomID int32, location *time.Location) {

	room.setID(utils.RandomString(16))
}

func (room *RoomInformation) Settings() *models.RoomSettings {
	room.settingsM.RLock()
	defer room.settingsM.RUnlock()
	return room._settings
}

func (room *RoomInformation) setSettings(rs *models.RoomSettings) {
	room.settingsM.Lock()
	defer room.settingsM.Unlock()
	room._settings = rs
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

func (room *RoomInformation) PlayingTime() time.Duration {
	room.playingTimeM.RLock()
	v := room._playingTime
	room.playingTimeM.RUnlock()
	return v
}

func (room *RoomInformation) setPlayingTime() {
	room.dateM.RLock()
	v := room._date
	room.dateM.RUnlock()

	t := time.Now().In(room.location)

	room.playingTimeM.Lock()
	room._playingTime = t.Sub(v)
	room.playingTimeM.Unlock()

	room.dateM.Lock()
	room._date = t
	room.dateM.Unlock()
}

// Date return date, when room was created
func (room *RoomInformation) Date() time.Time {
	room.dateM.RLock()
	v := room._date
	room.dateM.RUnlock()
	return v
}

func (room *RoomInformation) SetDate(date time.Time) {
	room.dateM.Lock()
	room._date = date
	room.dateM.Unlock()
}

func (room *RoomInformation) RecruitmentTime() time.Duration {
	room.recruitmentTimeM.RLock()
	v := room._recruitmentTime
	room.recruitmentTimeM.RUnlock()
	return v
}

func (room *RoomInformation) setRecruitmentTime() {
	room.dateM.RLock()
	v := room._date
	room.dateM.RUnlock()

	t := time.Now().In(room.location)

	room.recruitmentTimeM.Lock()
	room._recruitmentTime = t.Sub(v)
	room.recruitmentTimeM.Unlock()

	room.dateM.Lock()
	room._date = t
	room.dateM.Unlock()
}

func (room *RoomInformation) eventsRunning(synced.Msg) {
	room.setRecruitmentTime()
}

func (room *RoomInformation) eventsFinished(synced.Msg) {
	room.setPlayingTime()
}

func (room *RoomInformation) eventsSubscribe(events EventsI) {
	observer := synced.NewObserver(
		synced.NewPair(room_.StatusRunning, room.eventsRunning),
		synced.NewPair(room_.StatusFinished, room.eventsFinished))
	events.Observe(observer.AddPublisherCode(room_.UpdateStatus))
}
