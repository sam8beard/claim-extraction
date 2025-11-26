package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sam8beard/claim-extraction/internal/types"
	"github.com/sam8beard/claim-extraction/internal/workflows/acquisition"
	"github.com/sam8beard/claim-extraction/internal/workflows/conversion"
	"github.com/sam8beard/claim-extraction/internal/workflows/processing"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	// configure logger
	log.SetFlags(log.Lshortfile)
	// initialize workflows
	a := acquisition.Acquisition{}
	c := conversion.Conversion{}
	p := processing.Processing{}

	if err := a.InitializeClients(); err != nil {
		panic(err)
	} // if
	if err := c.InitializeClients(); err != nil {
		panic(err)
	} // if
	if err := p.InitializeClients(); err != nil {
		panic(err)
	} // if

	ctx := context.Background()

	// cli args for testing
	q := string(os.Args[1])
	fc, _ := strconv.Atoi(os.Args[2])

	acqInput := types.AcquisitionInput{
		Query:     q,
		FileCount: fc,
	}

	resultsFileName := strings.ReplaceAll(q, " ", "-")

	// open logging file for logs
	logDir := "logs"
	logPath := fmt.Sprintf("%s/%s.txt", logDir, resultsFileName)

	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("failed to create logs directory: %v", err)
	} // if

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	fmt.Printf("\nAttempting to download %d files for topic: %s...\n\n", fc, q)

	log.Println("\n\t\t------------ Starting acquisition... ------------ \n\n")
	acqResult, err := a.Run(ctx, acqInput)
	if err != nil {
		log.Fatalf("acquisition error: %v\n", err)
	} // if

	fmt.Printf("Converting %d files to text...\n\n", len(acqResult.SuccessFiles))
	conInput := types.ConversionInput(acqResult)

	log.Println("\n\t\t------------ Starting conversion... ------------ \n\n")
	conResult, err := c.Run(ctx, &conInput)
	if err != nil {
		log.Fatalf("conversion error: %v\n", err)
	} // if
	procInput := types.ProcessingInput{
		ConvertedFiles: conResult.ConvertedFiles,
	}

	fmt.Print("Processing files with spancat...\n\n")
	log.Println("\n\t\t------------ Starting processing... ------------ \n\n")
	procResult, err := p.Run(ctx, &procInput)
	if err != nil {
		log.Fatalf("processing error: %v\n", err)
	} // if

	resultHeader := "[SUCCESS]"
	resultMsg := "pipeline executed successfully"
	var failed bool
	if len(procResult.FileData) == 0 {
		resultHeader = "[FAIL]"
		resultMsg = "failed to process any files"
		failed = true
	} // if

	log.Printf("\n\n------------------------- %s -------------------------\n", resultHeader)
	log.Printf("\t\t%s\n\n", resultMsg)
	log.Println("-------------------------- [REPORT] ---------------------------")
	log.Printf("\t\t      %d files requested\n", fc)

	log.Printf("\t\t      %d files retrieved\n", len(acqResult.SuccessFiles))

	log.Printf("\t\t      %d files converted to text\n", len(conResult.ConvertedFiles))

	log.Printf("\t\t      %d files processed\n", len(procResult.FileData))

	if failed {
		fmt.Println("Pipeline failed: please rerun with a different topic or smaller filecount")
		os.Exit(1)
	}

	fmt.Print("Pipeline executed successfully\n\n")

	// logNLP(procResult, logFile)
	// Write results to file
	resultDir := "claimex-results"
	resultPath := fmt.Sprintf("%s/%s.json", resultDir, resultsFileName)

	if err := os.MkdirAll(resultDir, 0755); err != nil {
		log.Fatalf("failed to create results directory: %v", err)
	}
	output := buildOutput(procResult)
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatalf("failed to encode result JSON: %v", err)
	}

	if err := os.WriteFile(resultPath, data, 0644); err != nil {
		log.Fatalf("failed to write results file: %v", err)
	}

	fmt.Printf("Results written to %s\n", resultPath)
} // main

func buildOutput(nlpResult *processing.NLPResult) []map[string]any {
	output := []map[string]any{}
	fileData := nlpResult.FileData
	for _, data := range fileData {
		spans := []map[string]any{}
		for _, spanData := range data.ClaimSpans {
			spans = append(spans, map[string]any{
				"type":       spanData.Type,
				"text":       spanData.Text,
				"sent":       spanData.Sent,
				"confidence": spanData.Confidence,
			})
		}

		entry := map[string]any{
			"fileName":   data.FileName,
			"objectKey":  data.ObjectKey,
			"claimScore": data.ClaimScore,
			"claimSpans": spans,
		}
		output = append(output, entry)
	} // for

	return output
}
func logNLP(nr *processing.NLPResult, logFile *os.File) {
	log.SetOutput(logFile)
	fileData := nr.FileData
	for _, data := range fileData {
		log.Printf("\nSpan data for %s ---------------------------\n", data.FileName)
		for _, spanData := range data.ClaimSpans {
			fmt.Printf("\nType: %s\n", spanData.Type)
			fmt.Printf("\nText: %s\n", spanData.Text)
			fmt.Printf("\nSent: %s\n", spanData.Sent)
			fmt.Printf("\nConfidence: %f\n", spanData.Confidence)
		} // for
		log.Println("\nFile data ---------------------------")
		log.Printf("\nObject Key: %s\n", data.ObjectKey)
		log.Printf("\nFile Name: %s\n", data.FileName)
		log.Printf("\nClaim Score: %f\n", data.ClaimScore)
	} // for
} // logNLPResult
