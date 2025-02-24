package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"github.com/depy/RevenantRE/graphics"
	s "github.com/depy/RevenantRE/state"
	"github.com/depy/RevenantRE/ui"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
)

var state *s.State = &s.State{ImageScalingFactor: 1}

const (
	screenWidth  = 1920
	screenHeight = 1080
)

type GameState struct {
	ImageScalingFactor float64
}

type Game struct {
	ui *ebitenui.UI
}

func (g *Game) Update() error {
	g.ui.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(state.ImageScalingFactor, state.ImageScalingFactor)
	screen.DrawImage(state.Image, op)

	g.ui.Draw(screen)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(35, 35)
	op.GeoM.Translate(1320, 480)
	screen.DrawImage(state.Palette, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Cave\\RockPile.i2d")
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Cave\\cavbones1.i2d")
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Misc\\dragonent.i2d")
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\KeepInt\\kinrug.i2d")
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\resources\\book.dat")
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\resources\\scroll.dat")
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Equip\\scroll.i2d")
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Misc\\book.i2d")
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Magic\\death.i2d")
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Misc\\cheese.i2d")
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Misc\\potionblue.i2d")
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Misc\\bread.i2d")
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Forest\\formushrooms2.i2d")
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Forest\\forbirch001.i2d")
	file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Dungeon\\dunalcovee.i2d") // has weird artifacts

	if err != nil {
		log.Fatal(err)
	}

	fr, err := graphics.NewFileResource(file, false)
	if err != nil {
		log.Fatal(err)
		return
	}

	bm := fr.Bitmaps[0]

	img := image.NewRGBA(image.Rect(0, 0, int(bm.Width), int(bm.Height)))
	for i := 0; i < len(bm.Data); i++ {
		x := i % int(bm.Width)
		y := i / int(bm.Width)
		c := color.RGBA{bm.Data[i].R, bm.Data[i].G, bm.Data[i].B, bm.Data[i].A}
		img.Set(x, y, c)
	}

	eImg := ebiten.NewImageFromImage(img)

	palette := image.NewRGBA(image.Rect(0, 0, 16, 16))
	for i := 0; i < len(bm.Palette.Colors); i++ {
		x := i % 16
		y := i / 16
		c := bm.Palette.Colors[i]
		palette.Set(x, y, color.RGBA{c.R, c.G, c.B, c.A})
	}

	ePalette := ebiten.NewImageFromImage(palette)

	state.Image = eImg
	state.Palette = ePalette

	eui := ui.SetupUI(screenWidth, screenHeight, state)
	g := &Game{
		ui: eui,
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
