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
		// New topics below
		"ai for climate change modeling",
		"robotics in surgical procedures",
		"machine learning for stock market prediction",
		"ai in autonomous drone navigation",
		"computer vision for wildlife tracking",
		"nlp for mental health assessment",
		"ai-powered supply chain optimization",
		"deep learning for protein structure prediction",
		"ai in smart city infrastructure",
		"machine learning for personalized medicine",
		"ai for digital art generation",
		"computer vision for gesture recognition",
		"nlp for customer support automation",
		"ai in renewable energy forecasting",
		"machine learning for fraud prevention in e-commerce",
		"ai-powered document classification",
		"deep learning for speech synthesis",
		"ai in autonomous vehicle navigation",
		"computer vision for facial recognition security",
		"nlp for sentiment analysis in social media",
		"ai for predictive policing",
		"machine learning for disease outbreak prediction",
		"ai-powered chatbots for healthcare",
		"deep learning for video content analysis",
		"ai in agricultural yield prediction",
		"computer vision for traffic sign detection",
		"nlp for financial news analysis",
		"ai for personalized workout recommendations",
		"machine learning for anomaly detection in manufacturing",
		"ai-powered plagiarism detection",
		"deep learning for medical image classification",
		"ai in smart home automation",
		"computer vision for emotion detection",
		"nlp for resume screening",
		"ai for predictive asset management",
		"machine learning for energy grid optimization",
		"ai-powered virtual assistants",
		"deep learning for music genre classification",
		"ai in retail demand forecasting",
		"computer vision for defect detection in products",
		"nlp for legal case outcome prediction",
		"ai for personalized travel recommendations",
		"machine learning for credit scoring",
		"ai-powered fraud detection in insurance",
		"deep learning for object detection in videos",
		"ai in telemedicine diagnostics",
		"computer vision for sports analytics",
		"nlp for automated essay grading",
		"ai for predictive equipment maintenance",
		"machine learning for personalized news feeds",
		"ai-powered translation for scientific literature",
		"deep learning for handwriting generation",
		"ai in logistics route optimization",
		"computer vision for plant species identification",
		"nlp for chatbot conversation analysis",
		"ai for predictive weather forecasting",
		"machine learning for personalized shopping experiences",
		"ai-powered cybersecurity threat detection",
		"deep learning for anomaly detection in time series data",
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

<<<<<<< HEAD
type Lockbox struct {
	mu sync.Mutex
	r  Result
}
type Result struct {
	data map[shared.FileID][]byte
	err  map[shared.FileID]string
}

func (l *Lockbox) modifyLocker(f shared.FileID, u any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	switch val := u.(type) {
	case string:
		l.r.err[f] = val
	case []byte:
		l.r.data[f] = val
	}
}

func TestExtractTwo(t *testing.T) {
	msg := []byte("hello")
	fmsg := "hello"

	file := shared.FileID{
		Title: string("cool"),
	}
	fFile := shared.FileID{
		Title: string("not cool"),
	}

	l := Lockbox{
		r: Result{
			data: make(map[shared.FileID][]byte),
			err:  make(map[shared.FileID]string),
		},
	}

	var wg sync.WaitGroup
	wg.Add(2)

	doSomethingGood := func() {
		defer wg.Done()
		l.modifyLocker(file, msg)
	}
	go doSomethingGood()
	doSomethingBad := func() {
		defer wg.Done()
		l.modifyLocker(fFile, fmsg)
	}
	go doSomethingBad()

	wg.Wait()

	t.Log("firing")
	t.Log(l.r.data)
	t.Log(l.r.err)
}
=======
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
>>>>>>> testing/piping-issue
