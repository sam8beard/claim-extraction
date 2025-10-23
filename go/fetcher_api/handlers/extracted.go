package handlers 

import ( 
	"net/http"
	"encoding/json"
	// "github.com/minio/minio-go/v7"
	// "github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sam8beard/claim-extraction/go/db"
	"github.com/sam8beard/claim-extraction/go/utils"
	"os"
	"fmt"
	"strings"
	"path"
	"context"
)

type Response struct { 
	Keys []string `json:"keys"`
} // Response

func ExtractedHandler(w http.ResponseWriter, r *http.Request) { 
	// load env vars
	err := utils.LoadDotEnvUpwards()
	if err != nil { 
		fmt.Println("Could not load .env variables")
		panic(err)
	} // if 

	resp := Response{}

	if r.Method != "GET" { 
		http.Error(w, "Only GET requests allowed ", http.StatusMethodNotAllowed)
		return
	} // if 
	
	ctx := context.Background()

	// establish connection pool to pg db
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil { 
			fmt.Println("Unable to establish database connection")
			panic(err)
	} // if 
	defer pool.Close()

	// // create MinIO client
	// endpoint := "localhost:9000"
	// accessKeyID := "muel"
	// secretAccessKey := "password"
	// useSSL := false 
	// bucketName := "claim-pipeline-docstore"

	// minioClient, err := minio.New(endpoint, &minio.Options{ 
	// 	Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	// 	Secure: useSSL,
	// })

	if err != nil {
		fmt.Println("Unable to establish connection to MinIO server")
		panic(err)
	} // if 
	

	// get all keys of files whose text has been extracted 
	rows, err := db.GetAllExtractedKeys(ctx, pool)
	if err != nil {
		fmt.Println("Error scanning row: ", err)
		panic(err)
	} // if 

	keys := make([]string, 0)

	for rows.Next() { 
		var key string 
		err := rows.Scan(&key)
		if err != nil { 
			fmt.Println("Error scanning key:", err)
			continue
		} // if
		ext := path.Ext(key) 
		newKey := strings.Replace(key, ext, ".txt", 1)
		newKey = strings.Replace(newKey, "raw", "processed", 1)
		keys = append(keys, newKey)
	} // for

	// resp.Keys, err := json.Marshal(keys)
	// if err != nil { 
	// 	fmt.Println("Could not encode keys")
	// 	panice(err)
	// } // if 

	resp.Keys = keys

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

} // ExtractedHandler