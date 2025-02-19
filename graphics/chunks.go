package graphics

import "encoding/binary"

type ChunksHeader struct {
	Type    uint32
	Width   uint32
	Height  uint32
	Offsets []uint32
}

func NewChunksHeader(data []byte) ChunksHeader {
	ch := ChunksHeader{}
	ch.Type = binary.LittleEndian.Uint32(data[0:4])
	ch.Width = binary.LittleEndian.Uint32(data[4:8])
	ch.Height = binary.LittleEndian.Uint32(data[8:12])

	for i := range ch.Width * ch.Height {
		chunkOffset := binary.LittleEndian.Uint32(data[12+i*4 : 12+i*4+4])
		ch.Offsets = append(ch.Offsets, chunkOffset)
	}
	return ch
}
