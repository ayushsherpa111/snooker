package game

import (
	"image/color"

	config "github.com/ayushsherpa111/snooker/Config"
)

type circle struct {
	id uint8
	// Velocity in the x coordinate
	vx float32
	// Velocity in the y coordinate
	vy float32

	// Center X position of the circle
	cx float32
	// Center Y position of the circle
	cy float32

	color color.Color

	isSelectable bool
}

type cueStick struct {
	// position from cue ball
	cx, cy float32

	// width of cue stick
	strokeWidth float32

	// maximum velocity that can be given to cue
	maxPower float32

	// draw stick on the board
	drawStick bool
}

type cueStartState struct {
	x, y       float32
	color      color.Color
	selectable bool
}

var (
	bg          = color.RGBA{33, 136, 34, 0xff}
	white       = color.RGBA{0xff, 0xff, 0xf2, 0xff}
	red         = color.RGBA{0xff, 0x00, 0x00, 0xFF}
	yellow      = color.RGBA{0xff, 0xff, 0x00, 0xFF}
	green       = color.RGBA{0x01, 0x73, 0x01, 0xFF}
	brown       = color.RGBA{0x51, 0x3b, 0x16, 0xFF}
	blue        = color.RGBA{0x00, 0x00, 0xff, 0xFF}
	pink        = color.RGBA{0xEA, 0x33, 0xE7, 0xFF}
	black       = color.Black
	board_state = []cueStartState{
		{
			x:          config.PRE_MARKED_LINE - (14 * circRadius),
			y:          config.WIN_HEIGHT / 2,
			color:      white,
			selectable: true,
		},
		{
			x:     config.WIN_WIDTH/3 + config.WIN_WIDTH/3,
			y:     config.WIN_HEIGHT / 2,
			color: red,
		},
		{
			x:     config.PRE_MARKED_LINE,
			y:     (config.WIN_HEIGHT / 2) - config.SEMI_CIRC_RADI,
			color: green,
		},
		{
			x:     config.PRE_MARKED_LINE,
			y:     config.WIN_HEIGHT / 2,
			color: brown,
		},
		{
			x:     config.WIN_WIDTH / 5,
			y:     config.WIN_HEIGHT/2 + config.SEMI_CIRC_RADI,
			color: yellow,
		},
		{
			x:     config.WIN_WIDTH / 2,
			y:     config.WIN_HEIGHT / 2,
			color: blue,
		},
		{
			x:     820 - 2*circRadius - 4,
			y:     258,
			color: pink,
		},
		{
			x:     config.WIN_WIDTH - (4 * circRadius),
			y:     config.WIN_HEIGHT / 2,
			color: black,
		},
	}
)

func mirrorPoint(slope, coeff, y, x1, y1 float32) (float32, float32) {
	temp := -2 * (slope*x1 + y*y1 + coeff) / (slope*slope + y*y)
	return temp*slope + x1, temp*y + y1
}

func slope(x1, y1, x2, y2 float32) float32 {
	return (y2 - y1) / (x2 - x1)
}
