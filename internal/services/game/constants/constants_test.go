package constants

import (
	"fmt"
	"testing"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	rc = RoomConfiguration{
		NameMin:          10,
		NameMax:          30,
		TimeToPrepareMin: 0,
		TimeToPrepareMax: 60,
		TimeToPlayMin:    0,
		TimeToPlayMax:    7200,
		PlayersMin:       2,
		PlayersMax:       100,
		ObserversMax:     100,
		Set:              true,
	}
	fc = FieldConfiguration{
		WidthMin:  5,
		WidthMax:  100,
		HeightMin: 5,
		HeightMax: 100,
		Set:       true,
	}
)

type RepositoryFake struct {
	fieldPath, roomPath string
}

func (rfs *RepositoryFake) Init(field, room string) *RepositoryFake {
	rfs.fieldPath = field
	rfs.roomPath = room
	return rfs
}

func newRepositoryFake(field, room string) *RepositoryFake {
	return new(RepositoryFake).Init(field, room)
}

// getRoom load Room config(.json file) from FS by its path
func (rfs *RepositoryFake) GetRoom(path string) (RoomConfiguration, error) {
	if path != rfs.roomPath {
		return RoomConfiguration{}, fmt.Errorf(
			"Incorrect path to room configuration:%s", path)
	}
	return rc, nil
}

// getField load field config(.json file) from FS by its path
func (rfs *RepositoryFake) GetField(path string) (FieldConfiguration, error) {
	if path != rfs.fieldPath {
		return FieldConfiguration{}, fmt.Errorf(
			"Incorrect path to field configuration:%s", path)
	}
	return fc, nil
}

func rs() *models.RoomSettings {
	return &models.RoomSettings{
		Name:          "this room has no name",
		Width:         7,
		Height:        7,
		Players:       2,
		Observers:     10,
		TimeToPrepare: 5,
		TimeToPlay:    60,
		Mines:         2,
	}
}

// unit
func TestInitRoom(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Given correct and wrong path to room configuration", t, func() {
		var (
			correct = "right"
			wrong   = "another"
			rep     = newRepositoryFake("", correct)
		)
		ROOM = RoomConfiguration{}
		Convey("When initialize room constants with correct path", func() {
			err := InitRoom(rep, correct)
			Convey("Then no error should happen and room constants set", func() {
				So(err, ShouldBeNil)
				So(ROOM, ShouldResemble, rc)
			})
		})

		ROOM = RoomConfiguration{}
		Convey("When initialize room constants with wrong path", func() {
			err := InitRoom(rep, wrong)
			Convey("Then error should happen and room constants not set", func() {
				So(err, ShouldNotBeNil)
				So(ROOM, ShouldResemble, RoomConfiguration{})
			})
		})
	})
}

// unit
func TestInitField(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Given correct and wrong path to room configuration", t, func() {
		var (
			correct = "right"
			wrong   = "another"
			rep     = newRepositoryFake(correct, "")
		)
		FIELD = FieldConfiguration{}
		Convey("When initialize field constants with correct path", func() {
			err := InitField(rep, correct)
			Convey("Then no error should happen and field constants set", func() {
				So(err, ShouldBeNil)
				So(FIELD, ShouldResemble, fc)
			})
		})

		FIELD = FieldConfiguration{}
		Convey("When initialize field constants with wrong path", func() {
			err := InitField(rep, wrong)
			Convey("Then error should happen and field constants not set", func() {
				So(err, ShouldNotBeNil)
				So(FIELD, ShouldResemble, FieldConfiguration{})
			})
		})
	})
}

// unit
func TestCheckName(t *testing.T) {

	settings := rs()
	rep := newRepositoryFake("1", "2")
	InitField(rep, "1")
	InitRoom(rep, "2")

	Convey("Given roomsettings with too short name", t, func() {
		settings.Name = ""
		var i int32
		for i = 0; i < ROOM.NameMin-1; i++ {
			settings.Name += "a"
		}
		Convey("When check these settings", func() {
			err := Check(settings)
			Convey("Then the error will be given", func() {
				So(err, ShouldResemble, ErrorRoomName(settings))
			})
		})
	})

	Convey("Given roomsettings with valid short name", t, func() {
		settings.Name = ""
		var i int32
		for i = 0; i < ROOM.NameMin; i++ {
			settings.Name += "a"
		}
		Convey("When check these settings", func() {
			err := Check(settings)
			Convey("Then no error will be given", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given roomsettings with valid large name", t, func() {
		settings.Name = ""
		var i int32
		for i = 0; i < ROOM.NameMax; i++ {
			settings.Name += "a"
		}
		Convey("When check these settings", func() {
			err := Check(settings)
			Convey("Then no error will be given", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given roomsettings with too large name", t, func() {
		settings.Name = ""
		var i int32
		for i = 0; i < ROOM.NameMax+1; i++ {
			settings.Name += "a"
		}
		Convey("When check these settings", func() {
			err := Check(settings)
			Convey("Then error will be given", func() {
				So(err, ShouldResemble, ErrorRoomName(settings))
			})
		})
	})
}

// unit
func TestCheckWidth(t *testing.T) {

	settings := rs()
	rep := newRepositoryFake("1", "2")
	InitField(rep, "1")
	InitRoom(rep, "2")

	Convey("Given roomsettings with too small field width", t, func() {
		settings.Width = FIELD.WidthMin - 1
		Convey("When check these settings", func() {
			err := Check(settings)
			Convey("Then the error will be given", func() {
				So(err, ShouldResemble, ErrorFieldWidth(settings))
			})
		})
	})

	Convey("Given roomsettings with valid small field width", t, func() {
		settings.Width = FIELD.WidthMin
		Convey("When check these settings", func() {
			err := Check(settings)
			Convey("Then no error will be given", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given roomsettings with valid large field width", t, func() {
		settings.Width = FIELD.WidthMax
		Convey("When check these settings", func() {
			err := Check(settings)
			Convey("Then no error will be given", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given roomsettings with too large field width", t, func() {
		settings.Width = FIELD.WidthMax + 1
		Convey("When check these settings", func() {
			err := Check(settings)
			Convey("Then error will be given", func() {
				So(err, ShouldResemble, ErrorFieldWidth(settings))
			})
		})
	})
}

// unit
func TestCheckHeight(t *testing.T) {

	settings := rs()
	rep := newRepositoryFake("1", "2")
	InitField(rep, "1")
	InitRoom(rep, "2")

	Convey("Given roomsettings with too small field height", t, func() {
		settings.Height = FIELD.HeightMin - 1
		Convey("When check these settings", func() {
			err := Check(settings)
			Convey("Then the error will be given", func() {
				So(err, ShouldResemble, ErrorFieldHeight(settings))
			})
		})
	})

	Convey("Given roomsettings with valid small field height", t, func() {
		settings.Height = FIELD.HeightMin
		Convey("When check these settings", func() {
			err := Check(settings)
			Convey("Then no error will be given", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given roomsettings with valid large field height", t, func() {
		settings.Height = FIELD.HeightMax
		Convey("When check these settings", func() {
			err := Check(settings)
			Convey("Then no error will be given", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given roomsettings with too large field height", t, func() {
		settings.Height = FIELD.HeightMax + 1
		Convey("When check these settings", func() {
			err := Check(settings)
			Convey("Then error will be given", func() {
				So(err, ShouldResemble, ErrorFieldHeight(settings))
			})
		})
	})
}

// unit
func TestCheckPlayers(t *testing.T) {

	settings := rs()
	rep := newRepositoryFake("1", "2")
	InitField(rep, "1")
	InitRoom(rep, "2")

	Convey("Given roomsettings with too small player amount", t, func() {
		settings.Players = ROOM.PlayersMin - 1
		Convey("When check these settings", func() {
			err := Check(settings)
			Convey("Then the error will be given", func() {
				So(err, ShouldResemble, ErrorPlayers(settings))
			})
		})
	})

	Convey("Given roomsettings with valid small player amount", t, func() {
		settings.Players = ROOM.PlayersMin
		Convey("When check these settings", func() {
			err := Check(settings)
			Convey("Then no error will be given", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given roomsettings with valid large player amount, bigger then field size", t, func() {
		settings.Players = ROOM.PlayersMax
		Convey("When check these settings", func() {
			err := Check(settings)
			Convey("Then no error will be given", func() {
				So(err, ShouldResemble, ErrorPlayers(settings))
			})
		})
	})

	Convey("Given roomsettings with valid large player amount, equal as field size without mines", t, func() {
		settings.Players = settings.Height * settings.Width
		settings.Mines = 0
		Convey("When check these settings", func() {
			err := Check(settings)
			Convey("Then no error will be given", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given roomsettings with valid large player amount, less then field", t, func() {
		settings.Height = FIELD.HeightMax
		settings.Width = FIELD.WidthMax
		settings.Players = ROOM.PlayersMax
		Convey("When check these settings", func() {
			err := Check(settings)
			Convey("Then no error will be given", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given roomsettings with too large player amount", t, func() {
		settings.Height = FIELD.HeightMax
		settings.Width = FIELD.WidthMax
		settings.Players = ROOM.PlayersMax + 1
		Convey("When check these settings", func() {
			err := Check(settings)
			Convey("Then error will be given", func() {
				So(err, ShouldResemble, ErrorPlayers(settings))
			})
		})
	})
}

// integration
func TestUseCase(t *testing.T) {

	settings := rs()
	rep := newRepositoryFake("1", "2")
	InitField(rep, "1")
	InitRoom(rep, "2")

	Convey("Given roomsettings", t, func() {
		Convey("When check these settings", func() {
			err := Check(settings)
			Convey("Then no error will be given", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
