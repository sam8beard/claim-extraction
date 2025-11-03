package processing

import (
	"bytes"
	"context"
	"tui/backend/types"
	"tui/backend/types/shared"
)

type Processing struct {
	shared.Workflow
}

/*
Execute the processing workflow
*/
func (p *Processing) Run(ctx context.Context, input *types.ProcessingInput) (*types.ProcessingResult, error) {
	var err error
	pResult := types.ProcessingResult{
		ProcessedFiles: make(map[shared.File]*bytes.Buffer),
	}
	fetchResult, err := p.Fetch(ctx, input)
	if err != nil {
		return &pResult, err
	} // if
	pResult.ProcessedFiles = fetchResult.SuccessFiles

	return &pResult, err
} // Run
