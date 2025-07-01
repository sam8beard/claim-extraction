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
		if err == nil { 
			return godotenv.Load(envPath)
		} // if 
		parent := filepath.Dir(dir)
		if parent == dir { 
			break // we have reached the root dir
		} // if 
		dir = parent
	} // for 

	return os.ErrNotExist
} // LoadDotEnvUpwards