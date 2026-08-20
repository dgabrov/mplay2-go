package endpoint

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/amanagement24/mplay2-go/internal/data"
	"github.com/amanagement24/mplay2-go/internal/service"
)

type UpdateMediaRequest struct {
	Adding      bool   `json:"adding"`
	Id          string `json:"id"`
	Description string `json:"description"`
}

type PostUpdateMediaEndpoint struct {
	servr         *service.Servr
	uploadsFolder string
}

func NewPostUpdateMediaEndpoint(servr *service.Servr, uploadsFolder string) *PostUpdateMediaEndpoint {
	return &PostUpdateMediaEndpoint{
		servr:         servr,
		uploadsFolder: uploadsFolder,
	}
}

func (e *PostUpdateMediaEndpoint) Handle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	payload, err := e.process(ctx, r)
	writeJsonResponse(w, payload, err)

	return nil
}

func (e *PostUpdateMediaEndpoint) process(ctx context.Context, r *http.Request) (*data.Media, error) {
	token, err := getTokenFromRequest(r)
	if err != nil {
		return nil, fmt.Errorf("user not logged in")
	}

	userID, err := e.servr.ValidateToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return nil, fmt.Errorf("failed to parse multipart form: %w", err)
	}
	defer r.MultipartForm.RemoveAll()

	id := r.FormValue("id")
	strAdd := r.FormValue("adding")
	adding := strAdd == "true"
	description := r.FormValue("description")

	req := UpdateMediaRequest{Adding: adding, Id: id, Description: description}

	file, fileHeader, fileErr := r.FormFile("file")
	hasFile := fileErr == nil
	if hasFile {
		defer file.Close()
	}

	if req.Adding {
		return e.handleAddMedia(ctx, userID, req, file, fileHeader, hasFile)
	}

	return e.handleUpdateMedia(ctx, userID, req, file, fileHeader, hasFile)
}

func (e *PostUpdateMediaEndpoint) handleAddMedia(ctx context.Context, userID string, req UpdateMediaRequest, file io.ReadCloser, fileHeader *multipart.FileHeader, hasFile bool) (*data.Media, error) {
	if !hasFile {
		return nil, fmt.Errorf("cannot add media without the actual content")
	}

	description := req.Description
	if description == "" {
		description = strings.TrimSuffix(fileHeader.Filename, filepath.Ext(fileHeader.Filename))
	}

	detector := NewMediaTypeDetector()
	contentType, fileReader, err := e.probeAndGetReader(file, detector)
	if err != nil {
		return nil, fmt.Errorf("failed to probe media: %w", err)
	}

	media := &data.Media{
		Id:          req.Id,
		UserId:      userID,
		Description: description,
		ContentType: contentType,
	}

	if err := e.servr.AddMedia(ctx, media); err != nil {
		return nil, err
	}

	fileSize, err := e.saveMediaFileStream(media.Id, fileReader)
	if err != nil {
		return nil, fmt.Errorf("failed to save media file: %w", err)
	}

	media.Size = fileSize

	if err := e.servr.UpdateMediaSize(ctx, userID, media.Id, fileSize); err != nil {
		return nil, fmt.Errorf("failed to update media size: %w", err)
	}

	if detector.IsVideo(contentType) {
		prober := NewVideoProber()
		videoFile, err := os.Open(filepath.Join(e.uploadsFolder, media.Id+".dat"))
		if err == nil {
			defer videoFile.Close()
			dims := prober.ProbeMP4Dimensions(videoFile)
			if !dims.Valid {
				videoFile.Seek(0, 0)
				dims = prober.ProbeWebMDimensions(videoFile)
			}
			if dims.Valid {
				media.Width = dims.Width
				media.Height = dims.Height
				if err := e.servr.UpdateMediaDimensions(ctx, userID, media.Id, dims.Width, dims.Height); err != nil {
					return nil, fmt.Errorf("failed to update media dimensions: %w", err)
				}
			}
		}
	}

	return media, nil
}

func (e *PostUpdateMediaEndpoint) handleUpdateMedia(ctx context.Context, userID string, req UpdateMediaRequest, file io.ReadCloser, fileHeader *multipart.FileHeader, hasFile bool) (*data.Media, error) {
	existingMedia, err := e.servr.GetMedia(ctx, userID, req.Id)
	if err != nil {
		return nil, err
	}

	if existingMedia == nil {
		return nil, fmt.Errorf("media not found")
	}

	existingMedia.Description = req.Description

	if !hasFile {
		if err := e.servr.UpdateMedia(ctx, userID, req.Id, req.Description); err != nil {
			return nil, err
		}
		return existingMedia, nil
	}

	detector := NewMediaTypeDetector()
	contentType, fileReader, err := e.probeAndGetReader(file, detector)
	if err != nil {
		return nil, fmt.Errorf("failed to probe media: %w", err)
	}

	existingMedia.ContentType = contentType
	existingMedia.Description = req.Description

	if err := e.servr.UpdateMediaWithType(ctx, userID, req.Id, req.Description, contentType); err != nil {
		return nil, err
	}

	fileSize, err := e.saveMediaFileStream(req.Id, fileReader)
	if err != nil {
		return nil, fmt.Errorf("failed to save media file: %w", err)
	}

	existingMedia.Size = fileSize

	if err := e.servr.UpdateMediaSize(ctx, userID, req.Id, fileSize); err != nil {
		return nil, fmt.Errorf("failed to update media size: %w", err)
	}

	if detector.IsVideo(contentType) {
		prober := NewVideoProber()
		videoFile, err := os.Open(filepath.Join(e.uploadsFolder, req.Id+".dat"))
		if err == nil {
			defer videoFile.Close()
			dims := prober.ProbeMP4Dimensions(videoFile)
			if !dims.Valid {
				videoFile.Seek(0, 0)
				dims = prober.ProbeWebMDimensions(videoFile)
			}
			if dims.Valid {
				existingMedia.Width = dims.Width
				existingMedia.Height = dims.Height
				if err := e.servr.UpdateMediaDimensions(ctx, userID, req.Id, dims.Width, dims.Height); err != nil {
					return nil, fmt.Errorf("failed to update media dimensions: %w", err)
				}
			}
		}
	} else {
		existingMedia.Width = 0
		existingMedia.Height = 0
		if err := e.servr.UpdateMediaDimensions(ctx, userID, req.Id, 0, 0); err != nil {
			return nil, fmt.Errorf("failed to reset media dimensions: %w", err)
		}
	}

	return existingMedia, nil
}

func (e *PostUpdateMediaEndpoint) probeAndGetReader(src io.Reader, detector *MediaTypeDetector) (string, io.Reader, error) {
	signature := make([]byte, 512)
	n, err := src.Read(signature)
	if err != nil && err != io.EOF {
		return "", nil, fmt.Errorf("failed to read file signature: %w", err)
	}

	signatureData := signature[:n]
	contentType, _ := detector.detectFromSignaturePublic(signatureData)

	combinedReader := io.MultiReader(bytes.NewReader(signatureData), src)

	return contentType, combinedReader, nil
}

func (e *PostUpdateMediaEndpoint) saveMediaFileStream(id string, src io.Reader) (int64, error) {
	filePath := filepath.Join(e.uploadsFolder, id+".dat")
	dst, err := os.Create(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	written, err := io.Copy(dst, src)
	if err != nil {
		_ = os.Remove(filePath)
		return 0, fmt.Errorf("failed to write file: %w", err)
	}

	return written, nil
}
