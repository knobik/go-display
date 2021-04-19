package xwindow

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xinerama"
	"github.com/knobik/go-display"
)

func NumActiveDisplays() (num int) {
	defer func() {
		e := recover()
		if e != nil {
			num = 0
		}
	}()

	c, err := xgb.NewConn()
	if err != nil {
		return 0
	}
	defer c.Close()

	err = xinerama.Init(c)
	if err != nil {
		return 0
	}

	reply, err := xinerama.QueryScreens(c).Reply()
	if err != nil {
		return 0
	}

	num = int(reply.Number)
	return num
}

func GetDisplayBounds(displayIndex int) (rect display.Bounds) {
	defer func() {
		e := recover()
		if e != nil {
			rect = display.Bounds{}
		}
	}()

	c, err := xgb.NewConn()
	if err != nil {
		return display.Bounds{}
	}
	defer c.Close()

	err = xinerama.Init(c)
	if err != nil {
		return display.Bounds{}
	}

	reply, err := xinerama.QueryScreens(c).Reply()
	if err != nil {
		return display.Bounds{}
	}

	if displayIndex >= int(reply.Number) {
		return display.Bounds{}
	}

	primary := reply.ScreenInfo[0]
	x0 := int(primary.XOrg)
	y0 := int(primary.YOrg)

	screen := reply.ScreenInfo[displayIndex]
	x := int(screen.XOrg) - x0
	y := int(screen.YOrg) - y0
	w := int(screen.Width)
	h := int(screen.Height)
	rect = display.MakeBounds(x, y, x+w, y+h)
	return rect
}
