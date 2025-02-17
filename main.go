package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"github.com/ebitenui/ebitenui"
	euiimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

var img *image.RGBA

const (
	screenWidth  = 1920
	screenHeight = 1080
)

type Game struct {
	image *image.RGBA
	ui    *ebitenui.UI
}

func (g *Game) Update() error {
	g.ui.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(4, 4)
	i := ebiten.NewImageFromImage(g.image)
	screen.DrawImage(i, op)
	g.ui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func ReadBytes(file *os.File, n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := file.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func SetupUI() *ebitenui.UI {
	rootContainer := widget.NewContainer(
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			//Set how much padding before displaying content
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(30)),
		)),
	)

	innerContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(euiimage.NewNineSliceColor(color.NRGBA{32, 48, 64, 127})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				StretchHorizontal:  false,
				StretchVertical:    true,
			}),
			widget.WidgetOpts.MinSize(600, 100),
		),
	)
	rootContainer.AddChild(innerContainer)

	eui := &ebitenui.UI{
		Container: rootContainer,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("RevenantRE")
	return eui
}

func main() {
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Equip\\scroll.i2d")
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Misc\\bread.i2d")
	file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\resources\\book.dat")
	// file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\resources\\inventory.dat")

	if err != nil {
		log.Fatal(err)
	}

	fr, _ := NewFileResource(file)

	file.Seek(int64(fr.BitmapTable[0]), 1)
	bm, _ := NewBitmap(file)
	fmt.Println(bm.Height, bm.Width)

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
