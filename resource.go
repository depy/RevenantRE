package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

const FileResourceHeaderSize = 20

type FileResource struct {
	Header      FileResourceHeader
	BitmapTable []uint32
}

type FileResourceHeader struct {
	Magic       uint32
	Topbm       uint16
	CompType    uint8
	Version     uint8
	DataSize    uint32
	ObjSize     uint32
	HeaderSize  uint32
	ImgryHeader ImageryHeader
}

type ImageryHeader struct {
	ImageryId         uint32
	NumStates         uint32
	ImgryStateHeaders []ImageryStateHeader
}

type ImageryStateHeader struct {
	AnimName           [32]byte
	Walkmap            uint32
	Flags              uint32
	Animflags          uint16
	Frames             uint16 // Number of frames
	MaxWidth           uint16 // Graphics maximum width/height (for IsOnScreen and refresh rects)
	MaxHeight          uint16
	RegX               uint16 // Registration point x,y,z for graphics
	RegY               uint16
	RegZ               uint16
	AnimRegx           uint16 // Registration point of animation
	AnimRegy           uint16
	AnimRegz           uint16
	WorldRegX          uint16 // World registration x and y of walk and bounding box info
	WorldRegY          uint16
	WorldRegZ          uint16
	WorldWidth         uint16 // Object's world width, length, and height for walk map and bound box
	WorldLength        uint16
	WorldHeight        uint16
	InventoryAnimFlags uint16 // Animation flags for inventory animation
	InventoryFrames    uint16 // Number of frames of inventory animation
}

func NewFileResourceHeader(data []byte) FileResourceHeader {
	return FileResourceHeader{
		Magic:      binary.LittleEndian.Uint32(data[0:4]),
		Topbm:      binary.LittleEndian.Uint16(data[4:6]),
		CompType:   data[6],
		Version:    data[7],
		DataSize:   binary.LittleEndian.Uint32(data[8:12]),
		ObjSize:    binary.LittleEndian.Uint32(data[12:16]),
		HeaderSize: binary.LittleEndian.Uint32(data[16:20]),
	}
}

func NewImageryHeader(data []byte) ImageryHeader {
	ih := ImageryHeader{
		ImageryId: binary.LittleEndian.Uint32(data[0:4]),
		NumStates: binary.LittleEndian.Uint32(data[4:8]),
	}

	ih.ImgryStateHeaders = []ImageryStateHeader{}
	for i := 0; i < int(ih.NumStates); i++ {
		ish := ImageryStateHeader{}
		copy(ish.AnimName[:], data[0:1])
		// ...
		ih.ImgryStateHeaders = append(ih.ImgryStateHeaders, ish)
	}

	return ih
}

func NewFileResource(file *os.File) (FileResource, error) {
	fileResHdrData, err := ReadBytes(file, 20)
	if err != nil {
		fmt.Println("Error reading resource header")
		return FileResource{}, err
	}

	frh := NewFileResourceHeader(fileResHdrData)

	if frh.HeaderSize > 0 {
		imageryHdr, err := ReadBytes(file, int(frh.HeaderSize))
		if err != nil {
			fmt.Println("Error reading imagery header")
			return FileResource{}, err
		}

		frh.ImgryHeader = NewImageryHeader(imageryHdr)
	}

	bitmapOffsets := []uint32{}
	if frh.Topbm > 0 {
		for range frh.Topbm {
			offset, err := ReadBytes(file, 4)
			if err != nil {
				fmt.Println("Error reading bitmap offset")

			}
			bitmapOffsets = append(bitmapOffsets, binary.LittleEndian.Uint32(offset))
		}
	}
	return FileResource{Header: frh, BitmapTable: bitmapOffsets}, nil
}
