package model

// Share interface is handler for connection with share.
type Share interface {
	Echo() error                                      // Send echo request to server.
	OpenFile(filename string, mode int) (File, error) // Open shared file.
	OpenDir(path string) (Dir, error)                 // Open shared directory.
	Close()                                           // Close established connection.
}
