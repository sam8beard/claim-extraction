package conversion

import (
	"context"
	"tui/backend/types/shared"
)

type UpdateResult struct {
}

/*
Updates text_extracted for rows with extracted text
*/
func (c *Conversion) Update(ctx context.Context, u *shared.UploadResult) (*UpdateResult, error) {
	var err error
	updateResult := UpdateResult{}

	for _, file := range u.SuccessFiles {
		_, err := c.PGClient.Exec(
			ctx,
			`UPDATE documents SET text_extracted = true WHERE s3_key = $1`,
			file.ObjectKey,
		)
		if err != nil {

			// handle error
			continue
		} // if

	} // for
	return &updateResult, err
} // Update
