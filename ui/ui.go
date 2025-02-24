package ui

import (
	"bytes"
	"image/color"
	"log"

	s "github.com/depy/RevenantRE/state"
	"github.com/ebitenui/ebitenui"
	euiimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

func SetupUI(screenWidth int, screenHeight int, state *s.State) *ebitenui.UI {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(20)),
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
		widget.SliderOpts.Direction(widget.DirectionHorizontal),
		widget.SliderOpts.MinMax(1, 12),

		widget.SliderOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				HorizontalPosition: widget.GridLayoutPositionEnd,
				VerticalPosition:   widget.GridLayoutPositionStart,
			}),
			widget.WidgetOpts.MinSize(270, 25),
		),

		widget.SliderOpts.Images(
			&widget.SliderTrackImage{
				Idle:  euiimage.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: euiimage.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			&widget.ButtonImage{
				Idle:    euiimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Hover:   euiimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Pressed: euiimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
			},
		),
		widget.SliderOpts.FixedHandleSize(6),
		widget.SliderOpts.TrackOffset(0),
		widget.SliderOpts.PageSizeFunc(func() int {
			return 1
		}),
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			state.ImageScalingFactor = float64(args.Current)
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
				HorizontalPosition: widget.GridLayoutPositionStart,
				VerticalPosition:   widget.GridLayoutPositionStart,
			}),
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
