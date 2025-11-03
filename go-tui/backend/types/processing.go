package types

import (
	"bytes"
	"tui/backend/types/shared"
)

type ProcessingInput struct {
	ConvertedFiles []shared.File
}

type ProcessingResult struct {
	ProcessedFiles map[shared.File]bytes.Buffer
}
