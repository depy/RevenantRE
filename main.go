package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"github.com/depy/RevenantRE/graphics"
	"github.com/depy/RevenantRE/ui"
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

func (g *Game) Update() error {
	scaleFactor = 2
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
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Cave\\cavbones1.i2d") // 8bit, zbuffer, compressed, chunked
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Forest\\forbirch001.i2d") // 8bit, zbuffer, compressed, chunked
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Misc\\dragonent.i2d") // 8bit, zbuffer, compressed, chunked
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Forest\\formushrooms2.i2d") // 8bit, zbuffer, compressed, chunked
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\KeepInt\\kinrug.i2d") // 8bit, zbuffer, compressed, chunked
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\resources\\book.dat") // Works
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\resources\\scroll.dat") // Works
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Equip\\scroll.i2d") // Works
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Misc\\book.i2d") // Works
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Magic\\death.i2d") // Works
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Misc\\cheese.i2d") // Works
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Misc\\potionblue.i2d") // Works
	file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Misc\\bread.i2d") // Works

	if err != nil {
		log.Fatal(err)
	}

	fr, err := graphics.NewFileResource(file, false)
	if err != nil {
		log.Fatal(err)
		return
	}

	bm := fr.Bitmaps[0]

	img = image.NewRGBA(image.Rect(0, 0, int(bm.Width), int(bm.Height)))
	for i := 0; i < len(bm.Data); i++ {
		x := i % int(bm.Width)
		y := i / int(bm.Width)
		c := color.RGBA{bm.Data[i].R, bm.Data[i].G, bm.Data[i].B, bm.Data[i].A}
		img.Set(x, y, c)
	}

	eui := ui.SetupUI(screenWidth, screenHeight)
	g := &Game{
		image: img,
		ui:    eui,
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
