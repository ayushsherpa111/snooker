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
}

const (
	circ_Radius = 10
	antialias   = true
)

var set_board = false

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		fmt.Println("Mouse button pressed")
		fmt.Println(ebiten.CursorPosition())
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !set_board {
		g.setBoard()
		set_board = true
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
}

func (g *Game) drawBoard(target *ebiten.Image) {
	target.DrawImage(g.board, nil)
	for _, ball := range g.cueBalls {
		vector.DrawFilledCircle(target, ball.cx, ball.cy, circ_Radius, ball.color, antialias)
	}
}

func (g *Game) arrangePyramids(steps, depth, padding, idx int) {
	if steps < 1 {
		return
	}
	baseTriangleX := 900
	baseTriangleY := config.WIN_HEIGHT/3 + (circ_Radius * 2)
	for i := 1; i <= steps; i++ {
		g.cueBalls[i+idx].cx = float32(baseTriangleX) - float32(depth*circ_Radius)
		g.cueBalls[i+idx].cy = float32(
			baseTriangleY,
		) + float32(
			(2*circ_Radius)*i,
		) + float32(
			padding*circ_Radius,
		)
	}
	g.arrangePyramids(steps-1, depth+2, padding+1, steps + idx)
}
