package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct{
	fileserverHits atomic.Int32
}

func main(){
	mux:=http.NewServeMux();
	apiCfg:=apiConfig{
		fileserverHits: atomic.Int32{},
	}
	
	mux.Handle("/app/",apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./app")))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("POST /api/validate_chirp",handlerChirpsValidate)
	server:=&http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	log.Printf("Server starting on port %s", server.Addr)
	err:=server.ListenAndServe()
	if err!=nil{
		log.Fatalf("Error starting server: %v", err)
	}
}

