package conversion

import "context"

type UpdateResult struct {
}

/*
Updates text_extracted for rows with extracted text
*/
func (c *Conversion) Update(ctx context.Context, e *ExtractionResult) (*UpdateResult, error) {
	var err error
	updateResult := UpdateResult{}

	return &updateResult, err
} // UpdateRows
