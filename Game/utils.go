package game

import "math"

func findDistance(x1, y1, x2, y2 float64) (distance float64) {
	distance = math.Sqrt((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2))
	return
}
