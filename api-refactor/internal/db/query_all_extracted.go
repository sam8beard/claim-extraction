/*
Function to query object keys of all properly processed files
*/
package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllExtractedKeys(ctx context.Context, pool *pgxpool.Pool) (pgx.Rows, error) {

	return pool.Query(
		ctx,
		`SELECT s3_key FROM documents WHERE text_extracted=true ORDER BY uploaded_at;`,
	)

} // GetAllExtractedKeys

