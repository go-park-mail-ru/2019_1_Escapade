package engine

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

type RoomMetrics struct {
	r *Room
	s SyncI
	e *RoomEvents
	f *RoomField
	i *RoomInformation
}

func (room *RoomMetrics) Init(r *Room, s SyncI, e *RoomEvents,
	f *RoomField, i *RoomInformation) {
	room.r = r
	room.s = s
	room.e = e
	room.f = f
	room.i = i
}

func (room *RoomMetrics) Observe(needMetrics bool, cancel bool) {
	if !needMetrics {
		return
	}
	room.s.do(func() {
		var (
			roomType        string
			anonymous, mode int
		)
		if cancel {
			roomType = "aborted"
			metrics.AbortedRooms.Inc()
		} else {
			roomType = "finished"
			metrics.FinishedRooms.Inc()
		}
		if !room.r.Settings.NoAnonymous {
			anonymous = 1
		}
		if room.r.Settings.Deathmatch {
			mode = 1
		}

		size := float64(room.r.Settings.Width * room.r.Settings.Height)

		utils.Debug(false, "metrics RoomPlayers", room.r.Settings.Players)
		metrics.RoomPlayers.WithLabelValues(roomType).Observe(float64(room.r.Settings.Players))
		utils.Debug(false, "metrics difficult", room.f.Field.Difficult)
		metrics.RoomDifficult.WithLabelValues(roomType).Observe(float64(room.f.Field.Difficult))
		utils.Debug(false, "metrics size", size)
		metrics.RoomSize.WithLabelValues(roomType).Observe(size)
		utils.Debug(false, "metrics TimeToPlay", room.r.Settings.TimeToPlay)
		metrics.RoomTime.WithLabelValues(roomType).Observe(float64(room.r.Settings.TimeToPlay))
		if !cancel {
			openProcent := 1 - float64(float64(room.f.Field.cellsLeft())/size)
			utils.Debug(false, "metrics openProcent", openProcent)
			metrics.RoomOpenProcent.Observe(openProcent)

			utils.Debug(false, "metrics playing time", room.e.playingTime().Seconds())
			metrics.RoomTimePlaying.Observe(room.e.playingTime().Seconds())
		}
		metrics.RoomMode.WithLabelValues(roomType, utils.String(mode)).Inc()
		metrics.RoomAnonymous.WithLabelValues(roomType, utils.String(anonymous)).Inc()
		utils.Debug(false, "metrics recruitmentTime", room.e.recruitmentTime().Seconds())
		metrics.RoomTimeSearchingPeople.WithLabelValues(roomType).Observe(room.e.recruitmentTime().Seconds())
	})
}
