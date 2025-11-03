package processing

import (
	"bytes"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/minio/minio-go/v7"
	"io"
	"path"
	"strings"
	"tui/backend/db"
	"tui/backend/types"
	"tui/backend/types/shared"
)

type FetchResult struct {
	SuccessFiles map[shared.File]bytes.Buffer
}

/*
Fetch text bodies of converted files
*/
func (p *Processing) Fetch(ctx context.Context, input *types.ProcessingInput) (*FetchResult, error) {
	fetchResult := FetchResult{
		SuccessFiles: make(map[shared.File]bytes.Buffer, 0),
	}

	extractedRows, err := db.GetAllExtractedKeys(ctx, p.PGClient)
	if err != nil {
		err := errors.New("unable to query rows of extracted text")
		return &fetchResult, err
	} // if

	keys, err := FetchKeys(ctx, extractedRows)
	if err != nil {
		return &fetchResult, err
	} // if
	bucket := p.MinioClient.Bucket
	opts := minio.GetObjectOptions{}
	// download extracted files from processed/
	for _, key := range keys {
		object, err := p.MinioClient.Client.GetObject(
			ctx,
			bucket,
			key,
			opts,
		)
		if err != nil {
			continue
		} // if

		// use buffer to optimize streaming
		var newBuff bytes.Buffer
		_, err = io.Copy(&newBuff, object)
		if err != nil {
			err = errors.New("unable to copy object reader")
			return &fetchResult, err
		} // if
		newSFile := shared.File{
			ObjectKey: key,
		}
		fetchResult.SuccessFiles[newSFile] = newBuff
	} // for
	return &fetchResult, err
} // Fetch

func FetchKeys(ctx context.Context, rows pgx.Rows) ([]string, error) {
	var keys []string
	var err error

	keys = make([]string, 0)

	// get keys
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			continue
		} // if
		ext := path.Ext(key)
		newKey := strings.Replace(key, ext, ".txt", 1)
		newKey = strings.Replace(newKey, "raw", "processed", 1)
		keys = append(keys, newKey)
	} // for

	if len(keys) == 0 {
		err = errors.New("could not fetch any keys")
		return keys, err
	} // if

	return keys, err
} // FetchKeys
