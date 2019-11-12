package engine

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/synced"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// MetricsI handle sending metrics
// Strategy Pattern
type MetricsI interface {
	Observe(needMetrics bool, cancel bool)
}

// RoomMetrics implements MetricsStrategyI
type RoomMetrics struct {
	s synced.SyncI
	e EventsI
	f FieldProxyI

	size            float64
	anonymous, mode string

	settings *models.RoomSettings
}

// Init configure dependencies with other components of the room
func (room *RoomMetrics) Init(builder RBuilderI, rs *models.RoomSettings) {
	builder.BuildSync(&room.s)
	builder.BuildEvents(&room.e)
	builder.BuildField(&room.f)

	room.settings = rs
	room.size = float64(rs.Width * rs.Height)
	if !rs.NoAnonymous {
		room.anonymous = utils.String(1)
	}
	if rs.Deathmatch {
		room.mode = utils.String(1)
	}
}

func (room *RoomMetrics) Observe(needMetrics bool, cancel bool) {
	if !needMetrics {
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

			utils.Debug(false, "metrics playing time", room.e.playingTime().Seconds())
			metrics.RoomTimePlaying.Observe(room.e.playingTime().Seconds())
		}
		metrics.RoomMode.WithLabelValues(roomType, room.mode).Inc()
		metrics.RoomAnonymous.WithLabelValues(roomType, room.anonymous).Inc()
		utils.Debug(false, "metrics recruitmentTime", room.e.recruitmentTime().Seconds())
		metrics.RoomTimeSearchingPeople.WithLabelValues(roomType).Observe(room.e.recruitmentTime().Seconds())
	})
}
