package processing

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"tui/backend/types/shared"
)

type TestFile struct {
	ObjectKey string
	FileName  string
	Content   *bytes.Buffer
}

func getTestFetchResult() *FetchResult {
	return &FetchResult{
		SuccessFiles: map[shared.File]*bytes.Buffer{
			{ObjectKey: "processed/source1/file1.txt"}: bytes.NewBufferString(
				"Dr. Smith claims that the vaccine is highly effective. " +
					"However, recent studies suggest otherwise. " +
					"The WHO confirms that further testing is required. " +
					"Experts debate the methodology used in these studies. " +
					"Ultimately, the evidence remains inconclusive.",
			),
			{ObjectKey: "processed/source2/file2.txt"}: bytes.NewBufferString(
				"Alice asserts that the new policy will reduce emissions. " +
					"Bob counters that the economic impact will be severe. " +
					"The government releases official figures supporting Alice's claim. " +
					"Environmental groups respond positively, emphasizing long-term benefits. " +
					"Analysts remain skeptical about the short-term effects.",
			),
			{ObjectKey: "processed/source3/file3.txt"}: bytes.NewBufferString(
				"The article reports that inflation has risen by 3.2% over the past quarter. " +
					"Economists warn that this trend may continue if interest rates are not adjusted. " +
					"Consumer groups note the rising cost of living as evidence. " +
					"Some banks argue that this increase is within expected limits.",
			),
			{ObjectKey: "processed/source4/file4.txt"}: bytes.NewBufferString(
				"NASA confirms that the satellite has successfully entered orbit. " +
					"Scientists observe minor deviations in trajectory, which are under investigation. " +
					"Independent analysts suggest that the mission's success will advance space research. " +
					"The agency releases images and data for public review.",
			),
			{ObjectKey: "processed/source5/file5.txt"}: bytes.NewBufferString(
				"According to several reports, the technology startup achieved record revenue last year. " +
					"Investors praise the management team for strategic decisions. " +
					"Competitors question the sustainability of such growth. " +
					"Financial analysts provide a cautious outlook for the upcoming fiscal year.",
			),
		},
	}
} // getFetchResult

func printNLPResult(nr *NLPResult) {
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

func TestNLP(t *testing.T) {
	p := Processing{}
	var err error
	ctx := context.Background()
	fetchResult := getTestFetchResult()

	nlpResult, err := p.NLP(ctx, fetchResult)
	if err != nil {
		panic(err)
	} // if
	printNLPResult(nlpResult)

} // TestNLP
