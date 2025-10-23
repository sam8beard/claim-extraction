package main

import ( 
	"net/http"
	"github.com/sam8beard/claim-extraction/go/fetcher_api/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"fmt"
)

func main() { 
	router := chi.NewRouter() 
	router.Use(middleware.Logger)
	router.Get("/documents/extracted", handlers.ExtractedHandler)

	server := &http.Server { 
		Addr: ":60000",
		Handler: router,
	} // server 

	fmt.Println("Listening on port", server.Addr)

	err := server.ListenAndServe()
	if err != nil { 
		fmt.Println("Failed to start server", err)
	} // if 
	
} // main