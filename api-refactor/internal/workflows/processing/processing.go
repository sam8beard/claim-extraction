package processing

import (
	"context"
	"fmt"

	"github.com/sam8beard/claim-extraction/api-refactor/internal/types"
	"github.com/sam8beard/claim-extraction/api-refactor/internal/types/shared"
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

func PrintNLPResult(nr *NLPResult) {
	fileData := nr.FileData
	for _, data := range fileData {
		fmt.Println("\nFile data ---------------------------")
		fmt.Printf("\nObject Key: %s\n", data.ObjectKey)
		fmt.Printf("\nFile Name: %s\n", data.FileName)
		fmt.Printf("\nClaim Score: %f\n", data.ClaimScore)
		fmt.Printf("\nSpan data for %s ---------------------------\n", data.FileName)
		for _, spanData := range data.ClaimSpans {
			fmt.Printf("\nType: %s\n", spanData.Type)
			fmt.Printf("\nText: %s\n", spanData.Text)
			fmt.Printf("\nSent: %s\n", spanData.Sent)
			fmt.Printf("\nConfidence: %f\n", spanData.Confidence)
		} // for
	} // for
} // printNLPResult
