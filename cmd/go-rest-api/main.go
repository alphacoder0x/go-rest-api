package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alphacoder0x/go-rest-api/internal/config"
)

func main() {
	//load config

	cfg := config.MustLoad()
	//database setup
	//setup router

	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to go rest api"))
	})

	router.HandleFunc("GET /api/data",func(w http.ResponseWriter,r *http.Request){
		w.Write([]byte("serving for data"))
	})
	//setup http server

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	done := make(chan os.Signal,1)

	signal.Notify(done,os.Interrupt,syscall.SIGINT,syscall.SIGTERM)
	go func(){
		fmt.Println("server started", cfg.Addr)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Serving http return error %s", err.Error())
		}
	}()

	<-done

	slog.Info("shutting down the server")
	
	ctx,cancel := context.WithTimeout(context.Background(),5 * time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server",slog.String("error",err.Error()))
	}

	slog.Info("server shutdown successfully")
	
}
