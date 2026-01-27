package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Starting up balancebeam load balancer")

	beam := balance.NewRouter()

	go startServer(&server)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	shutdownContext, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()

	if err := server.Shutdown(shutdownContext); err != nil {
		panic(err)
	}
	fmt.Println("Shutdown successfully")
}

func startServer(server *http.Server) {
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
