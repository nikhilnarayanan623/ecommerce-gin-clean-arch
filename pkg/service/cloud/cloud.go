package cloud

import (
	"context"
	"mime/multipart"
)

type CloudService interface {
	SaveFile(ctx context.Context, fileHeader *multipart.FileHeader) (uploadId string, err error)
	GetFileUrl(ctx context.Context, uploadID string) (url string, err error)
}
