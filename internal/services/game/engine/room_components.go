package engine

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// ComponentBuilderI create room components, bind them, and add them to room
// 	to the room
// ABuilder Pattern
type ComponentBuilderI interface {
	Build(r *Room, field *Field, lobby *Lobby,
		rs *models.RoomSettings, timeouts Timeouts, id string, chatID,
		roomID int32)

	BuildInformation(i *RoomInformationI)
	BuildField(f *FieldProxyI)
	BuildSync(s *SyncI)
	BuildAPI(a *APIStrategyI)
	BuildLobby(l *LobbyProxyI)
	BuildModelsAdapter(m *ModelsAdapterI)
	BuildSender(s *SendStrategyI)
	BuildRoomConnectionEvents(c *ConnectionEventsI)
	BuildPeople(p *PeopleI)
	BuildEvents(e *EventsI)
	BuildRecorder(r *ActionRecorderProxyI)
	BuildMetrics(m *MetricsStrategyI)
	BuildMessages(m *MessagesProxyI)
	BuildGarbageCollector(g *GarbageCollectorI)
}

// RoomBuilder implements ComponentBuilderI
type RoomBuilder struct {
	info             *RoomInformation
	field            *RoomField
	sync             *SyncWgroup
	api              *RoomAPI
	lobby            *RoomLobby
	models           *RoomModelsAdapter
	sender           *RoomSender
	connEvents       *RoomConnectionEvents
	people           *RoomPeople
	events           *RoomEvents
	record           *RoomRecorder
	metrics          *RoomMetrics
	messages         *RoomMessages
	garbageCollector *RoomGarbageCollector
}

// Build create Room components and set them to Room object
func (room *RoomBuilder) Build(r *Room, field *Field, lobby *Lobby,
	rs *models.RoomSettings, timeouts Timeouts, id string, chatID,
	roomID int32) {

	room.createComponents()
	room.configureDependencies(r, field, lobby, rs,
		timeouts, id, chatID, roomID)
}

func (room *RoomBuilder) createComponents() {
	room.sync = &SyncWgroup{}
	room.info = &RoomInformation{}
	room.api = &RoomAPI{}
	room.lobby = &RoomLobby{}
	room.field = &RoomField{}
	room.models = &RoomModelsAdapter{}
	room.sender = &RoomSender{}
	room.people = &RoomPeople{}
	room.connEvents = &RoomConnectionEvents{}
	room.events = &RoomEvents{}
	room.metrics = &RoomMetrics{}
	room.record = &RoomRecorder{}
	room.messages = &RoomMessages{}
	room.garbageCollector = &RoomGarbageCollector{}
}

func (room *RoomBuilder) configureDependencies(r *Room,
	field *Field, lobby *Lobby, rs *models.RoomSettings,
	timeouts Timeouts, id string, chatID, roomID int32) {
	room.sync.Init()
	room.info.Init(rs, id, roomID)
	room.api.Init(room)
	room.lobby.Init(room, r, lobby)
	room.field.Init(room, field, rs.Deathmatch)
	room.models.Init(room)
	room.sender.Init(room)
	room.people.Init(room, rs)
	room.connEvents.Init(room, rs.Deathmatch)
	room.events.Init(room, rs)
	room.metrics.Init(room, rs)
	room.record.Init(room)
	room.messages.Init(room, chatID)
	room.garbageCollector.Init(room, timeouts)
}

func (room *RoomBuilder) set(build *Room) {
	build.sync = room.sync
	build.info = room.info
	build.api = room.api
	build.lobby = room.lobby
	build.field = room.field
	build.models = room.models
	build.sender = room.sender
	build.people = room.people
	build.connEvents = room.connEvents
	build.events = room.events
	build.metrics = room.metrics
	build.record = room.record
	build.messages = room.messages
	build.garbageCollector = room.garbageCollector
}

// BuildInformation set RoomInformationI implementation
func (room *RoomBuilder) BuildInformation(i *RoomInformationI) {
	*i = room.info
}

// BuildField set FieldProxyI implementation
func (room *RoomBuilder) BuildField(f *FieldProxyI) {
	*f = room.field
}

// BuildSync set SyncI implementation
func (room *RoomBuilder) BuildSync(s *SyncI) {
	*s = room.sync
}

// BuildAPI set APIStrategyI implementation
func (room *RoomBuilder) BuildAPI(a *APIStrategyI) {
	*a = room.api
}

// BuildLobby set APIStrategyI implementation
func (room *RoomBuilder) BuildLobby(l *LobbyProxyI) {
	*l = room.lobby
}

// BuildModelsAdapter set ModelsAdapterI implementation
func (room *RoomBuilder) BuildModelsAdapter(m *ModelsAdapterI) {
	*m = room.models
}

// BuildSender set SendStrategyI implementation
func (room *RoomBuilder) BuildSender(s *SendStrategyI) {
	*s = room.sender
}

// BuildRoomConnectionEvents set ConnectionEventsI implementation
func (room *RoomBuilder) BuildRoomConnectionEvents(c *ConnectionEventsI) {
	*c = room.connEvents
}

// BuildPeople set PeopleI implementation
func (room *RoomBuilder) BuildPeople(p *PeopleI) {
	*p = room.people
}

// BuildEvents set EventsI implementation
func (room *RoomBuilder) BuildEvents(e *EventsI) {
	*e = room.events
}

// BuildRecorder set ActionRecorderProxyI implementation
func (room *RoomBuilder) BuildRecorder(r *ActionRecorderProxyI) {
	*r = room.record
}

// BuildMetrics set MetricsStrategyI implementation
func (room *RoomBuilder) BuildMetrics(m *MetricsStrategyI) {
	*m = room.metrics
}

// BuildMessages set MessagesProxyI implementation
func (room *RoomBuilder) BuildMessages(m *MessagesProxyI) {
	*m = room.messages
}

// BuildGarbageCollector set GarbageCollectorI implementation
func (room *RoomBuilder) BuildGarbageCollector(g *GarbageCollectorI) {
	*g = room.garbageCollector
}
