package model

import "os"

// File is interface for smb file handler.
type File interface {
	Read(p []byte) (n int, err error)             // Read from file no buffer (impl io.Reader).
	Seek(offset int64, whence int) (int64, error) // Seek position into file (impl io.Seeker).
	Stat() (os.FileInfo, error)                   // Get file statistic.
	Close() error                                 // Close current file.
}
