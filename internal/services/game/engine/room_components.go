package engine

import "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"

// RBuilderI create room components, bind them, and add them to room
// 	to the room
// ABuilder Pattern
type RBuilderI interface {
	Build(r *Room, ra *RoomArgs)

	BuildInformation(i *RoomInformationI)
	BuildField(f *FieldProxyI)
	BuildSync(s *synced.SyncI)
	BuildAPI(a *RoomRequestsI)
	BuildLobby(l *LobbyProxyI)
	BuildModelsAdapter(m *RModelsI)
	BuildSender(s *RSendI)
	BuildConnectionEvents(c *RClientI)
	BuildPeople(p *PeopleI)
	BuildEvents(e *EventsI)
	BuildRecorder(r *ActionRecorderI)
	BuildMetrics(m *RoomMetrics)
	BuildMessages(m *MessagesI)
	BuildGarbageCollector(g *GarbageCollectorI)
}

// RoomBuilder implements ComponentBuilderI
type RoomBuilder struct {
	info             *RoomInformation
	field            *RoomField
	sync             *synced.SyncWgroup
	api              *RoomAPI
	lobby            *RoomLobby
	models           *RoomModels
	sender           *RoomSender
	client           *RClient
	people           *RoomPeople
	events           *RoomEvents
	record           *RoomRecorder
	metrics          *RoomMetrics
	messages         *RoomMessages
	garbageCollector *RoomGarbageCollector
}

// Build create Room components and set them to Room object
func (builder *RoomBuilder) Build(room *Room, args *RoomArgs) {

	builder.createComponents()
	builder.configureDependencies(room, args)
	builder.set(room)
}

func (builder *RoomBuilder) createComponents() {
	builder.sync = &synced.SyncWgroup{}
	builder.info = &RoomInformation{}
	builder.api = &RoomAPI{}
	builder.lobby = &RoomLobby{}
	builder.field = &RoomField{}
	builder.models = &RoomModels{}
	builder.sender = &RoomSender{}
	builder.people = &RoomPeople{}
	builder.client = &RClient{}
	builder.events = &RoomEvents{}
	builder.metrics = &RoomMetrics{}
	builder.record = &RoomRecorder{}
	builder.messages = &RoomMessages{}
	builder.garbageCollector = &RoomGarbageCollector{}
}

func (builder *RoomBuilder) configureDependencies(r *Room, args *RoomArgs) {
	builder.sync.Init(args.c.Wait.Duration)
	builder.info.Init(args.rs, args.id, args.DBRoomID, args.lobby.location())
	builder.api.Init(builder)
	builder.lobby.Init(builder, r, args.lobby)
	builder.field.Init(builder, args.Field, args.rs.Deathmatch)
	builder.models.Init(builder)
	builder.sender.Init(builder)
	builder.people.Init(builder, args.rs)
	builder.client.Init(builder, args.rs.Deathmatch)
	builder.events.Init(builder, args.rs, args.c.CanClose)
	builder.metrics.Init(builder, args.rs, args.lobby.config().Metrics)
	builder.record.Init(builder)
	builder.messages.Init(builder, args.lobby.ChatService,
		args.DBchatID, args.lobby.location())
	builder.garbageCollector.Init(builder, args.c.GarbageCollector.Duration,
		args.c.Timeouts)
}

func (builder *RoomBuilder) set(room *Room) {
	room.sync = builder.sync
	room.info = builder.info
	room.api = builder.api
	room.lobby = builder.lobby
	room.field = builder.field
	room.models = builder.models
	room.sender = builder.sender
	room.people = builder.people
	room.client = builder.client
	room.events = builder.events
	room.metrics = builder.metrics
	room.record = builder.record
	room.messages = builder.messages
	room.garbageCollector = builder.garbageCollector
}

// BuildInformation set RoomInformationI implementation
func (builder *RoomBuilder) BuildInformation(i *RoomInformationI) {
	*i = builder.info
}

// BuildField set FieldProxyI implementation
func (builder *RoomBuilder) BuildField(f *FieldProxyI) {
	*f = builder.field
}

// BuildSync set SyncI implementation
func (builder *RoomBuilder) BuildSync(s *synced.SyncI) {
	*s = builder.sync
}

// BuildAPI set APIStrategyI implementation
func (builder *RoomBuilder) BuildAPI(a *RoomRequestsI) {
	*a = builder.api
}

// BuildLobby set APIStrategyI implementation
func (builder *RoomBuilder) BuildLobby(l *LobbyProxyI) {
	*l = builder.lobby
}

// BuildModelsAdapter set ModelsAdapterI implementation
func (builder *RoomBuilder) BuildModelsAdapter(m *RModelsI) {
	*m = builder.models
}

// BuildSender set SendStrategyI implementation
func (builder *RoomBuilder) BuildSender(s *RSendI) {
	*s = builder.sender
}

// BuildConnectionEvents set ConnectionEventsStrategyI implementation
func (builder *RoomBuilder) BuildConnectionEvents(c *RClientI) {
	*c = builder.client
}

// BuildPeople set PeopleI implementation
func (builder *RoomBuilder) BuildPeople(p *PeopleI) {
	*p = builder.people
}

// BuildEvents set EventsI implementation
func (builder *RoomBuilder) BuildEvents(e *EventsI) {
	*e = builder.events
}

// BuildRecorder set ActionRecorderProxyI implementation
func (builder *RoomBuilder) BuildRecorder(r *ActionRecorderI) {
	*r = builder.record
}

// BuildMetrics set MetricsStrategyI implementation
func (builder *RoomBuilder) BuildMetrics(m *RoomMetrics) {
	*m = *builder.metrics
}

// BuildMessages set MessagesProxyI implementation
func (builder *RoomBuilder) BuildMessages(m *MessagesI) {
	*m = builder.messages
}

// BuildGarbageCollector set GarbageCollectorI implementation
func (builder *RoomBuilder) BuildGarbageCollector(g *GarbageCollectorI) {
	*g = builder.garbageCollector
}
