package endpoint

import (
	"bytes"
	"io"
)

const (
	ContentTypeAudio = "audio"
	ContentTypeVideo = "video"
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

func (d *MediaTypeDetector) detectFromSignaturePublic(sig []byte) (string, bool) {
	contentType := d.detectFromSignature(sig)
	return contentType, contentType != ""
}

func (d *MediaTypeDetector) detectFromSignature(sig []byte) string {
	if len(sig) == 0 {
		return ""
	}

	// MP4/MOV (video)
	if len(sig) > 8 && bytes.Equal(sig[4:8], []byte("ftyp")) {
		return ContentTypeVideo
	}

	// WebM/Matroska (video)
	if len(sig) > 3 && bytes.Equal(sig[0:4], []byte{0x1A, 0x45, 0xDF, 0xA3}) {
		return ContentTypeVideo
	}

	// AVI (video)
	if len(sig) > 8 && bytes.Equal(sig[0:4], []byte("RIFF")) && bytes.Contains(sig, []byte("AVI ")) {
		return ContentTypeVideo
	}

	// MPEG-TS (video)
	if len(sig) > 0 && sig[0] == 0x47 {
		return ContentTypeVideo
	}

	// MPEG-PS/VOB (video)
	if len(sig) > 3 && bytes.Equal(sig[0:3], []byte{0x00, 0x00, 0x01}) {
		return ContentTypeVideo
	}

	// FLV (video)
	if len(sig) > 3 && bytes.Equal(sig[0:3], []byte("FLV")) {
		return ContentTypeVideo
	}

	// WAV (audio)
	if len(sig) > 8 && bytes.Equal(sig[0:4], []byte("RIFF")) && bytes.Contains(sig, []byte("WAVE")) {
		return ContentTypeAudio
	}

	// MP3 (audio)
	if len(sig) > 2 && (bytes.Equal(sig[0:2], []byte{0xFF, 0xFB}) ||
		bytes.Equal(sig[0:2], []byte{0xFF, 0xFA}) ||
		bytes.Equal(sig[0:3], []byte("ID3"))) {
		return ContentTypeAudio
	}

	// OGG (could be audio or video, treat as audio by default)
	if len(sig) > 3 && bytes.Equal(sig[0:4], []byte("OggS")) {
		return ContentTypeAudio
	}

	// FLAC (audio)
	if len(sig) > 3 && bytes.Equal(sig[0:4], []byte("fLaC")) {
		return ContentTypeAudio
	}

	// AAC (audio)
	if len(sig) > 1 && bytes.Equal(sig[0:2], []byte{0xFF, 0xF1}) ||
		bytes.Equal(sig[0:2], []byte{0xFF, 0xF9}) {
		return ContentTypeAudio
	}

	// M4A/M4B (audio) - similar to MP4 but audio
	if len(sig) > 8 && bytes.Equal(sig[4:8], []byte("ftyp")) &&
		(bytes.Contains(sig, []byte("M4A")) || bytes.Contains(sig, []byte("M4B"))) {
		return ContentTypeAudio
	}

	return ""
}
