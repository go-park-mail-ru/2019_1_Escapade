package engine

import re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"

// Borders - the boundaries of the min area surrounding the specified cell
type Borders struct {
	left, right, down, up   int32
	fieldWidth, fieldHeight int32
}

// Init the Boders struct with cell(center of area), field's weight and height,
// 	and area radius
func (b *Borders) Init(center Cell,
	fieldWidth, fieldHeight, radius int32) error {

	if fieldWidth <= 0 || fieldHeight <= 0 || radius <= 0 {
		return re.ErrorWrongBordersParams(fieldWidth, fieldHeight, radius)
	}

	if !center.AreCoordinatesValid(fieldWidth, fieldHeight) {
		return re.ErrorCellOutside()
	}
	var (
		x = center.X
		y = center.Y
	)
	b.fieldWidth = fieldWidth
	b.fieldHeight = fieldHeight
	b.left = b.fixWidth(x - radius)
	b.right = b.fixWidth(x + radius)
	b.down = b.fixHeight(y - radius)
	b.up = b.fixHeight(y + radius)
	return nil
}

func (b *Borders) fixWidth(i int32) int32 {
	if i < 0 {
		return 0
	}
	if i >= b.fieldWidth {
		return b.fieldWidth - 1
	}
	return i
}

func (b *Borders) fixHeight(i int32) int32 {
	if i < 0 {
		return 0
	}
	if i >= b.fieldHeight {
		return b.fieldHeight - 1
	}
	return i
}
