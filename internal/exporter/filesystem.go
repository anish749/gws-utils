package exporter

import "os"

// Exporter writes content to a destination.
type Exporter interface {
	WriteFile(path string, content string) error
	EnsureDir(path string) error
}

// FilesystemExporter writes content to the local filesystem.
type FilesystemExporter struct{}

func NewFilesystemExporter() *FilesystemExporter {
	return &FilesystemExporter{}
}

func (e *FilesystemExporter) WriteFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func (e *FilesystemExporter) EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
