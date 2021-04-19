package display

import (
	"github.com/knobik/go-display/internal/util"
	"github.com/knobik/go-display/internal/xwindow"
)

// NumActiveDisplays returns the number of active displays.
func NumActiveDisplays() int {
	return xwindow.NumActiveDisplays()
}

// GetDisplayBounds returns the bounds of displayIndex'th display.
// The main display is displayIndex = 0.
func GetDisplayBounds(displayIndex int) util.Bounds {
	return xwindow.GetDisplayBounds(displayIndex)
}
