package processing

import (
	"bytes"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/minio/minio-go/v7"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/db"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/types"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/types/shared"
	"io"
	"log"
	"path"
	"strings"
)

type FetchResult struct {
	SuccessFiles map[shared.File]*bytes.Buffer
}

/*
Fetch text bodies of converted files
*/
func (p *Processing) Fetch(ctx context.Context, input *types.ProcessingInput) (*FetchResult, error) {
	fetchResult := FetchResult{
		SuccessFiles: make(map[shared.File]*bytes.Buffer, 0),
	}

	// get rows of files that have been converted
	extractedRows, err := db.GetAllExtractedKeys(ctx, p.PGClient, input)
	if err != nil {
		log.Println("error in GetAllExtractedKeys block")
		err := errors.New("unable to query rows of extracted text")
		return nil, err
	} // if

	keys, err := FetchKeys(ctx, extractedRows)
	log.Printf("all keys retrieved by FetchKeys: %v", keys)

	if err != nil {
		log.Println("error in fetch keys block")
		return nil, err
	} // if
	bucket := p.MinioClient.Bucket
	opts := minio.GetObjectOptions{}
	// download extracted files from processed
	for _, key := range keys {
		log.Printf("attempting to get %s\n", key)
		object, err := p.MinioClient.Client.GetObject(
			ctx,
			bucket,
			key,
			opts,
		)
		/*
			Apparently, err is not populated if the retrieval process executes, but the key is determined to not exist in bucket. So this err check only covers the case that the call to GetObject was not successful.
		*/
		// could not get object from MinIO
		if err != nil {
			continue
		} // if

		// TESTING
		info, _ := object.Stat()
		obK := info.Key
		// The specified key does not exist
		if obK == "" {
			log.Printf("warning: no key found in bucket for %s\n", key)
			log.Printf("check bucket to see if key exists\n")
			continue
		} // if
		log.Printf("key found in bucket: %s\n", obK)

		// buf for object
		var newBuff bytes.Buffer
		_, err = io.Copy(&newBuff, object)
		if err != nil {
			log.Printf("unable to copy contents: %s\n", key)
			continue
		} // if
		newSFile := shared.File{
			ObjectKey: key,
		}
		fetchResult.SuccessFiles[newSFile] = &newBuff
	} // for
	return &fetchResult, nil
} // Fetch

func FetchKeys(ctx context.Context, rows pgx.Rows) ([]string, error) {
	var err error

	keys := make([]string, 0)
	//log.Print("firing right before row scanning loop")
	// get keys
	for rows.Next() {
		var key string
		if err = rows.Scan(&key); err != nil {
			log.Println("firing in err block on scan")
			continue
		} // if

		ext := path.Ext(key)
		newKey := strings.Replace(key, ext, ".txt", 1)
		newKey = strings.Replace(newKey, "raw", "processed", 1)
		keys = append(keys, newKey)
		//log.Printf("%v", keys)
	} // for

	if len(keys) == 0 {
		err = errors.New("could not fetch any keys")
		return keys, err
	} // if

	return keys, err
} // FetchKeys
