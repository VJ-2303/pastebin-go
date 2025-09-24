package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type config struct {
	addr string
}

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLog,
		Handler:  nil,
	}

	infoLog.Printf("Starting server on %s", cfg.addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
