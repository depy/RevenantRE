package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

var img *image.RGBA
var scaleFactor = 4.0
var slider *widget.Slider

const (
	screenWidth  = 1920
	screenHeight = 1080
)

type Game struct {
	image *image.RGBA
	ui    *ebitenui.UI
}

func ReadBytes(file *os.File, n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := file.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (g *Game) Update() error {
	scaleFactor = float64(slider.Current)
	g.ui.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleFactor, scaleFactor)
	i := ebiten.NewImageFromImage(g.image)
	screen.DrawImage(i, op)
	g.ui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Cave\\cavbones1.i2d")
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Forest\\forbirch001.i2d")
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Misc\\dragonent.i2d")
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Equip\\scroll.i2d")
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Misc\\bread.i2d")
	file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\resources\\book.dat")
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\resources\\inventory.dat")

	if err != nil {
		log.Fatal(err)
	}

	fr, _ := NewFileResource(file)
	bm := fr.Bitmaps[0]

	img = image.NewRGBA(image.Rect(0, 0, int(bm.Width), int(bm.Height)))
	for i := 0; i < len(bm.Data); i++ {
		x := i % int(bm.Width)
		y := i / int(bm.Width)
		c := color.RGBA{bm.Data[i].R, bm.Data[i].G, bm.Data[i].B, bm.Data[i].A}
		img.Set(x, y, c)
	}

	eui := SetupUI()
	g := &Game{
		image: img,
		ui:    eui,
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
