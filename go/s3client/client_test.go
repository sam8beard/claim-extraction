package s3client

import (
	"testing"
	"fmt"
)

func TestNewClient(t *testing.T) { 
	result, err := NewClient()
	if err != nil { 
		t.Log(err)
	}
} // TestNewClient