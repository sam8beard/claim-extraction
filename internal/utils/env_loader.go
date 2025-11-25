package utils

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadDotEnvUpwards() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	} // if

	for {
		envPath := filepath.Join(dir, ".env")
		_, err := os.Stat(envPath)
		// fmt.Println(err)
		if err == nil {
			// fmt.Println("here")
			// fmt.Println(envPath)
			return godotenv.Load(envPath)
		} // if
		parent := filepath.Dir(dir)
		if parent == dir {
			// fmt.Println("here")
			break // we have reached the root dir
		} // if
		dir = parent
	} // for
	// fmt.Println(dir)
	return os.ErrNotExist
} // LoadDotEnvUpwards
