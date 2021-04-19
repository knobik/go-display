package display

type Bounds struct {
	Left   int
	Top    int
	Right  int
	Bottom int
}

// MakeBounds The returned rectangle has minimum and maximum coordinates swapped if necessary so that
// it is well-formed.
func MakeBounds(left, top, right, bottom int) Bounds {
	if left > right {
		left, right = right, left
	}
	if top > bottom {
		top, bottom = bottom, top
	}
	return Bounds{
		Left:   left,
		Top:    top,
		Right:  right,
		Bottom: bottom,
	}
}
