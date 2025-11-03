package conversion

import (
	"context"
	"errors"

	"github.com/sam8beard/claim-extraction/api-refactor/internal/types/shared"
)

/*
Updates text_extracted for rows with extracted text
*/
func (c *Conversion) Update(ctx context.Context, f shared.FileID) error {
	var err error
	_, err = c.PGClient.Exec(
		ctx,
		`UPDATE documents SET text_extracted = true WHERE s3_key = $1`,
		f.OriginalKey,
	)
	if err != nil {
		return errors.New("unable to update row in documents")
	} // if

	return err
} // Update
