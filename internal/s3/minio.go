package s3

import (
	"bytes"
	"context"
	"fmt"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/enum"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/nickalie/go-webpbin"
	"image"
	_ "image/gif"  // Добавляем для поддержки формата GIF
	_ "image/jpeg" // Добавляем для поддержки формата JPEG
	_ "image/png"  // Добавляем для поддержки формата PNG
	"io"
	"net/http"
	"net/url"
	"strings"

	"time"
)

type UploadInput struct {
	File        io.Reader
	Size        int64
	ContentType string
}

// MinioService communicates with minio (s3).
type MinioService struct {
	client *minio.Client
}

func NewMinioService(client *minio.Client) *MinioService {
	return &MinioService{
		client: client,
	}
}

func (s *MinioService) Upload(ctx context.Context, bucket string, input UploadInput) (*uuid.UUID, *base.ServiceError) {
	opts := minio.PutObjectOptions{
		ContentType:  input.ContentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	key := uuid.New()
	keyStr := key.String() + "." + strings.Split(input.ContentType, "/")[1]
	_, err := s.client.PutObject(ctx, bucket, keyStr, input.File, input.Size, opts)
	if err != nil {
		return nil, unexpectedServiceError(err)
	}

	return &key, nil
}

func (s *MinioService) UploadAsWebP(ctx context.Context, bucket enum.Bucket, file io.Reader) (*uuid.UUID, *base.ServiceError) {
	// Декодируем изображение из входного файла
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, unexpectedServiceError(err)
	}

	// Создаем буфер для сохранения изображения в формате WebP
	var webpBuffer bytes.Buffer

	err = webpbin.NewCWebP(
		webpbin.SetSkipDownload(true),
		webpbin.SetVendorPath("")).
		Quality(80).
		InputImage(img).
		Output(&webpBuffer).
		Run()

	if err != nil {
		return nil, unexpectedServiceError(err)
	}

	// Определяем новый контент-тайп для WebP
	contentType := "image/webp"

	// Генерируем новый ключ для WebP файла
	key := uuid.New()
	keyStr := key.String() + ".webp"

	// Определяем параметры загрузки для WebP
	opts := minio.PutObjectOptions{
		ContentType:  contentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	// Загружаем WebP изображение в MinIO
	_, err = s.client.PutObject(ctx, string(bucket), keyStr, &webpBuffer, int64(webpBuffer.Len()), opts)
	if err != nil {
		return nil, unexpectedServiceError(err)
	}

	return &key, nil
}

func (s *MinioService) GetWebPFileURL(ctx context.Context, bucket enum.Bucket, fileID uuid.UUID) (string, *base.ServiceError) {
	// Генерируем ключ на основе fileID (предполагаем, что файл сохраняется с расширением .webp)
	keyStr := fileID.String() + ".webp"

	// Определяем срок действия URL (например, 1 час)
	expiry := time.Hour

	// Генерируем подписанную ссылку на объект
	reqParams := make(url.Values)
	reqParams.Set("response-content-type", "image/webp")

	resignedURL, err := s.client.PresignedGetObject(ctx, string(bucket), keyStr, expiry, reqParams)
	if err != nil {
		return "", unexpectedServiceError(err)
	}

	return resignedURL.String(), nil
}

func (s *MinioService) DeleteWebPFile(ctx context.Context, bucket enum.Bucket, fileID uuid.UUID) *base.ServiceError {
	// Генерируем ключ на основе fileID (предполагаем, что файл сохранен с расширением .webp)
	keyStr := fileID.String() + ".webp"

	// Удаляем файл из MinIO
	err := s.client.RemoveObject(ctx, string(bucket), keyStr, minio.RemoveObjectOptions{})
	if err != nil {
		return unexpectedServiceError(err)
	}

	return nil
}

func (s *MinioService) GetImageFileURL(ctx context.Context, bucket string, fileKey string) (*string, *base.ServiceError) {

	_, err := s.client.StatObject(context.Background(), bucket, fileKey+".png", minio.StatObjectOptions{})
	if err != nil {
		_, err := s.client.StatObject(context.Background(), bucket, fileKey+".jpeg", minio.StatObjectOptions{})
		if err != nil {
			return nil, base.NewNotFoundError(err)
		} else {
			fileKey += ".jpeg"
		}
	} else {
		fileKey += ".png"
	}

	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\""+fileKey+"\"")

	resignedURL, err := s.client.PresignedGetObject(ctx, bucket, fileKey, time.Duration(1000)*time.Second, reqParams)
	if err != nil {
		return nil, unexpectedServiceError(err)
	}

	urlString := resignedURL.String()
	return &urlString, nil
}

func (s *MinioService) GetDocumentFileURL(ctx context.Context, bucket string, fileKey string, fileType string) (*string, *base.ServiceError) {

	_, err := s.client.StatObject(context.Background(), bucket, fmt.Sprintf("%s.%s", fileKey, fileType), minio.StatObjectOptions{})
	if err != nil {
		if err != nil {
			return nil, base.NewNotFoundError(err)
		}
	}

	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\""+fmt.Sprintf("%s.%s", fileKey, fileType)+"\"")

	resignedURL, err := s.client.PresignedGetObject(ctx, bucket, fmt.Sprintf("%s.%s", fileKey, fileType), time.Duration(1000)*time.Second, reqParams)
	if err != nil {
		return nil, unexpectedServiceError(err)
	}

	urlString := resignedURL.String()
	return &urlString, nil
}

func (s *MinioService) RemoveDocument(ctx context.Context, fileID uuid.UUID, fileType string) *base.ServiceError {
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}

	if err := s.client.RemoveObject(ctx, "document", fmt.Sprintf("%s.%s", fileID.String(), fileType), opts); err != nil {
		return unexpectedServiceError(err)
	}

	return nil
}

func (s *MinioService) RemoveImage(ctx context.Context, bucket string, fileID uuid.UUID) *base.ServiceError {
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}

	if err := s.client.RemoveObject(ctx, bucket, fileID.String()+".png", opts); err != nil {
		return unexpectedServiceError(err)
	}

	if err := s.client.RemoveObject(ctx, bucket, fileID.String()+".jpeg", opts); err != nil {
		return unexpectedServiceError(err)
	}

	return nil
}

// unexpectedServiceError returns any unclassified service error.
func unexpectedServiceError(err error) *base.ServiceError {
	return &base.ServiceError{
		Err:     err,
		Blame:   base.BlameServer,
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
}
