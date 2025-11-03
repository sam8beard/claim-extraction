package types

import (
	"github.com/sam8beard/claim-extraction/api-refactor/internal/types/shared"
)

type ProcessingInput struct {
	ConvertedFiles []shared.File
}
