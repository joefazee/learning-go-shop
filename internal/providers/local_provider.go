package providers

import (
	"mime/multipart"
	"os"
	"path/filepath"
)

type LocalUploadProvider struct {
	basePath string
}

func NewLocalUploadProvider(basePath string) *LocalUploadProvider {
	return &LocalUploadProvider{basePath: basePath}
}

func (p *LocalUploadProvider) UploadFile(file *multipart.FileHeader, path string) (string, error) {

	fullPath := filepath.Join(p.basePath, path)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return "", err
	}

	// Open source
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// create destination
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// read from source to destination
	if _, err := dst.ReadFrom(src); err != nil {
		return "", err
	}

	return path, nil

}

func (p *LocalUploadProvider) DeleteFile(path string) error {
	fullPath := filepath.Join(p.basePath, path)
	return os.Remove(fullPath)
}
