package engine

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"

	room_ "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/engine/room"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/metrics"
)

// RoomMetrics implements MetricsStrategyI
type RoomMetrics struct {
	s synced.SyncI
	i RoomInformationI
	f FieldProxyI

	size            float64
	anonymous, mode string
	needMetrics     bool

	settings *models.RoomSettings
}

// init struct's values
func (room *RoomMetrics) init(settings *models.RoomSettings, needMetrics bool) {
	room.needMetrics = needMetrics
	room.settings = settings
	room.size = float64(settings.Width * settings.Height)
	if !settings.NoAnonymous {
		room.anonymous = utils.String(1)
	}
	if settings.Deathmatch {
		room.mode = utils.String(1)
	}
}

// build components
func (room *RoomMetrics) build(builder RBuilderI) {
	builder.BuildSync(&room.s)
	builder.BuildField(&room.f)
	builder.BuildInformation(&room.i)
}

// subscribe to room events
func (room *RoomMetrics) subscribe(builder RBuilderI) {
	var (
		events   EventsI
		messages MessagesI
	)

	builder.BuildEvents(&events)
	builder.BuildMessages(&messages)

	room.eventsSubscribe(events)
	room.messagesSubscribe(messages)
}

// Init configure dependencies with other components of the room
func (room *RoomMetrics) Init(builder RBuilderI, rs *models.RoomSettings, needMetrics bool) {
	room.init(rs, needMetrics)
	room.build(builder)
	room.subscribe(builder)
}

// record metrics
func (room *RoomMetrics) record(cancel bool) {
	if !room.needMetrics {
		return
	}
	room.s.Do(func() {
		var (
			roomType string
		)
		if cancel {
			roomType = "aborted"
			metrics.AbortedRooms.Inc()
		} else {
			roomType = "finished"
			metrics.FinishedRooms.Inc()
		}

		utils.Debug(false, "metrics RoomPlayers", room.settings.Players)
		metrics.RoomPlayers.WithLabelValues(roomType).Observe(float64(room.settings.Players))
		utils.Debug(false, "metrics difficult", room.f.Field().difficult())
		metrics.RoomDifficult.WithLabelValues(roomType).Observe(float64(room.f.Field().difficult()))
		utils.Debug(false, "metrics size", room.size)
		metrics.RoomSize.WithLabelValues(roomType).Observe(room.size)
		utils.Debug(false, "metrics TimeToPlay", room.settings.TimeToPlay)
		metrics.RoomTime.WithLabelValues(roomType).Observe(float64(room.settings.TimeToPlay))
		if !cancel {
			openProcent := 1 - float64(float64(room.f.Field().cellsLeft())/room.size)
			utils.Debug(false, "metrics openProcent", openProcent)
			metrics.RoomOpenProcent.Observe(openProcent)

			utils.Debug(false, "metrics playing time", room.i.PlayingTime().Seconds())
			metrics.RoomTimePlaying.Observe(room.i.PlayingTime().Seconds())
		}
		metrics.RoomMode.WithLabelValues(roomType, room.mode).Inc()
		metrics.RoomAnonymous.WithLabelValues(roomType, room.anonymous).Inc()
		utils.Debug(false, "metrics recruitmentTime", room.i.RecruitmentTime().Seconds())
		metrics.RoomTimeSearchingPeople.WithLabelValues(roomType).Observe(room.i.RecruitmentTime().Seconds())
	})
}

// eventsFinished is called when game finished
func (room *RoomMetrics) eventsFinished(msg synced.Msg) {
	result, ok := msg.Extra.(room_.FinishResults)
	if !ok {
		return
	}
	room.record(result.Cancel)
}

// eventsSubscribe subscibe to events associated with room's status
func (room *RoomMetrics) eventsSubscribe(e EventsI) {
	observer := synced.NewObserver(
		synced.NewPair(room_.StatusFinished, room.eventsFinished))
	e.Observe(observer.AddPublisherCode(room_.UpdateStatus))
}

// messagesSubscribe subscibe to events associated with room's chat
func (room *RoomMetrics) messagesSubscribe(m MessagesI) {
	observer := synced.NewObserver(
		synced.NewPairNoArgs(room_.Delete, metrics.RoomsMessages.Dec),
		synced.NewPairNoArgs(room_.Add, metrics.RoomsMessages.Inc))
	m.Observe(observer.AddPublisherCode(room_.UpdateChat))
}
