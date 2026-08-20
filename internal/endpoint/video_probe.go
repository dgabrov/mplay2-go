package endpoint

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
)

type VideoDimensions struct {
	Width  int
	Height int
	Valid  bool
}

type VideoProber struct{}

func NewVideoProber() *VideoProber {
	return &VideoProber{}
}

func (p *VideoProber) ProbeMP4Dimensions(reader io.Reader) VideoDimensions {
	data := make([]byte, 4*1024*1024)
	n, err := reader.Read(data)
	if err != nil && err != io.EOF {
		return VideoDimensions{}
	}
	data = data[:n]

	return p.parseMP4Dimensions(data)
}

func (p *VideoProber) parseMP4Dimensions(data []byte) VideoDimensions {
	pos := 0
	for pos < len(data)-8 {
		boxSize := binary.BigEndian.Uint32(data[pos : pos+4])
		boxType := string(data[pos+4 : pos+8])

		if boxSize == 0 || boxSize > uint32(len(data)-pos) {
			break
		}

		if boxType == "moov" {
			return p.parseMoovBox(data, pos+8, pos+int(boxSize))
		}

		pos += int(boxSize)
		if pos < 0 {
			break
		}
	}

	return VideoDimensions{}
}

func (p *VideoProber) parseMoovBox(data []byte, start, end int) VideoDimensions {
	pos := start
	for pos < end-8 {
		if pos+8 > len(data) {
			break
		}

		boxSize := binary.BigEndian.Uint32(data[pos : pos+4])
		boxType := string(data[pos+4 : pos+8])

		if boxSize == 0 || int(boxSize) > end-pos {
			break
		}

		if boxType == "trak" {
			dims := p.parseTrakBox(data, pos+8, pos+int(boxSize))
			if dims.Valid {
				return dims
			}
		}

		pos += int(boxSize)
		if pos < 0 {
			break
		}
	}

	return VideoDimensions{}
}

func (p *VideoProber) parseTrakBox(data []byte, start, end int) VideoDimensions {
	pos := start
	for pos < end-8 {
		if pos+8 > len(data) {
			break
		}

		boxSize := binary.BigEndian.Uint32(data[pos : pos+4])
		boxType := string(data[pos+4 : pos+8])

		if boxSize == 0 || int(boxSize) > end-pos {
			break
		}

		if boxType == "tkhd" {
			return p.parseTkhdBox(data, pos+8, pos+int(boxSize))
		}

		pos += int(boxSize)
		if pos < 0 {
			break
		}
	}

	return VideoDimensions{}
}

func (p *VideoProber) parseTkhdBox(data []byte, start, end int) VideoDimensions {
	if end-start < 80 {
		return VideoDimensions{}
	}

	version := data[start]
	var widthPos, heightPos int

	if version == 0 {
		widthPos = start + 76
		heightPos = start + 80
	} else {
		widthPos = start + 84
		heightPos = start + 88
	}

	if heightPos+4 > len(data) {
		return VideoDimensions{}
	}

	widthFixed := binary.BigEndian.Uint32(data[widthPos : widthPos+4])
	heightFixed := binary.BigEndian.Uint32(data[heightPos : heightPos+4])

	width := int(widthFixed >> 16)
	height := int(heightFixed >> 16)

	if width > 0 && height > 0 {
		return VideoDimensions{
			Width:  width,
			Height: height,
			Valid:  true,
		}
	}

	return VideoDimensions{}
}

func (p *VideoProber) ProbeWebMDimensions(reader io.Reader) VideoDimensions {
	data := make([]byte, 4*1024*1024)
	n, err := reader.Read(data)
	if err != nil && err != io.EOF {
		return VideoDimensions{}
	}
	data = data[:n]

	return p.parseWebMDimensions(data)
}

func (p *VideoProber) parseWebMDimensions(data []byte) VideoDimensions {
	if len(data) < 4 || !bytes.Equal(data[0:4], []byte{0x1A, 0x45, 0xDF, 0xA3}) {
		return VideoDimensions{}
	}

	pos := 4
	for pos < len(data)-2 {
		elementID := p.readVint(data[pos:])
		if elementID.length == 0 {
			break
		}
		pos += elementID.length

		if pos >= len(data) {
			break
		}

		size := p.readVint(data[pos:])
		if size.length == 0 {
			break
		}
		pos += size.length

		if pos+size.value > len(data) {
			break
		}

		if bytes.Equal(elementID.data[:elementID.length], []byte{0xAE}) {
			segment := data[pos : pos+size.value]
			dims := p.parseSegment(segment)
			if dims.Valid {
				return dims
			}
		}

		pos += size.value
	}

	return VideoDimensions{}
}

type vint struct {
	data   [8]byte
	value  int
	length int
}

func (p *VideoProber) readVint(data []byte) vint {
	if len(data) == 0 {
		return vint{}
	}

	firstByte := data[0]
	length := 0
	mask := byte(0x80)

	for i := 0; i < 8; i++ {
		if firstByte&mask != 0 {
			length = i + 1
			break
		}
		mask >>= 1
	}

	if length == 0 || length > len(data) {
		return vint{}
	}

	result := vint{
		length: length,
	}
	copy(result.data[:], data[:length])

	value := int(firstByte &^ mask)
	for i := 1; i < length; i++ {
		if i < len(data) {
			value = (value << 8) | int(data[i])
		}
	}

	result.value = value
	return result
}

func (p *VideoProber) parseSegment(data []byte) VideoDimensions {
	pos := 0
	for pos < len(data)-2 {
		elementID := p.readVint(data[pos:])
		if elementID.length == 0 {
			break
		}
		pos += elementID.length

		if pos >= len(data) {
			break
		}

		size := p.readVint(data[pos:])
		if size.length == 0 {
			break
		}
		pos += size.length

		if pos+size.value > len(data) {
			break
		}

		if bytes.Equal(elementID.data[:elementID.length], []byte{0xA0}) {
			segment := data[pos : pos+size.value]
			dims := p.parseTrack(segment)
			if dims.Valid {
				return dims
			}
		}

		pos += size.value
	}

	return VideoDimensions{}
}

func (p *VideoProber) parseTrack(data []byte) VideoDimensions {
	pos := 0
	for pos < len(data)-2 {
		elementID := p.readVint(data[pos:])
		if elementID.length == 0 {
			break
		}
		pos += elementID.length

		if pos >= len(data) {
			break
		}

		size := p.readVint(data[pos:])
		if size.length == 0 {
			break
		}
		pos += size.length

		if pos+size.value > len(data) {
			break
		}

		if bytes.Equal(elementID.data[:elementID.length], []byte{0xE0}) {
			segment := data[pos : pos+size.value]
			dims := p.parseTrackVideo(segment)
			if dims.Valid {
				return dims
			}
		}

		pos += size.value
	}

	return VideoDimensions{}
}

func (p *VideoProber) parseTrackVideo(data []byte) VideoDimensions {
	pos := 0
	for pos < len(data)-2 {
		elementID := p.readVint(data[pos:])
		if elementID.length == 0 {
			break
		}
		pos += elementID.length

		if pos >= len(data) {
			break
		}

		size := p.readVint(data[pos:])
		if size.length == 0 {
			break
		}
		pos += size.length

		if pos+size.value > len(data) {
			break
		}

		if bytes.Equal(elementID.data[:elementID.length], []byte{0xB0}) {
			width, _ := p.readFloat(data[pos : pos+size.value])
			if width == 0 {
				pos += size.value
				continue
			}

			pos += size.value

			if pos < len(data)-2 {
				elementID = p.readVint(data[pos:])
				if elementID.length > 0 {
					pos += elementID.length
					if pos < len(data) {
						size = p.readVint(data[pos:])
						if size.length > 0 && pos+size.length+size.value <= len(data) {
							pos += size.length
							height, _ := p.readFloat(data[pos : pos+size.value])
							if height > 0 {
								return VideoDimensions{
									Width:  int(width),
									Height: int(height),
									Valid:  true,
								}
							}
						}
					}
				}
			}
		}

		pos += size.value
	}

	return VideoDimensions{}
}

func (p *VideoProber) readFloat(data []byte) (float64, bool) {
	if len(data) >= 8 {
		bits := binary.BigEndian.Uint64(data[:8])
		return math.Float64frombits(bits), true
	}
	if len(data) >= 4 {
		bits := binary.BigEndian.Uint32(data[:4])
		return float64(math.Float32frombits(bits)), true
	}
	return 0, false
}
