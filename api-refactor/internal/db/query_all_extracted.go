/*
Function to query object keys of all properly processed files
*/
package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/types"
)

func GetAllExtractedKeys(ctx context.Context, pool *pgxpool.Pool, input *types.ProcessingInput) (pgx.Rows, error) {
	keysToFetch := make([]string, 0)
	files := input.ConvertedFiles
	// grab all keys of converted files
	for _, f := range files {
		// we want the ORIGINAL KEY of the file, not the processed one
		keysToFetch = append(keysToFetch, f.OriginalKey)
	} // for
	return pool.Query(
		ctx,
		`SELECT s3_key FROM documents WHERE text_extracted=true AND s3_key = ANY($1)`,
		keysToFetch,
	)

} // GetAllExtractedKeys
