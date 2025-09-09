package models

type FileSize int64

const (
	B  FileSize = 1
	KB FileSize = 1024
	MB FileSize = KB * 1024
	GB FileSize = MB * 1024
)
