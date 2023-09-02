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
	displacement_v := Vector{x: current_V.x - previous_V.x, y: current_V.y - previous_V.y}

	*previous_V = *current_V

	current_V.x = current_V.x + displacement_v.x + accel_V.x*(config.FIXED_TIMESTAMP_LOOP*config.FIXED_TIMESTAMP_LOOP)
	current_V.y = current_V.y + displacement_v.y + accel_V.y*(config.FIXED_TIMESTAMP_LOOP*config.FIXED_TIMESTAMP_LOOP)
	*accel_V = Vector{}
}

func (g *Game) SetVelocity(id int8, v_v Vector) {
	ball := g.cueBalls[id]
	ball.p_v.x = ball.c_v.x - (v_v.x * config.FIXED_TIMESTAMP_LOOP)
	ball.p_v.y = ball.c_v.y - (v_v.y * config.FIXED_TIMESTAMP_LOOP)
}

// func (g *Game) accumulateForces() {
// 	for _ := range g.cueBalls {
// 		// drag
// 	}
// }

func (g *Game) checkConstraints() {
	for _, cue := range g.cueBalls {
        down := Vector{x: 0, y: config.WIN_HEIGHT - cue.c_v.y}
        up := Vector{x: 0, y: 0 - cue.c_v.y}
        left := Vector{x: 0 - cue.c_v.x, y: 0}
        right := Vector{x: config.WIN_WIDTH - cue.c_v.x, y: 0}

        dist_dwn := down.distance()
        dist_up := up.distance()
        dist_left := left.distance()
        dist_right := right.distance()

        if dist_dwn <= circRadius {
            fmt.Println("BOUNCE UP")
            cue.p_v.y += 1
        }

        if dist_up <= circRadius {
            fmt.Println("BOUNCE DOWN")
            cue.p_v.y -= 1
        }

        if dist_left <= circRadius+5 {
            cue.p_v.x -= 1
        }

        if dist_right <= circRadius {
            cue.p_v.x += 1
        }
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
