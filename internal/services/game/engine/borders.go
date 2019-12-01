package engine

type Borders struct {
	left   int32
	right  int32
	down   int32
	up     int32
	Width  int32
	Height int32
}

func (b *Borders) Init(cell Cell, width, height, mineArea int32) {
	var (
		x = cell.X
		y = cell.Y
	)
	b.Width = width
	b.Height = height
	b.left = b.fixWidth(x - mineArea)
	b.right = b.fixWidth(x + mineArea)
	b.down = b.fixHeight(y - mineArea)
	b.up = b.fixHeight(y + mineArea)
}

func (b *Borders) fixWidth(i int32) int32 {
	if i < 0 {
		return 0
	}
	if i >= b.Width {
		return b.Width - 1
	}
	return i
}

func (b *Borders) fixHeight(i int32) int32 {
	if i < 0 {
		return 0
	}
	if i >= b.Height {
		return b.Height - 1
	}
	return i
}
