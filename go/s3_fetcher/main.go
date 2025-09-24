package main

import ( 
	"github.com/sam8beard/claim-extraction/go/db"
	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sam8beard/claim-extraction/go/utils"
	"context"
	"os"
	"fmt"
	"path"
	"strings"
	"encoding/json"
)

func main() { 
	// load env vars
	err := utils.LoadDotEnvUpwards()
	if err != nil { 
		fmt.Println("Could not load .env variables")
		return 
	} // if 

	// establish connection pool to pg db
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("Unable to establish database connection")
		return
	} // if 
	defer pool.Close() 

	// get rows that have extracted text available
	rows, err := db.GetAllExtractedKeys(context.Background(), pool)
	
	keys := make([]string, 0)

	// create json file for python 
	file, err := os.Create("../../python/nlp/training/s3-keys.json")
	if err != nil { 
		fmt.Println("Error creating file:", err)
		return 
	} // if 
	defer file.Close()
	
	// modify keys to .txt format e.g. path/to/file.pdf -> path/to/file.txt

	// iterate through rows 
		// scan key into string var 
		// replace extension with .txt
		// add to .csv file
	for rows.Next() { 

		var key string 
		rows.Scan(&key)
		ext := path.Ext(key) 
		newKey := strings.Replace(key, ext, ".txt", 1)
		newKey = strings.Replace(newKey, "raw", "processed", 1)
		keys = append(keys, newKey)

	} // for
	
	// marshal the keys 
	json_keys, err := json.Marshal(keys)
	if err != nil { 
		fmt.Println("Could not encode to json")
		return
	} // if 
	
	// write json file
	_, err = file.Write(json_keys)
	if err != nil { 
		fmt.Println("Could not write to file")
		return
	} // if 

} // main

// will this structure work:
// - make slice of byte arrays
// - iterate through rows
// - for each row: make string var, scan row contents into string var, replace file extention in string var with .txt, encode string, add to slice