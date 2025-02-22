package entities

import "github.com/hajimehoshi/ebiten/v2"

type Sprite struct {
	Image *ebiten.Image
	X     float64
	Y     float64
}
