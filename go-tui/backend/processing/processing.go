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
func (p *Processing) Run(ctx context.Context, input *types.ProcessingInput) (*NLPResult, error) {
	var err error
	fetchResult, err := p.Fetch(ctx, input)
	if err != nil {
		return nil, err
	} // if

	nlpResult, err := p.NLP(ctx, fetchResult)
	if err != nil {
		return nil, err
	} // if

	return nlpResult, err
} // Run
