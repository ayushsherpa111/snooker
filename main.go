package main

import (
	"fmt"

	config "github.com/ayushsherpa111/snooker/Config"
	game "github.com/ayushsherpa111/snooker/Game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	snooker := &game.Game{}

	ebiten.SetWindowSize(config.WIN_WIDTH, config.WIN_HEIGHT)
	if err := ebiten.RunGame(snooker); err != nil {
		fmt.Println(err.Error())
		return
	}
}