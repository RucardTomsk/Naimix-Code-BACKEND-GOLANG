package s3

import (
	"context"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/enum"
	"github.com/google/uuid"
	"io"
)

// ObjectStoreService is used to communicate with s3 object storage.
type ObjectStoreService interface {
	Upload(ctx context.Context, bucket string, input UploadInput) (*uuid.UUID, *base.ServiceError)
	GetImageFileURL(ctx context.Context, bucket string, fileName string) (*string, *base.ServiceError)
	GetDocumentFileURL(ctx context.Context, bucket string, fileName string, fileType string) (*string, *base.ServiceError)
	RemoveImage(ctx context.Context, bucket string, fileID uuid.UUID) *base.ServiceError
	RemoveDocument(ctx context.Context, fileID uuid.UUID, fileType string) *base.ServiceError
	UploadAsWebP(ctx context.Context, bucket enum.Bucket, file io.Reader) (*uuid.UUID, *base.ServiceError)
	GetWebPFileURL(ctx context.Context, bucket enum.Bucket, fileID uuid.UUID) (string, *base.ServiceError)
	DeleteWebPFile(ctx context.Context, bucket enum.Bucket, fileID uuid.UUID) *base.ServiceError
}
