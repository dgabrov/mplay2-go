package endpoint

import (
	"bytes"
	"io"
	"strings"
)

type MediaTypeDetector struct{}

func NewMediaTypeDetector() *MediaTypeDetector {
	return &MediaTypeDetector{}
}

func (d *MediaTypeDetector) DetectMediaType(reader io.Reader) (string, error) {
	signature := make([]byte, 512)
	n, err := reader.Read(signature)
	if err != nil && err != io.EOF {
		return "", err
	}
	signature = signature[:n]

	contentType := d.detectFromSignature(signature)
	return contentType, nil
}

func (d *MediaTypeDetector) IsVideo(mimeType string) bool {
	return strings.HasPrefix(mimeType, "video/")
}

func (d *MediaTypeDetector) detectFromSignaturePublic(sig []byte) (string, bool) {
	contentType := d.detectFromSignature(sig)
	return contentType, contentType != ""
}

func (d *MediaTypeDetector) detectFromSignature(sig []byte) string {
	if len(sig) == 0 {
		return ""
	}

	// MP4/MOV/M4A/M4B (video or audio)
	if len(sig) > 8 && bytes.Equal(sig[4:8], []byte("ftyp")) {
		return d.detectMP4Type(sig)
	}

	// WebM (video or audio)
	if len(sig) > 3 && bytes.Equal(sig[0:4], []byte{0x1A, 0x45, 0xDF, 0xA3}) {
		if bytes.Contains(sig, []byte("V_")) {
			return "video/webm"
		}
		if bytes.Contains(sig, []byte("A_")) {
			return "audio/webm"
		}
		return "video/webm"
	}

	// AVI (video)
	if len(sig) > 8 && bytes.Equal(sig[0:4], []byte("RIFF")) && bytes.Contains(sig, []byte("AVI ")) {
		return "video/x-msvideo"
	}

	// MPEG-TS (video)
	if len(sig) > 0 && sig[0] == 0x47 {
		return "video/mp2t"
	}

	// MPEG-PS/VOB (video)
	if len(sig) > 3 && bytes.Equal(sig[0:3], []byte{0x00, 0x00, 0x01}) {
		return "video/mpeg"
	}

	// FLV (video)
	if len(sig) > 3 && bytes.Equal(sig[0:3], []byte("FLV")) {
		return "video/x-flv"
	}

	// WAV (audio)
	if len(sig) > 8 && bytes.Equal(sig[0:4], []byte("RIFF")) && bytes.Contains(sig, []byte("WAVE")) {
		return "audio/wav"
	}

	// MP3 (audio)
	if len(sig) > 2 && (bytes.Equal(sig[0:2], []byte{0xFF, 0xFB}) ||
		bytes.Equal(sig[0:2], []byte{0xFF, 0xFA})) {
		return "audio/mpeg"
	}

	// ID3 tag (MP3)
	if len(sig) > 3 && bytes.Equal(sig[0:3], []byte("ID3")) {
		return "audio/mpeg"
	}

	// OGG Vorbis/Opus (audio)
	if len(sig) > 3 && bytes.Equal(sig[0:4], []byte("OggS")) {
		if len(sig) > 28 && bytes.Equal(sig[28:35], []byte("OpusHead")) {
			return "audio/opus"
		}
		return "audio/ogg"
	}

	// FLAC (audio)
	if len(sig) > 3 && bytes.Equal(sig[0:4], []byte("fLaC")) {
		return "audio/flac"
	}

	// AAC (audio)
	if len(sig) > 1 && (bytes.Equal(sig[0:2], []byte{0xFF, 0xF1}) ||
		bytes.Equal(sig[0:2], []byte{0xFF, 0xF9})) {
		return "audio/aac"
	}

	return ""
}

func (d *MediaTypeDetector) detectMP4Type(sig []byte) string {
	// Check the ftyp brand to determine if it's audio or video
	if len(sig) > 11 {
		// Check for audio-specific brands
		if bytes.Contains(sig, []byte("M4A")) || bytes.Contains(sig, []byte("M4B")) {
			return "audio/mp4"
		}
	}

	// Default to video for MP4
	return "video/mp4"
}
