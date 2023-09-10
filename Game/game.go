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
	cueBalls []*circle
	pockets  []circle
	board    *ebiten.Image
	cue      *circle
	cueStick cueStick
	Debug    bool
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

var isCueInMotion = false

func (g *Game) Update() error {
	if !isBoardSet {
		g.setBoard()
		isBoardSet = true
	}
	mouseX, mouseY := ebiten.CursorPosition()
	fMouseX, fMouseY := float64(mouseX), float64(mouseY)

	// fmt.Printf("(%d %d) (%d %d)\n", mouseX, mouseY, mouseX, -mouseY)

	if e := g.handleInputs(fMouseX, fMouseY); e != nil {
		return e
	}

	if isCueSelected && !gameStarted {
		g.cue.c_v.x = fMouseX
		g.cue.c_v.y = fMouseY
	}

	if g.cueStick.drawStick {
		m := 1 / slope(fMouseX, -1*fMouseY, g.cue.c_v.x, -1*g.cue.c_v.y)
		coeff := (-1*m)*g.cue.c_v.x + g.cue.c_v.y
		g.cueStick.cx, g.cueStick.cy = mirrorPoint(m, coeff, 1, fMouseX, -1*fMouseY)
		g.cueStick.cy *= -1

		if g.Debug {
			fmt.Printf("Mouse: (%f,%f) Cue(%f,%f)\n", fMouseX, fMouseY, g.cue.c_v.x, g.cue.c_v.y)
			fmt.Printf("Slope: %f\n", m)
			fmt.Printf("Coeff: %f\n", coeff)
			fmt.Printf("CueStick: (%f,%f)\n", g.cueStick.cx, g.cueStick.cy)
		}
	}

	g.accumulateForces()
	g.checkConstraints()
	g.move()
	return nil
}

func (g *Game) isOverlapping(x1, y1, x2, y2 float64, distance float64) bool {
	return findDistance(x1, y1, x2, y2) < distance
}

func (g *Game) Draw(screen *ebiten.Image) {
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

	g.cueBalls = make([]*circle, 0, 22)
	for i := 0; i < 8; i++ {
		if i == 1 {
			// red balls
			for x := 1; x < 16; x++ {
				red_ball := circle{
					id:           uint8(x),
					color:        board_state[1].color,
					c_v:          Vector{0, 0},
					isSelectable: board_state[1].selectable,
				}
				g.cueBalls = append(g.cueBalls, &red_ball)
			}
		} else {
			id := 0
			if i != 0 {
				id = i + 15
			}
			other_balls := circle{
				id:           uint8(id),
				color:        board_state[i].color,
				c_v:          Vector{board_state[i].x, board_state[i].y},
				p_v:          Vector{board_state[i].x, board_state[i].y},
				isSelectable: board_state[i].selectable,
			}
			g.cueBalls = append(g.cueBalls, &other_balls)
		}
	}
	g.arrangePyramids(5, 0, 0, 0)
	g.cue = g.cueBalls[config.CUE_BALL_IDX]
	g.cueStick = cueStick{
		drawStick:   false,
		maxPower:    config.BASE_POWER,
		strokeWidth: 2,
	}
}

func (g *Game) drawBoard(target *ebiten.Image) {
	target.DrawImage(g.board, nil)
	for _, ball := range g.cueBalls {
		vector.DrawFilledCircle(
			target,
			float32(ball.c_v.x),
			float32(ball.c_v.y),
			circRadius,
			ball.color,
			antialias,
		)
	}
	if g.cueStick.drawStick {
		vector.StrokeLine(
			target,
			float32(g.cue.c_v.x),
			float32(g.cue.c_v.y),
			float32(g.cueStick.cx),
			float32(g.cueStick.cy),
			float32(g.cueStick.strokeWidth),
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
		x := float64(baseTriangleX) - float64(depth*circRadius)
		y := float64(
			baseTriangleY,
		) + float64(
			(2*circRadius)*i,
		) + float64(
			padding*circRadius,
		)
		g.cueBalls[i+idx].c_v.x = x
		g.cueBalls[i+idx].c_v.y = y

		g.cueBalls[i+idx].p_v.x = x
		g.cueBalls[i+idx].p_v.y = y
	}
	g.arrangePyramids(steps-1, depth+2, padding+1, steps+idx)
}

func (g *Game) handleInputs(fMouseX, fMouseY float64) error {
	switch {
	case inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) && !isCueSelected:
		if g.isOverlapping(fMouseX, fMouseY, g.cue.c_v.x, g.cue.c_v.y, circRadius) {
			isCueSelected = true
		}
	case inpututil.IsKeyJustPressed(ebiten.KeyQ):
		return ebiten.Termination
	case inpututil.IsKeyJustPressed(ebiten.KeyR):
		g.cue.c_v.x = 250
		g.cue.c_v.y = 250
		g.cue.p_v.x = 250
		g.cue.p_v.y = 250
		g.cue.a_v.x = 0
		g.cue.a_v.y = 0
        isCueInMotion = false
	case inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight):
		isCueSelected = false
	case inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft):
		g.cueStick.drawStick = true
	case inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft):
		// shoot ball
		if !isCueInMotion {
			velocity := Vector{
				x: config.BASE_POWER * (g.cueStick.cx - g.cue.c_v.x),
				y: config.BASE_POWER * (g.cueStick.cy - g.cue.c_v.y),
			}
			g.SetVelocity(config.CUE_BALL_IDX, velocity)
		}
		g.cueStick.drawStick = false
		isCueInMotion = true
	default:
	}
	return nil
}
