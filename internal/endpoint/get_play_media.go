package endpoint

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/amanagement24/mplay2-go/internal/service"
)

type GetPlayMediaEndpoint struct {
	servr         *service.Servr
	uploadsFolder string
	mediaSlice    int64
}

func NewGetPlayMediaEndpoint(servr *service.Servr, uploadsFolder string, mediaSlice int) *GetPlayMediaEndpoint {
	return &GetPlayMediaEndpoint{
		servr:         servr,
		uploadsFolder: uploadsFolder,
		mediaSlice:    int64(mediaSlice),
	}
}

func (e *GetPlayMediaEndpoint) Handle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	err := e.process(ctx, r, w)
	if err != nil {
		writeJsonResponse(w, nil, err)
	}

	return nil
}

func (e *GetPlayMediaEndpoint) process(ctx context.Context, r *http.Request, w http.ResponseWriter) error {
	token, err := getTokenFromRequest(r)
	if err != nil {
		return fmt.Errorf("user not logged in")
	}

	userID, err := e.servr.ValidateToken(ctx, token)
	if err != nil {
		return err
	}

	mediaID := r.URL.Query().Get("id")
	if mediaID == "" {
		return fmt.Errorf("missing 'id' query parameter")
	}

	media, err := e.servr.GetMedia(ctx, userID, mediaID)
	if err != nil {
		return err
	}

	if media == nil {
		return fmt.Errorf("media not found")
	}

	filePath := filepath.Join(e.uploadsFolder, mediaID+".dat")
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open media file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat media file: %w", err)
	}

	fileSize := fileInfo.Size()

	rangeHeader := r.Header.Get("Range")
	start, end := e.parseRangeHeader(rangeHeader, fileSize)

	if start < 0 || end < start || start >= fileSize {
		w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", fileSize))
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return nil
	}

	if end >= fileSize {
		end = fileSize - 1
	}

	sliceSize := end - start + 1
	if sliceSize > e.mediaSlice {
		sliceSize = e.mediaSlice
		end = start + sliceSize - 1
	}

	if _, err := file.Seek(start, 0); err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	w.Header().Set("Content-Type", media.ContentType)
	w.Header().Set("Content-Length", strconv.FormatInt(sliceSize, 10))
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	w.Header().Set("Accept-Ranges", "bytes")
	w.WriteHeader(http.StatusPartialContent)

	written, err := io.CopyN(w, file, sliceSize)
	if err != nil && err != io.EOF {
		slog.Error("failed to write media slice", "error", err)
		return err
	}

	if written != sliceSize {
		slog.Warn("incomplete media slice written", "expected", sliceSize, "written", written)
	}

	return nil
}

func (e *GetPlayMediaEndpoint) parseRangeHeader(rangeHeader string, fileSize int64) (int64, int64) {
	if rangeHeader == "" {
		return 0, fileSize - 1
	}

	if !strings.HasPrefix(rangeHeader, "bytes=") {
		return 0, fileSize - 1
	}

	rangeSpec := strings.TrimPrefix(rangeHeader, "bytes=")

	if strings.Contains(rangeSpec, ",") {
		return 0, fileSize - 1
	}

	parts := strings.Split(rangeSpec, "-")
	if len(parts) != 2 {
		return 0, fileSize - 1
	}

	var start, end int64

	if parts[0] == "" {
		if parts[1] == "" {
			return 0, fileSize - 1
		}

		suffix, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil || suffix <= 0 {
			return 0, fileSize - 1
		}
		start = fileSize - suffix
		if start < 0 {
			start = 0
		}
		end = fileSize - 1
	} else {
		var err error
		start, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil || start < 0 {
			return 0, fileSize - 1
		}

		if parts[1] == "" {
			end = fileSize - 1
		} else {
			end, err = strconv.ParseInt(parts[1], 10, 64)
			if err != nil || end < start {
				end = fileSize - 1
			}
		}
	}

	return start, end
}
