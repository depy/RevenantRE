package graphics

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/depy/RevenantRE/utils"
)

// Bitmap Flags
const (
	BM_8BIT       uint16 = 0x0001 // Bitmap data is 8 bit.
	BM_15BIT      uint16 = 0x0002 // Bitmap data is 15 bit.
	BM_16BIT      uint16 = 0x0004 // Bitmap data is 16 bit.
	BM_24BIT      uint16 = 0x0008 // Bitmap data is 24 bit.
	BM_32BIT      uint16 = 0x0010 // Bitmap data is 24 bit.
	BM_ZBUFFER    uint16 = 0x0020 // Bitmap has ZBuffer.
	BM_NORMALS    uint16 = 0x0040 // Bitmap has Normal Buffer.
	BM_ALIAS      uint16 = 0x0080 // Bitmap has Alias Buffer.
	BM_ALPHA      uint16 = 0x0100 // Bitmap has Alpha Buffer.
	BM_PALETTE    uint16 = 0x0200 // Bitmap has 256 Color SPalette Structure.
	BM_REGPOINT   uint16 = 0x0400 // Bitmap has registration point
	BM_NOBITMAP   uint16 = 0x0800 // Bitmap has no pixel data
	BM_5BITPAL    uint16 = 0x1000 // Bitmap palette is 5 bit for r,g,b instead of 8 bit
	BM_COMPRESSED uint16 = 0x4000 // Bitmap is compressed.
	BM_CHUNKED    uint32 = 0x8000 // Bitmap is chunked out
)

type BitmapFlags struct {
	Is8bit        bool
	Is15bit       bool
	Is16bit       bool
	Is24bit       bool
	Is32bit       bool
	HasZBuffer    bool
	HasNormals    bool
	HasAlias      bool
	HasAlpha      bool
	HasPalette    bool
	HasRegPoint   bool
	NoBitmap      bool
	Is5bitPalette bool
	IsCompressed  bool
	IsChunked     bool
}

type RGBA struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

type BitmapHeader struct {
	Width         uint32
	Height        uint32
	RegPointX     uint32
	RegPointY     uint32
	Flags         uint32
	DrawingMode   uint32
	KeyColor      uint32
	AliasSize     uint32
	AliasOffset   uint32
	AlphaSize     uint32
	Alpha         uint32
	ZBufferSize   uint32
	ZBuffer       uint32
	NormalSize    uint32
	Normal        uint32
	PaletteSize   uint32
	PaletteOffset uint32
	DataSize      uint32
}

type Bitmap struct {
	Width   uint32
	Height  uint32
	Header  BitmapHeader
	Palette Palette
	Data    []RGBA
}

type Palette struct {
	Colors []RGBA
}

func NewPalette(data []byte) Palette {
	p := Palette{}
	for i := 0; i < len(data); i += 2 {
		c := binary.LittleEndian.Uint16(data[i : i+2])
		pR := uint8((c>>10)&0x1F) * 8
		pG := uint8((c>>5)&0x1F) * 8
		pB := uint8(c&0x1F) * 8
		p.Colors = append(p.Colors, RGBA{pR, pG, pB, 255})
	}
	return p
}

func NewBitmapFlags(flags uint32) BitmapFlags {
	return BitmapFlags{
		Is8bit:        flags&uint32(BM_8BIT) != 0,
		Is15bit:       flags&uint32(BM_15BIT) != 0,
		Is16bit:       flags&uint32(BM_16BIT) != 0,
		Is24bit:       flags&uint32(BM_24BIT) != 0,
		Is32bit:       flags&uint32(BM_32BIT) != 0,
		HasZBuffer:    flags&uint32(BM_ZBUFFER) != 0,
		HasNormals:    flags&uint32(BM_NORMALS) != 0,
		HasAlias:      flags&uint32(BM_ALIAS) != 0,
		HasAlpha:      flags&uint32(BM_ALPHA) != 0,
		HasPalette:    flags&uint32(BM_PALETTE) != 0,
		HasRegPoint:   flags&uint32(BM_REGPOINT) != 0,
		NoBitmap:      flags&uint32(BM_NOBITMAP) != 0,
		Is5bitPalette: flags&uint32(BM_5BITPAL) != 0,
		IsCompressed:  flags&uint32(BM_COMPRESSED) != 0,
		IsChunked:     flags&uint32(BM_CHUNKED) != 0,
	}
}

func PrintBitmapFlags(flags *BitmapFlags) {
	fmt.Println("----- Bitmap Flags -----")
	fmt.Println("Is8bit:\t\t", flags.Is8bit)
	fmt.Println("Is16bit:\t", flags.Is16bit)
	fmt.Println("Is15bit:\t", flags.Is15bit)
	fmt.Println("Is24bit:\t", flags.Is24bit)
	fmt.Println("Is32bit:\t", flags.Is32bit)
	fmt.Println("HasZBuffer:\t", flags.HasZBuffer)
	fmt.Println("HasNormals:\t", flags.HasNormals)
	fmt.Println("HasAlias:\t", flags.HasAlias)
	fmt.Println("HasAlpha:\t", flags.HasAlpha)
	fmt.Println("HasPalette:\t", flags.HasPalette)
	fmt.Println("HasRegPoint:\t", flags.HasRegPoint)
	fmt.Println("NoBitmap:\t", flags.NoBitmap)
	fmt.Println("Is5bitPalette:\t", flags.Is5bitPalette)
	fmt.Println("IsCompressed:\t", flags.IsCompressed)
	fmt.Println("IsChunked:\t", flags.IsChunked)
	fmt.Println("------------------------")
}

func NewBitmapHeader(data []byte) BitmapHeader {
	return BitmapHeader{
		Width:         binary.LittleEndian.Uint32(data[0:4]),
		Height:        binary.LittleEndian.Uint32(data[4:8]),
		RegPointX:     binary.LittleEndian.Uint32(data[8:12]),
		RegPointY:     binary.LittleEndian.Uint32(data[12:16]),
		Flags:         binary.LittleEndian.Uint32(data[16:20]),
		DrawingMode:   binary.LittleEndian.Uint32(data[20:24]),
		KeyColor:      binary.LittleEndian.Uint32(data[24:28]),
		AliasSize:     binary.LittleEndian.Uint32(data[28:32]),
		AliasOffset:   binary.LittleEndian.Uint32(data[32:36]),
		AlphaSize:     binary.LittleEndian.Uint32(data[36:40]),
		Alpha:         binary.LittleEndian.Uint32(data[40:44]),
		ZBufferSize:   binary.LittleEndian.Uint32(data[44:48]),
		ZBuffer:       binary.LittleEndian.Uint32(data[48:52]),
		NormalSize:    binary.LittleEndian.Uint32(data[52:56]),
		Normal:        binary.LittleEndian.Uint32(data[56:60]),
		PaletteSize:   binary.LittleEndian.Uint32(data[60:64]),
		PaletteOffset: binary.LittleEndian.Uint32(data[64:68]),
		DataSize:      binary.LittleEndian.Uint32(data[68:72]),
	}
}

func PrintBitmapHeader(bmh *BitmapHeader) {
	fmt.Println("------ Bitmap Header -----")
	fmt.Println("Width:\t", bmh.Width)
	fmt.Println("Height:\t", bmh.Height)
	fmt.Println("RegPointX:\t", bmh.RegPointX)
	fmt.Println("RegPointY:\t", bmh.RegPointY)
	fmt.Println("Flags:\t", bmh.Flags)
	fmt.Println("DrawingMode:\t", bmh.DrawingMode)
	fmt.Println("KeyColor:\t", bmh.KeyColor)
	fmt.Println("AliasSize:\t", bmh.AliasSize)
	fmt.Println("AliasOffset:\t", bmh.AliasOffset)
	fmt.Println("AlphaSize:\t", bmh.AlphaSize)
	fmt.Println("Alpha:\t", bmh.Alpha)
	fmt.Println("ZBufferSize:\t", bmh.ZBufferSize)
	fmt.Println("ZBuffer:\t", bmh.ZBuffer)
	fmt.Println("NormalSize:\t", bmh.NormalSize)
	fmt.Println("Normal:\t", bmh.Normal)
	fmt.Println("PaletteSize:\t", bmh.PaletteSize)
	fmt.Println("PaletteOffset:\t", bmh.PaletteOffset)
	fmt.Println("DataSize:\t", bmh.DataSize)
	fmt.Println("------------------------")
}

func NewBitmap(file *os.File, readOnlyHeaders bool) (Bitmap, error) {
	bmhData, err := utils.ReadBytes(file, 72) // Seems like the header is 72 bytes when there's no chunking header following
	if err != nil {
		fmt.Println("Error reading bitmap header for: ", file.Name())
		return Bitmap{}, err
	}

	bmHeader := NewBitmapHeader(bmhData)
	//PrintBitmapHeader(&bmHeader)

	bmFlags := NewBitmapFlags(bmHeader.Flags)
	//PrintBitmapFlags(&bmFlags)

	rgbData := []RGBA{}
	palette := Palette{}

	if !readOnlyHeaders {
		bmapData, err := utils.ReadBytes(file, int(bmHeader.DataSize))
		if err != nil {
			fmt.Println("Error reading bitmap data for: ", file.Name())
			return Bitmap{}, err
		}

		if bmFlags.Is15bit {
			rgbData = RenderBitmap15bit(bmHeader, bmapData)
		} else if bmFlags.Is8bit {
			paletteData, err := utils.ReadBytes(file, 512)
			palette = NewPalette(paletteData)
			if err != nil {
				fmt.Println("Error reading palette data for: ", file.Name())
				return Bitmap{}, err
			}
			if bmFlags.IsCompressed {
				chunks := Decompress(bmapData, bmFlags.IsChunked)
				bmHeader.Width = chunks.ChunksHeader.Width * CHUNK_WIDTH
				bmHeader.Height = chunks.ChunksHeader.Height * CHUNK_HEIGHT
				rgbData = RenderChunkedBitmap8bit(bmHeader, chunks, palette)
			} else {
				rgbData = RenderBitmap8bit(bmHeader, bmapData, palette)
			}
		}
	}

	return Bitmap{bmHeader.Width, bmHeader.Height, bmHeader, palette, rgbData}, nil
}

func RenderBitmap15bit(bmh BitmapHeader, data []byte) []RGBA {
	result := make([]RGBA, bmh.Width*bmh.Height)

	for i := range bmh.Height {
		for j := range bmh.Width {
			d := data[i*bmh.Width*2+j*2 : i*bmh.Width*2+j*2+2]

			pixelData := binary.LittleEndian.Uint16(d)
			convPixelData := pixelData
			pR := uint8((convPixelData&0b0111110000000000)>>10) << 3
			pG := uint8((convPixelData&0b0000001111100000)>>5) << 3
			pB := uint8((convPixelData & 0b0000000000011111)) << 3
			pA := uint8(convPixelData & 0b1000000000000000)

			result[i*bmh.Width+j] = RGBA{pR, pG, pB, pA}
		}
	}
	return result
}

func RenderBitmap8bit(bmh BitmapHeader, data []byte, palette Palette) []RGBA {
	result := make([]RGBA, bmh.Width*bmh.Height)
	for i := range bmh.Height {
		for j := range bmh.Width {
			d := data[i*bmh.Width+j]
			c := palette.Colors[d]
			result[i*bmh.Width+j] = RGBA{c.R, c.G, c.B, c.A}
		}
	}
	return result
}

func RenderChunkedBitmap8bit(bmh BitmapHeader, cbd ChunkedBitmapData, palette Palette) []RGBA {
	size := int(cbd.ChunksHeader.Width) * int(cbd.ChunksHeader.Height) * int(CHUNK_WIDTH) * int(CHUNK_HEIGHT)
	result := make([]RGBA, size)

	for i := range cbd.Chunks {
		chunk := cbd.Chunks[i]
		for k := 0; k < CHUNK_HEIGHT; k++ {
			for l := 0; l < CHUNK_WIDTH; l++ {
				xOff := (i % int(cbd.ChunksHeader.Width)) * CHUNK_WIDTH
				yOff := (i / int(cbd.ChunksHeader.Width)) * CHUNK_HEIGHT
				ri := k*CHUNK_HEIGHT*int(cbd.ChunksHeader.Width) + yOff*CHUNK_HEIGHT*int(cbd.ChunksHeader.Width) + l + xOff

				dPos := k*64 + l
				d := chunk.DecompData[dPos]
				c := palette.Colors[d]
				result[ri] = RGBA{c.R, c.G, c.B, c.A}
			}
		}
	}

	return result
}
