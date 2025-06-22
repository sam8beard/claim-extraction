package s3

import (
	"testing"
	"fmt"
)

func TestNewClient(t *testing.T) { 
	result := NewClient()
	fmt.Printf("%T\n", result)
} // TestNewClient