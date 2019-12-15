package engine

import (
	"testing"

	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	. "github.com/smartystreets/goconvey/convey"
)

// integration
func TestWrongBordersInit(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Test borders Init", t, func() {
		var (
			b                          = &Borders{}
			width, height, mines, x, y int32
		)
		width, height = 10, 10

		for mines = -width; mines < width*2; mines++ {
			for x = -width; x < width*2; x++ {
				for y = -height; y < height*2; y++ {
					_, err := b.Init(Cell{X: x, Y: y}, width, height, mines)
					if mines <= 0 {
						So(err, ShouldResemble, re.ErrorWrongBordersParams(width,
							height, mines))
					} else if x < 0 || y < 0 || x >= width || y >= height {
						So(err, ShouldResemble, re.ErrorCellOutside())
					} else {
						So(err, ShouldBeNil)
					}
				}
			}
		}
	})
}

// unit
func TestFixWidth(t *testing.T) {

	var currX, currY, maxW, maxY, rad int32
	currX, currY, maxW, maxY, rad = 20, 10, 50, 70, 5
	var borders, err = NewBorders(Cell{X: currX, Y: currY}, maxW, maxY, rad)
	Convey("Given borders", t, func() {
		var x int32
		So(err, ShouldBeNil)
		Convey("X < less then 0", func() {
			x = -1
			Convey("When fix X", func() {
				x = borders.fixWidth(x)
				Convey("Then X = 0 ", func() {
					So(x, ShouldEqual, 0)
				})
			})
		})
		Convey("X > more then max", func() {
			x = maxW
			Convey("When fix X", func() {
				x = borders.fixWidth(x)
				Convey("Then X = maxW-1 ", func() {
					So(x, ShouldEqual, maxW-1)
				})
			})
		})
		Convey("X inside", func() {
			x = maxW / 2
			Convey("When fix X", func() {
				x = borders.fixWidth(x)

				Convey("Then X didnt chande ", func() {
					So(x, ShouldEqual, maxW/2)
				})
			})
		})
	})
}

// unit
func TestFixheight(t *testing.T) {

	var currX, currY, maxW, maxY, rad int32
	currX, currY, maxW, maxY, rad = 20, 10, 50, 70, 5
	var borders, err = NewBorders(Cell{X: currX, Y: currY}, maxW, maxY, rad)
	Convey("Given borders", t, func() {
		var y int32
		So(err, ShouldBeNil)
		Convey("Y < less then 0", func() {
			y = -1
			Convey("When fix Y", func() {
				y = borders.fixWidth(y)
				Convey("Then Y = 0 ", func() {
					So(y, ShouldEqual, 0)
				})
			})
		})
		Convey("Y > more then max", func() {
			y = maxY
			Convey("When fix Y", func() {
				y = borders.fixHeight(y)
				Convey("Then Y = maxY-1 ", func() {
					So(y, ShouldEqual, maxY-1)
				})
			})
		})
		Convey("Y inside", func() {
			y = maxY / 2
			Convey("When fix Y", func() {
				y = borders.fixHeight(y)

				Convey("Then Y didnt chande ", func() {
					So(y, ShouldEqual, maxY/2)
				})
			})
		})
	})
}
