package main

import (
	"encoding/binary"
	"fmt"
	"os"
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
	Width  uint32
	Height uint32
	Header *BitmapHeader
	Data   []RGBA
}

func NewBitmapHeader(data []byte) *BitmapHeader {
	return &BitmapHeader{
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

func NewBitmap(file *os.File) (Bitmap, error) {
	bmhData, err := ReadBytes(file, 72) // Seems like the header is 72 bytes when there's no chunking header following
	if err != nil {
		fmt.Println("Error reading bitmap header")
		return Bitmap{}, err
	}

	bmHeader := NewBitmapHeader(bmhData)

	bmapData, err := ReadBytes(file, int(bmHeader.DataSize))
	if err != nil {
		fmt.Println("Error reading bitmap data")
		return Bitmap{}, err
	}

	rgbData := RenderBitmap(bmHeader, bmapData)
	return Bitmap{bmHeader.Width, bmHeader.Height, bmHeader, rgbData}, nil
}

func RenderBitmap(bmh *BitmapHeader, data []byte) []RGBA {
	result := make([]RGBA, bmh.Width*bmh.Height)
	if bmh.Flags&uint32(BM_15BIT) != 0 {
		for i := range bmh.Height {
			for j := range bmh.Width {
				d := data[i*bmh.Width*2+j*2 : i*bmh.Width*2+j*2+2]

				//fmt.Println(hex.Dump(d))
				pixelData := binary.LittleEndian.Uint16(d)
				convPixelData := pixelData
				pR := uint8((convPixelData&0b0111110000000000)>>10) << 3
				pG := uint8((convPixelData&0b0000001111100000)>>5) << 3
				pB := uint8((convPixelData & 0b0000000000011111)) << 3
				pA := uint8(convPixelData & 0b1000000000000000)

				result[i*bmh.Width+j] = RGBA{pR, pG, pB, pA}
			}
		}
	}
	return result
}
