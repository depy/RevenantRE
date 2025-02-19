package ui

import (
	"bytes"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	euiimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

func SetupUI(screenWidth int, screenHeight int) *ebitenui.UI {
	rootContainer := widget.NewContainer(
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			//Set how much padding before displaying content
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(30)),
		)),
	)

	rightSidePanel := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Spacing(20, 10),
				widget.GridLayoutOpts.Stretch([]bool{false}, []bool{false}),
				widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(10)),
			),
		),
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

	slider := widget.NewSlider(
		// Set the slider orientation - n/s vs e/w
		widget.SliderOpts.Direction(widget.DirectionHorizontal),
		// Set the minimum and maximum value for the slider
		widget.SliderOpts.MinMax(1, 12),

		widget.SliderOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				HorizontalPosition: widget.GridLayoutPositionEnd,
				VerticalPosition:   widget.GridLayoutPositionStart,
			}),
			widget.WidgetOpts.MinSize(270, 25),
		),

		widget.SliderOpts.Images(
			// Set the track images
			&widget.SliderTrackImage{
				Idle:  euiimage.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: euiimage.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			// Set the handle images
			&widget.ButtonImage{
				Idle:    euiimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Hover:   euiimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Pressed: euiimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
			},
		),
		// Set the size of the handle
		widget.SliderOpts.FixedHandleSize(6),
		// Set the offset to display the track
		widget.SliderOpts.TrackOffset(0),
		// Set the size to move the handle
		widget.SliderOpts.PageSizeFunc(func() int {
			return 1
		}),
		// Set the callback to call when the slider value is changed
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			// fmt.Println(args.Current, "dragging", args.Dragging)
		}),
	)

	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
	}

	fontFace := &text.GoTextFace{
		Source: s,
		Size:   18,
	}

	sliderLabel := widget.NewText(
		widget.TextOpts.Text("Scaling factor", fontFace, color.White),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				HorizontalPosition: widget.GridLayoutPositionEnd,
				VerticalPosition:   widget.GridLayoutPositionStart,
			}),
			widget.WidgetOpts.MinSize(280, 25),
		),
	)

	rightSidePanel.AddChild(sliderLabel, slider)
	rootContainer.AddChild(rightSidePanel)

	eui := &ebitenui.UI{
		Container: rootContainer,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("RevenantRE")
	return eui
}
