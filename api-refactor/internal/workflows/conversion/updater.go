package conversion

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/sam8beard/claim-extraction/api-refactor/internal/types/shared"
)

/*
Updates text_extracted for rows with extracted text
*/
func (c *Conversion) Update(ctx context.Context, f shared.FileID) error {
	var err error
	log.Fatalf("Original key in update: %s\n", f.OriginalKey)
	res, err := c.PGClient.Exec(
		ctx,
		`UPDATE documents SET text_extracted = true WHERE s3_key = $1`,
		f.OriginalKey,
	)
	if err != nil {
		return errors.New("unable to update row in documents")
	} // if

	count := res.RowsAffected()
	log.Printf("Rows updated: %d for key: %s", count, f.OriginalKey)
	if count == 0 {
		return fmt.Errorf("no rows updated for key: %s", f.OriginalKey)
	}
	return nil
} // Update
