package models 

import "time"

type Document struct { 
	ID int
    FileName string
	UploadedAt time.Time
    Source string
    TextExtracted bool
    ContentHash string
    S3Key string
    FileSizeBytes int
}
