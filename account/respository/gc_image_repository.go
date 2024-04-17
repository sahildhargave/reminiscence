package repository

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"cloud.google.com/go/storage"
	"github.com/sahildhargave/memories/account/model"
	"github.com/sahildhargave/memories/account/model/apperrors"
)

type gcImageRepository struct {
	Storage    *storage.Client
	BucketName string
}

// NewImageRepository creates a new instance of gcImageRepository.
func NewImageRepository(gcClient *storage.Client, bucketName string) model.ImageRepository {
	return &gcImageRepository{
		Storage:    gcClient,
		BucketName: bucketName,
	}
}

func (r *gcImageRepository) DeleteProfile(ctx context.Context, objName string) error {
	bckt := r.Storage.Bucket(r.BucketName)

	object := bckt.Object(objName)

	if err := object.Delete(ctx); err != nil {
		log.Printf("Failed to delete image object with ID: %s from GOOGLE CLOUD Storage\n", objName)
		return apperrors.NewInternal()
	}
	return nil
}

// UpdateProfile updates a profile image in Google Cloud Storage.
func (r *gcImageRepository) UpdateProfile(ctx context.Context, objName string, imageFile multipart.File) (string, error) {
	bckt := r.Storage.Bucket(r.BucketName)
	object := bckt.Object(objName)
	wc := object.NewWriter(ctx)

	wc.ObjectAttrs.CacheControl = "no-cache, max-age=0"

	_, err := io.Copy(wc, imageFile)
	if err != nil {
		log.Printf("Unable to write file to Google Cloud Storage: %v\n", err)
		return "", apperrors.NewInternal()
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", r.BucketName, objName)
	return imageURL, nil
}
