package game

import (
	"fmt"
	"math"
	"time"

	config "github.com/ayushsherpa111/snooker/Config"
)

var t = time.Date(2001, time.September, 9, 1, 46, 40, 0, time.UTC)

func (g *Game) move() bool {
	ballsSkipped := 0
	for _, ball := range g.cueBalls {
		if ball.vx == 0 && ball.vy == 0 {
			ballsSkipped++
			continue
		}
		if g.Debug {
			fmt.Println(ball.ax, ball.ay, ball.vx, ball.vy, ball.cx, ball.cy)
		}

		ball.ax -= ball.vx * config.DRAG
		ball.ay -= ball.vy * config.DRAG

		ball.vx += ball.ax * config.FIXED_TIMESTAMP_LOOP 
		ball.vy += ball.ay * config.FIXED_TIMESTAMP_LOOP 

		ball.cx += ball.vx * config.FIXED_TIMESTAMP_LOOP 
		ball.cy += ball.vy * config.FIXED_TIMESTAMP_LOOP 

		if math.Abs(ball.vx*ball.vx) < 0.8 {
			ball.vx = 0
			ball.ax = 0
		}

		if math.Abs(ball.vy*ball.vy) < 0.8 {
			ball.vy = 0
			ball.ay = 0
		}

		if ball.cx > (config.WIN_WIDTH - circRadius) {
			ball.vx *= -1
			ball.ax *= -1
		}

		if ball.cy > (config.WIN_HEIGHT - circRadius) {
			ball.vy *= -1
			ball.ay *= -1
		}

		if ball.cx < 0 {
			ball.vx *= -1
			ball.ax *= -1
		}

		if ball.cy < 0 {
			ball.vy *= -1
			ball.ay *= -1
		}

		ball.ax = capOff(ball.ax, config.MAX_ACCEL)
		ball.ay = capOff(ball.ay, config.MAX_ACCEL)
	}

	return ballsSkipped != len(g.cueBalls)
}
