package main

import (
	"cdn_app/database"
	api "cdn_app/pkg/files/http"
	fservice "cdn_app/pkg/files/service"
	fstorage "cdn_app/pkg/files/storage"
	"cdn_app/settings"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
	port := os.Getenv("PORT")
	settings.Config = settings.Settings{
		Port: port,
	}
	fmt.Printf("Application access port: %s\n", port)

	// Database init
	rsc, err := database.InitDB()
	if err != nil {
		fmt.Printf("Can't initialize resources. %v\n", err)
	}
	defer func() {
		err := rsc.Release()
		if err != nil {
			fmt.Printf("Got an error during resources release. %v\n", err)
		}
	}()

	filesDb := fstorage.New(rsc.DB)
	filesCtrl := fservice.NewController(filesDb)


	r := mux.NewRouter()
	r.HandleFunc("/api/files/user/{name}", api.APIHandler(filesCtrl)).Methods("GET")
	r.HandleFunc("/api/files/server/{name}", api.APIHandler(filesCtrl)).Methods("GET")
	r.HandleFunc("/api/files/area/{name}", api.APIHandler(filesCtrl)).Methods("GET")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	shutdown := make(chan error, 1)

	server := http.Server{
		Addr:    net.JoinHostPort("", settings.Config.Port),
		Handler: r,
	}

	go func() {
		err := server.ListenAndServe()
		fmt.Printf("Got an error during ListenAndServe: %v\n", err)
		shutdown <- err
	}()
	fmt.Println("The service is ready to listen and serve")

	select {
	case killSignal := <-interrupt:
		switch killSignal {
		case os.Interrupt:
			fmt.Println("Got SIGINT...")
		case syscall.SIGTERM:
			fmt.Println("Got SIGTERM...")
		}
	case <-shutdown:
		fmt.Println("Got an error...")
	}

	fmt.Println("The service is stopping...")
	err = server.Shutdown(context.Background())
	if err != nil {
		fmt.Printf("Got an error during service shutdown: %v\n", err)
	}

}
