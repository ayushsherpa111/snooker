package game

import (
	"fmt"
	"time"

	config "github.com/ayushsherpa111/snooker/Config"
)

var t = time.Date(2001, time.September, 9, 1, 46, 40, 0, time.UTC)

func (g *Game) move() {
	for _, ball := range g.cueBalls {
		g.verlet(&ball.c_v, &ball.p_v, &ball.a_v)
	}
}

func (g *Game) verlet(current_V, previous_V, accel_V *Vector) {
	temp_V := *current_V

	current_V.x += current_V.x - previous_V.x + accel_V.x*(config.FIXED_TIMESTAMP_LOOP*config.FIXED_TIMESTAMP_LOOP)
	current_V.y += current_V.y - previous_V.y + accel_V.y*(config.FIXED_TIMESTAMP_LOOP*config.FIXED_TIMESTAMP_LOOP)

	*previous_V = temp_V
	fmt.Println(current_V)
	fmt.Println(previous_V)
}

func (g *Game) accumulateForces() {
	for _, i := range g.cueBalls {
		i.a_v = Vector{0, 9.8}
	}
}

// if math.Abs(ball.vx*ball.vx) < 0.8 {
// 	ball.vx = 0
// 	ball.ax = 0
// }
//
// if math.Abs(ball.vy*ball.vy) < 0.8 {
// 	ball.vy = 0
// 	ball.ay = 0
// }
//
// if ball.cx > (config.WIN_WIDTH - circRadius) {
// 	ball.vx *= -1
// 	ball.ax *= -1
// }
//
// if ball.cy > (config.WIN_HEIGHT - circRadius) {
// 	ball.vy *= -1
// 	ball.ay *= -1
// }
//
// if ball.cx < 0 {
// 	ball.vx *= -1
// 	ball.ax *= -1
// }
//
// if ball.cy < 0 {
// 	ball.vy *= -1
// 	ball.ay *= -1
// }
//
// ball.ax = capOff(ball.ax, config.MAX_ACCEL)
// ball.ay = capOff(ball.ay, config.MAX_ACCEL)
