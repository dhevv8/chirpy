package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/dhevv8/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db 		  *database.Queries
	platform string
	jwtSecret string
}

func main(){
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform:=os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	dbConn,err:=sql.Open("postgres",dbURL)
	if err!=nil{
		log.Fatalf("Error connecting to database: %s", err)
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	dbQueries:=database.New(dbConn)
	mux:=http.NewServeMux();
	apiCfg:=apiConfig{
		fileserverHits: atomic.Int32{},
		db:	dbQueries,
	}
	
	mux.Handle("/app/",apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./app")))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsers)
	mux.HandleFunc("POST /api/chirps", apiCfg.handleChirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpByID)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	server:=&http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	log.Printf("Server starting on port %s", server.Addr)
	err=server.ListenAndServe()
	if err!=nil{
		log.Fatalf("Error starting server: %v", err)
	}
}

