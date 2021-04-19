package xwindow

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xinerama"
	"github.com/knobik/go-display/internal/util"
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

func GetDisplayBounds(displayIndex int) (rect util.Bounds) {
	defer func() {
		e := recover()
		if e != nil {
			rect = util.Bounds{}
		}
	}()

	c, err := xgb.NewConn()
	if err != nil {
		return util.Bounds{}
	}
	defer c.Close()

	err = xinerama.Init(c)
	if err != nil {
		return util.Bounds{}
	}

	reply, err := xinerama.QueryScreens(c).Reply()
	if err != nil {
		return util.Bounds{}
	}

	if displayIndex >= int(reply.Number) {
		return util.Bounds{}
	}

	primary := reply.ScreenInfo[0]
	x0 := int(primary.XOrg)
	y0 := int(primary.YOrg)

	screen := reply.ScreenInfo[displayIndex]
	x := int(screen.XOrg) - x0
	y := int(screen.YOrg) - y0
	w := int(screen.Width)
	h := int(screen.Height)
	rect = util.MakeBounds(x, y, x+w, y+h)
	return rect
}
