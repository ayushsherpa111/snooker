package game

import (
	"fmt"
	"image/png"
	"os"

	config "github.com/ayushsherpa111/snooker/Config"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	cueBalls []circle
	pockets  []circle
	board    *ebiten.Image
	cue      *circle
	cueStick cueStick
    Debug bool
}

const (
	circRadius = 10
	antialias  = true
	cueBall    = 0
)

var (
	isBoardSet    = false
	isCueSelected = false
	gameStarted   = false
)

func (g *Game) Update() error {
	mouseX, mouseY := ebiten.CursorPosition()
	fMouseX, fMouseY := float32(mouseX), float32(mouseY)

	// fmt.Printf("(%d %d) (%d %d)\n", mouseX, mouseY, mouseX, -mouseY)

	switch {
	case inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) && !isCueSelected:
		if g.isOverlapping(fMouseX, fMouseY, g.cue.cx, g.cue.cy, circRadius) {
			isCueSelected = true
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyQ):
		return ebiten.Termination
	case inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight):
		isCueSelected = false
	case inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft):
		g.cueStick.drawStick = true
	case inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft):
		g.cueStick.drawStick = false
	default:
	}

	if isCueSelected && !gameStarted {
		g.cue.cx = fMouseX
		g.cue.cy = fMouseY
	}

	if g.cueStick.drawStick {
		m := 1 / slope(fMouseX, -1*fMouseY, g.cue.cx, -1*g.cue.cy)
		coeff := (-1*m)*g.cue.cx + g.cue.cy
		// y = mx + c
		// y - mx - c = 0 -> line equation of cursor
		// y - (-1/m) x - c = 0
		// y + 1/m x - c = 0
		g.cueStick.cx, g.cueStick.cy = mirrorPoint(
			m,
			coeff,
			1,
			fMouseX,
			-1*fMouseY,
		)
		g.cueStick.cy *= -1
        if g.Debug {
            fmt.Printf("Mouse: (%f,%f) Cue(%f,%f)\n", fMouseX, fMouseY, g.cue.cx, g.cue.cy)
            fmt.Printf("Slope: %f\n", m)
            fmt.Printf("Coeff: %f\n", coeff)
            fmt.Printf("CueStick: (%f,%f)\n", g.cueStick.cx, g.cueStick.cy)
        }
	}

	return nil
}

func (g *Game) isOverlapping(x1, y1, x2, y2 float32, distance float32) bool {
	return findDistance(x1, y1, x2, y2) < float64(distance)
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !isBoardSet {
		g.setBoard()
		isBoardSet = true
	}
	g.drawBoard(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.WIN_WIDTH, config.WIN_HEIGHT
}

func (g *Game) setBoard() {
	file, e := os.OpenFile("./pool.png", os.O_RDONLY, 0o777)
	if e != nil {
		fmt.Println("File not found")
	}
	board, _ := png.Decode(file)
	g.board = ebiten.NewImageFromImage(board)

	g.cueBalls = make([]circle, 0, 22)
	for i := 0; i < 8; i++ {
		if i == 1 {
			// red balls
			for x := 1; x < 16; x++ {
				red_ball := circle{
					id:           uint8(x),
					color:        board_state[1].color,
					cx:           0,
					cy:           0,
					isSelectable: board_state[1].selectable,
				}
				g.cueBalls = append(g.cueBalls, red_ball)
			}
		} else {
			id := 0
			if i != 0 {
				id = i + 15
			}
			other_balls := circle{
				id:           uint8(id),
				color:        board_state[i].color,
				cx:           board_state[i].x,
				cy:           board_state[i].y,
				isSelectable: board_state[i].selectable,
			}
			g.cueBalls = append(g.cueBalls, other_balls)
		}
	}
	g.arrangePyramids(5, 0, 0, 0)
	fmt.Println(g.cueBalls[1:6])
	g.cue = &g.cueBalls[0]
	g.cueStick = cueStick{
		drawStick:   false,
		maxPower:    config.BASE_POWER,
		strokeWidth: 2,
	}
}

func (g *Game) drawBoard(target *ebiten.Image) {
	target.DrawImage(g.board, nil)
	for _, ball := range g.cueBalls {
		vector.DrawFilledCircle(target, ball.cx, ball.cy, circRadius, ball.color, antialias)
	}
	if g.cueStick.drawStick {
		vector.StrokeLine(
			target,
			g.cue.cx,
			g.cue.cy,
			g.cueStick.cx,
			g.cueStick.cy,
			g.cueStick.strokeWidth,
			blue,
			true,
		)
	}
}

func (g *Game) arrangePyramids(steps, depth, padding, idx int) {
	if steps < 1 {
		return
	}
	baseTriangleX := 900
	baseTriangleY := config.WIN_HEIGHT/3 + (circRadius * 2)
	for i := 1; i <= steps; i++ {
		g.cueBalls[i+idx].cx = float32(baseTriangleX) - float32(depth*circRadius)
		g.cueBalls[i+idx].cy = float32(
			baseTriangleY,
		) + float32(
			(2*circRadius)*i,
		) + float32(
			padding*circRadius,
		)
	}
	g.arrangePyramids(steps-1, depth+2, padding+1, steps+idx)
}
