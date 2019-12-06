package engine

import (
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// unit
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
					err := b.Init(Cell{X: x, Y: y}, width, height, mines)
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
