package model

// DirEntryType is directory entry type.
type DirEntryType string

// Available dir entry types.
const (
	DirEntryTypeLink    DirEntryType = "link"
	DirEntryTypeFile    DirEntryType = "file"
	DirEntryTypeDir     DirEntryType = "dir"
	DirEntryTypeUnknown DirEntryType = "unknown"
)

// DirEntry is entry in directory list.
type DirEntry struct {
	Type DirEntryType
	Name string
	Size uint64 // 0 for dir.
}

// Dir is interface for smb2 directory.
type Dir interface {
	List() ([]*DirEntry, error) // List returns directory files.
	Close() error               // Closes current directory.
}
