package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"golang.org/x/image/font/gofont/goregular"

	"github.com/ebitenui/ebitenui"
	euiimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var debug bool = false

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

func ReadNByes(file *os.File, n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := file.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

var img *image.RGBA

func main() {
	f := flag.String("debug", "", "Debug mode")
	flag.Parse()
	if *f != "" {
		debug = true
	}

	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Equip\\scroll.i2d")
	file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\imagery\\Imagery\\Misc\\bread.i2d")
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\resources\\book.dat")
	//file, err := os.Open("D:\\Games\\RevenantRE\\__1extracted\\resources\\inventory.dat")

	if err != nil {
		log.Fatal(err)
	}

	mainHeader, err := ReadNByes(file, 20)
	if err != nil {
		fmt.Println("Error reading main header")
	}

	mhMagic := mainHeader[0:4]
	mhTopbm := int(binary.LittleEndian.Uint16(mainHeader[4:6]))
	mhCompType := mainHeader[6:7][0]
	mhVersion := mainHeader[7:8][0]
	mgDataSize := int(binary.LittleEndian.Uint32(mainHeader[8:12]))
	mhObjSize := int(binary.LittleEndian.Uint32(mainHeader[12:16]))
	mhHeaderSize := int(binary.LittleEndian.Uint32(mainHeader[16:20]))

	if debug {
		fmt.Println(hex.Dump(mainHeader))
	}

	fmt.Println("---------- Base header ----------")
	fmt.Println("Magic:\t\t", mhMagic)
	fmt.Println("Topbm:\t\t", mhTopbm)
	fmt.Println("CompType:\t", mhCompType)
	fmt.Println("Version:\t", mhVersion)
	fmt.Println("DataSize:\t", mgDataSize)
	fmt.Println("ObjSize:\t", mhObjSize)
	fmt.Println("HeaderSize:\t", mhHeaderSize)
	fmt.Println("------------------------------\n")

	if mhHeaderSize > 0 {
		resHeader, err := ReadNByes(file, mhHeaderSize)
		if err != nil {
			fmt.Println("Error reading resource header")
		}
		if debug {
			fmt.Println(hex.Dump(resHeader))
		}
	}

	bitmapOffsets := []int{}
	if mhTopbm > 0 {
		for range mhTopbm {
			offset, err := ReadNByes(file, 4)
			if err != nil {
				fmt.Println("Error reading bitmap offset")

			}
			bitmapOffsets = append(bitmapOffsets, int(binary.LittleEndian.Uint32(offset)))
		}
	}
	if debug {
		fmt.Println("== Bitmap offsets: %v\n", bitmapOffsets)
	}

	file.Seek(int64(bitmapOffsets[0]), 1)

	bmapHeader, err := ReadNByes(file, 72) // Seems like the header is 72 bytes when there's no chunking header following
	if err != nil {
		fmt.Println("Error reading bitmap header")
	}

	bmWidth := int(binary.LittleEndian.Uint32(bmapHeader[0:4]))
	bmHeight := int(binary.LittleEndian.Uint32(bmapHeader[4:8]))
	bmRegPointX := int(binary.LittleEndian.Uint32(bmapHeader[8:12]))
	bmRegPointY := int(binary.LittleEndian.Uint32(bmapHeader[12:16]))
	bmFlags := int(binary.LittleEndian.Uint32(bmapHeader[16:20]))
	bmDrawingMode := int(binary.LittleEndian.Uint32(bmapHeader[20:24]))
	bmKeyColor := int(binary.LittleEndian.Uint32(bmapHeader[24:28]))
	bmAliasSize := int(binary.LittleEndian.Uint32(bmapHeader[28:32]))
	bmAliasOffset := int(binary.LittleEndian.Uint32(bmapHeader[32:36]))
	bmAlphaSize := int(binary.LittleEndian.Uint32(bmapHeader[36:40]))
	bmAlpha := int(binary.LittleEndian.Uint32(bmapHeader[40:44]))
	bmZBufferSize := int(binary.LittleEndian.Uint32(bmapHeader[44:48]))
	bmZBuffer := int(binary.LittleEndian.Uint32(bmapHeader[48:52]))
	bmNormalSize := int(binary.LittleEndian.Uint32(bmapHeader[52:56]))
	bmNormal := int(binary.LittleEndian.Uint32(bmapHeader[56:60]))
	bmPaletteSize := int(binary.LittleEndian.Uint32(bmapHeader[60:64]))
	bmPaletteOffset := int(binary.LittleEndian.Uint32(bmapHeader[64:68]))
	bmDataSize := int(binary.LittleEndian.Uint32(bmapHeader[68:72]))
	// chunkDecompFlag := int(binary.LittleEndian.Uint32(bmapHeader[72:76]))
	// chunkWidth := int(binary.LittleEndian.Uint32(bmapHeader[76:80]))
	// chunkHeight := int(binary.LittleEndian.Uint32(bmapHeader[80:84]))

	fmt.Println("---------- Bitmap header ----------")
	fmt.Println("Width:\t\t", bmWidth)
	fmt.Println("Height:\t\t", bmHeight)
	fmt.Println("RegPointX:\t", bmRegPointX)
	fmt.Println("RegPointY:\t", bmRegPointY)
	fmt.Println("Flags:\t\t", bmFlags)
	fmt.Println("DrawingMode:\t", bmDrawingMode)
	fmt.Println("KeyColor:\t", bmKeyColor)
	fmt.Println("AliasSize:\t", bmAliasSize)
	fmt.Println("AliasOffset:\t", bmAliasOffset)
	fmt.Println("AlphaSize:\t", bmAlphaSize)
	fmt.Println("Alpha:\t\t", bmAlpha)
	fmt.Println("ZBufferSize:\t", bmZBufferSize)
	fmt.Println("ZBuffer:\t", bmZBuffer)
	fmt.Println("NormalSize:\t", bmNormalSize)
	fmt.Println("Normal:\t\t", bmNormal)
	fmt.Println("PaletteSize:\t", bmPaletteSize)
	fmt.Println("Palette:\t", bmPaletteOffset)
	fmt.Println("DataSize:\t", bmDataSize)
	fmt.Println("DataSize:\t", bmDataSize)
	// fmt.Println("ChunkDecompFlag:\t", chunkDecompFlag)
	// fmt.Println("ChunkWidth:\t", chunkWidth)
	// fmt.Println("ChunkHeight:\t", chunkHeight)
	fmt.Println("------------------------------\n")

	if debug {
		fmt.Println(hex.Dump(bmapHeader))
	}

	img = image.NewRGBA(image.Rect(0, 0, bmWidth, bmHeight))

	// 15bit data (2bytes per pixel)
	// for i := range bmHeight {
	// 	for j := range bmWidth {
	// 		d, err := ReadNByes(file, 2)
	// 		if err != nil {
	// 			fmt.Println("Error reading pixel data")
	// 		}
	// 		//fmt.Println(hex.Dump(d))
	// 		pixelData := binary.LittleEndian.Uint16(d)
	// 		convPixelData := pixelData
	// 		pR := uint8((convPixelData&0b0111110000000000)>>10) << 3
	// 		pG := uint8((convPixelData&0b0000001111100000)>>5) << 3
	// 		pB := uint8((convPixelData & 0b0000000000011111)) << 3
	// 		pA := uint8(convPixelData & 0b1000000000000000)

	// 		img.Pix[i*bmWidth*4+j*4+0] = pR
	// 		img.Pix[i*bmWidth*4+j*4+1] = pG
	// 		img.Pix[i*bmWidth*4+j*4+2] = pB
	// 		img.Pix[i*bmWidth*4+j*4+3] = pA

	// 		// fmt.Printf("Raw pixelData: %d\t%b\t", pixelData, pixelData)
	// 		// fmt.Printf("\tCnv pixelData: %d\t%b", convPixelData, convPixelData)
	// 		// fmt.Printf("\tPixel: %d\t%d\t%d\t%d\t%b\t%b\t%b\t%b\n", pR, pG, pB, pA, pR, pG, pB, pA)
	// 	}
	// }

	// Reading 8 bit data (1 byte per pixel) - i2d files
	bitmapData, err := ReadNByes(file, bmWidth*bmHeight)
	palette, err := ReadNByes(file, bmPaletteSize)

	if err != nil {
		fmt.Println("Error reading pixel data")
	}
	for i := range bmHeight {
		for j := range bmWidth {
			d := bitmapData[i*bmWidth+j]
			if err != nil {
				fmt.Println("Error reading pixel data")
			}
			colorData := binary.LittleEndian.Uint16(palette[d : d+2])

			img.Pix[i*bmWidth*4+j*4+0] = uint8((colorData&0b1111100000000000)>>11) << 3
			img.Pix[i*bmWidth*4+j*4+1] = uint8((colorData&0b0000011111100000)>>5) << 3
			img.Pix[i*bmWidth*4+j*4+2] = uint8(colorData&0b0000000000011111) << 3
			img.Pix[i*bmWidth*4+j*4+3] = 0

			// fmt.Printf("Raw pixelData: %d\t%b\t", pixelData, pixelData)
			// fmt.Printf("\tCnv pixelData: %d\t%b", convPixelData, convPixelData)
			// fmt.Printf("\tPixel: %d\t%d\t%d\t%d\t%b\t%b\t%b\t%b\n", pR, pG, pB, pA, pR, pG, pB, pA)
		}
	}

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
	g := &Game{
		image: img,
		ui:    eui,
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func loadButtonImage() (*widget.ButtonImage, error) {
	idle := euiimage.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255})

	hover := euiimage.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255})

	pressed := euiimage.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}

func loadFont(size float64) (text.Face, error) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}
