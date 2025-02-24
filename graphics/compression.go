package graphics

import (
	"encoding/binary"
)

const (
	CHUNK_HEIGHT = 64
	CHUNK_WIDTH  = 64
)

type ChunksHeader struct {
	Type    uint32
	Width   uint32
	Height  uint32
	Offsets []uint32
}

type Chunk struct {
	ChunkId    byte
	RleMarker  byte
	LzMarker   byte
	DecompData []byte
}

type ChunkedBitmapData struct {
	ChunksHeader ChunksHeader
	Chunks       []Chunk
}

func PrintChunksHeader(ch ChunksHeader) {
	println("----- ChunksHeader -----")
	println("Type: ", ch.Type)
	println("Width: ", ch.Width)
	println("Height: ", ch.Height)
	for i, offset := range ch.Offsets {
		print(i, ": ", offset, ", ")
	}
	println("\n--------------------------")
}

func NewChunksHeader(data []byte) ChunksHeader {
	ch := ChunksHeader{}
	ch.Type = binary.LittleEndian.Uint32(data[0:4])
	ch.Width = binary.LittleEndian.Uint32(data[4:8])
	ch.Height = binary.LittleEndian.Uint32(data[8:12])

	for i := range ch.Width * ch.Height {
		chunkOffset := binary.LittleEndian.Uint32(data[12+i*4 : 12+i*4+4]) // 12 is for previously read 3 field (3 * 4)
		ch.Offsets = append(ch.Offsets, chunkOffset)
	}
	return ch
}

func Decompress(data []byte, isChunked bool) ChunkedBitmapData {
	if isChunked {
		return DecompressChunked(data)
	} else {
		return ChunkedBitmapData{}
	}
}

func DecompressChunked(data []byte) ChunkedBitmapData {
	h := NewChunksHeader(data)

	cbd := ChunkedBitmapData{ChunksHeader: h}
	cbd.Chunks = []Chunk{}
	for i := range len(h.Offsets) {
		if h.Offsets[i] != 0 {
			chunkOffsetValueOffset := 12 + 4*i
			ofs := chunkOffsetValueOffset + int(h.Offsets[i])
			chunkStart := data[ofs:]
			chunk := decode(chunkStart)
			cbd.Chunks = append(cbd.Chunks, chunk)
		} else {
			c := Chunk{}
			c.DecompData = make([]byte, CHUNK_HEIGHT*CHUNK_WIDTH)
			cbd.Chunks = append(cbd.Chunks, c)
		}
	}
	return cbd
}

func decode(data []byte) Chunk {
	chunkId := data[0]
	// unknownValue := chunkStart[1:4]
	rleMarker := data[4]
	lzMarker := data[5]
	compData := data[6:]

	chunk := Chunk{
		ChunkId:   chunkId,
		RleMarker: rleMarker,
		LzMarker:  lzMarker,
	}

	chunk.DecompData = make([]byte, CHUNK_HEIGHT*CHUNK_WIDTH)

	si := 0 // Source index
	di := 0 // Destination index
	row := 0
	for row < CHUNK_HEIGHT {
		b := compData[si]
		si++

		if b == rleMarker {
			count := compData[si]
			si++

			if count == 0 {
				// End of line
				row++
				continue
			} else if count < 0x80 {
				// Normal RLE mode
				count &= 0x7f
				val := compData[si]
				si++
				for i := 0; i < int(count); i++ {
					chunk.DecompData[di] = val
					di++
				}
			} else if count >= 0x80 {
				// Skip mode
				count &= 0x7f
				di += int(count)
			}
		} else if b == lzMarker {
			count := compData[si]
			si++
			offset := binary.LittleEndian.Uint16(compData[si : si+2])
			si += 2
			lzOffset := di - int(offset) - 4
			for i := 0; i < int(count); i++ {
				chunk.DecompData[di] = chunk.DecompData[lzOffset]
				di++
				lzOffset++
			}
		} else {
			chunk.DecompData[di] = b
			di++
		}
	}
	return chunk
}
