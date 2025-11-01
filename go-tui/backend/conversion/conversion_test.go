package conversion

import (
	"context"
	"fmt"
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
func randomAcqInput() types.AcquisitionInput {
	queries := []string{
		// New topics below
		"ai for wildfire risk prediction",
		"machine learning for autonomous warehouse robots",
		"ai in real-time language translation",
		"computer vision for underwater exploration",
		"nlp for medical record summarization",
		"ai-powered energy consumption optimization",
		"deep learning for protein-ligand interaction prediction",
		"ai in traffic flow optimization",
		"machine learning for satellite image classification",
		"ai for early disease diagnosis",
		"computer vision for drone-based inspection",
		"nlp for automated meeting transcription",
		"ai-powered recommendation systems for education",
		"deep learning for 3d object reconstruction",
		"ai in predictive maintenance for aircraft",
		"machine learning for customer churn prediction",
		"ai for environmental monitoring with IoT sensors",
		"computer vision for automated quality control",
		"nlp for legal contract analysis",
		"ai-powered robotic process automation",
		"deep learning for natural disaster detection",
		"ai in wildlife population monitoring",
		"machine learning for financial fraud detection",
		"ai for smart energy management in buildings",
		"computer vision for traffic accident detection",
		"nlp for opinion mining in product reviews",
		"ai-powered personalized learning platforms",
		"deep learning for video game character animation",
		"ai in predictive inventory management",
		"machine learning for disease risk assessment",
		"ai for autonomous farming vehicles",
		"computer vision for historical document digitization",
		"nlp for scientific paper summarization",
		"ai-powered smart home security",
		"deep learning for autonomous robot navigation",
		"ai in urban air quality prediction",
		"machine learning for loan default prediction",
		"ai for personalized diet recommendations",
		"computer vision for warehouse item tracking",
		"nlp for multilingual chat translation",
		"ai-powered digital marketing optimization",
		"deep learning for 3d medical imaging",
		"ai in automated traffic signal control",
		"machine learning for social media trend prediction",
		"ai for smart wearable devices",
		"computer vision for drone mapping",
		"nlp for intelligent tutoring systems",
		"ai-powered document summarization",
		"deep learning for weather pattern analysis",
		"ai in autonomous shipping logistics",
		"machine learning for predictive maintenance in factories",
		"ai for personalized content curation",
		"computer vision for surgical tool tracking",
		"nlp for mental health chatbots",
		"ai-powered supply chain risk assessment",
		"deep learning for facial expression recognition",
		"ai in renewable energy load forecasting",
		"machine learning for personalized investment advice",
		"ai for automated vehicle inspection",
		"computer vision for packaging defect detection",
		"nlp for automated translation in real-time",
		"ai-powered smart city monitoring systems",
		"deep learning for gesture-controlled devices",
		"ai in autonomous traffic management",
		"machine learning for predictive healthcare interventions",
		"ai for precision agriculture monitoring"}
	fileCount := 1 + rand.Intn(4) // random number between 1 and 4

	return types.AcquisitionInput{
		Query:     queries[rand.Intn(len(queries))],
		FileCount: fileCount,
	}
} // randomAcqInput

func NewConversionInput() types.ConversionInput {
	a := acquisition.Acquisition{}
	err := a.InitializeClients()
	if err != nil {
		panic(err)
	}

	input := randomAcqInput()

	acResult, err := a.Run(ctx, input)
	if err != nil {
		log.Fatal(err)
	} // if

	conInput := types.ConversionInput(acResult)
	return conInput
} // NewConversionInput

func TestDownload(t *testing.T) {
	input := NewConversionInput()
	con := NewConversion()

	_, err := con.Download(ctx, input)
	if err != nil {

		panic(err)
	}
} // TestDownload

func (r ExtractionResult) String() string {
	sucF := r.SuccessFiles
	failF := r.FailedFiles
	var s string
	s += "\nSUCCESS FILES: \n"
	for fileID, body := range sucF {
		printableBody := string(body)
		truncBody := printableBody[:15] + "..."
		title := fileID.Title
		objectKey := fileID.ObjectKey
		url := fileID.URL
		s += fmt.Sprintf("\nTitle: %s, ObjectKey: %s, URL: %s, Body: %s\n", title, objectKey, url, truncBody)
	} // for

	s += "\nFAILED FILES: \n"
	for fileID, body := range failF {
		truncBody := body[:15] + "..."
		title := fileID.Title
		objectKey := fileID.ObjectKey
		url := fileID.URL
		s += fmt.Sprintf("\nTitle: %s, ObjectKey: %s, URL: %s, Body: %s\n", title, objectKey, url, truncBody)
	} // for
	return s
} // String

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
	t.Log(eResult.String())

} // TestExtract
