/*
Functions for inserting metadata into Postgres tables
*/
package db

import (
	"context"
	"tui/backend/db/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertDocumentMetadata(ctx context.Context, pool *pgxpool.Pool, doc *models.Document) error {
	return pool.QueryRow(
		ctx,
		`INSERT INTO documents 
		(file_name, source, text_extracted, content_hash, s3_key, file_size_bytes)
		VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING
		RETURNING id, uploaded_at;`, doc.FileName, doc.Source, doc.TextExtracted, doc.ContentHash, doc.S3Key,
		doc.FileSizeBytes,
	).Scan(&doc.ID, &doc.UploadedAt)

} // InsertDocumentMetadata
