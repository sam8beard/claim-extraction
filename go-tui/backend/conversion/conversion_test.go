package conversion

import (
	"context"
	"log"
	"math/rand"
	"testing"

	"tui/backend/acquisition"
	"tui/backend/types"
)

var ctx = context.Background()

func NewConversion() *Conversion {
	c := Conversion{}
	err := c.InitializeClients()
	if err != nil {
		log.Fatalf("failed to initialize clients")
	}
	return &c
} // NewAcquisition

// randomAcquisitionInput generates a random AcquisitionInput with a random query and file count <= 4.
func randomAcquisitionInput() types.AcquisitionInput {
	queries := []string{
		"generative ai for music composition",
		"computer vision in medical imaging",
		"nlp for legal document summarization",
		"ai-powered fraud detection in banking",
		"predictive analytics for energy consumption",
		"speech recognition for smart home devices",
		"machine learning in sports injury prediction",
		"ai for personalized learning platforms",
		"deep learning for crop disease detection",
		"autonomous robots in warehouse logistics",
		"ai-driven customer sentiment analysis",
		"virtual reality for remote collaboration",
		"blockchain for digital identity management",
		"ai in wildlife conservation monitoring",
		"data mining for healthcare diagnostics",
		"reinforcement learning for game development",
		"ai for predictive maintenance in aviation",
		"image segmentation in satellite imagery",
		"ai-powered translation for global business",
		"machine learning for traffic flow optimization",
		"ai in personalized nutrition planning",
		"computer vision for retail inventory tracking",
		"nlp for automated contract review",
		"ai for disaster response coordination",
		"deep learning for handwriting recognition",
		"ai in smart grid management",
		"machine learning for personalized advertising",
		"ai-powered recommendation in streaming services",
		"computer vision for food quality inspection",
		"ai for remote patient monitoring",
	}
	fileCount := 1 + rand.Intn(4) // random number between 1 and 4

	return types.AcquisitionInput{
		Query:     queries[rand.Intn(len(queries))],
		FileCount: fileCount,
	}
}

func NewConversionInput() types.ConversionInput {
	a := acquisition.Acquisition{}
	a.InitializeClients()
	// input := types.AcquisitionInput{
	// 	Query:     "artificial intelligence in healthcare",
	// 	FileCount: 5,
	// }

	// input := types.AcquisitionInput{
	// 	Query:     "blockchain applications in finance",
	// 	FileCount: 8,
	// }

	// input := types.AcquisitionInput{
	// 	Query:     "privacy regulations in europe",
	// 	FileCount: 3,
	// }

	input := randomAcquisitionInput()

	// input := types.AcquisitionInput{
	// 	Query:     "natural language processing trends",
	// 	FileCount: 6,
	// }

	// input := types.AcquisitionInput{
	// 	Query:     "cybersecurity best practices",
	// 	FileCount: 9,
	// }

	// input := types.AcquisitionInput{
	// 	Query:     "cloud computing adoption",
	// 	FileCount: 11,
	// }

	// input := types.AcquisitionInput{
	// 	Query:     "data science in retail",
	// 	FileCount: 4,
	// }

	// input := types.AcquisitionInput{
	// 	Query:     "robotics in manufacturing",
	// 	FileCount: 13,
	// }

	// input := types.AcquisitionInput{
	// 	Query:     "autonomous vehicles safety",
	// 	FileCount: 3,
	// }

	// input := types.AcquisitionInput{
	// 	Query:     "quantum computing breakthroughs",
	// 	FileCount: 14,
	// }

	// input := types.AcquisitionInput{
	// 	Query:     "edge computing use cases",
	// 	FileCount: 2,
	// }

	// input := types.AcquisitionInput{
	// 	Query:     "internet of things security",
	// 	FileCount: 15,
	// }

	// input := types.AcquisitionInput{
	// 	Query:     "big data analytics",
	// 	FileCount: 10,
	// }

	// input := types.AcquisitionInput{
	// 	Query:     "virtual reality in education",
	// 	FileCount: 6,
	// }

	// input := types.AcquisitionInput{
	// 	Query:     "ethical issues in ai",
	// 	FileCount: 8,
	// }

	acResult, err := a.Run(ctx, input)
	if err != nil {
		log.Fatal(err)
	} // if

	conInput := types.ConversionInput(acResult)
	return conInput
}

func TestDownload(t *testing.T) {
	input := NewConversionInput()
	con := NewConversion()

	_, err := con.Download(ctx, input)
	if err != nil {

		panic(err)
	}
}

func TestExtract(t *testing.T) {
	input := NewConversionInput()
	con := NewConversion()

	downloadResult, err := con.Download(ctx, input)

	if err != nil {
		panic(err)
	}
	eResult, err := con.Extract(ctx, downloadResult)
	if err != nil {
		t.Fatalf("Extraction failed")
	}
	t.Logf("%v", eResult)
}

// type testData struct {
// 	data     map[shared.FileID]string
// 	byteData map[shared.FileID][]byte
// }
// type testLocker struct {
// 	mu sync.Mutex
// 	td testData
// }

// func (l *testLocker) logEntry(id shared.FileID, u any) {
// 	l.mu.Lock()
// 	defer l.mu.Unlock()

// 	switch val := u.(type) {
// 	case string:
// 		l.td.data[id] = val
// 	case []byte:
// 		l.td.byteData[id] = val
// 	} // switch
// } // logEntry
