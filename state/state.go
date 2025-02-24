package state

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type State struct {
	ImageScalingFactor float64
	Image              *ebiten.Image
	Palette            *ebiten.Image
}
