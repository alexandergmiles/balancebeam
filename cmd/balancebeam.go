package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexandergmiles/balancebeam/internal/balance"
)

func main() {
	fmt.Println("Starting up balancebeam load balancer")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	beam := balance.NewBeam(*logger)

	configureRoutes(beam)

	server := http.Server{Addr: ":8000", Handler: beam.Router}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go startServer(&server)
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		panic(err)
	}
	fmt.Println("Shutdown gracefully")
}

func startServer(server *http.Server) {
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		fmt.Println(err)
	}
}

func configureRoutes(beam *balance.Beam) error {
	if beam == nil {
		return fmt.Errorf("beam is nil")
	}

	beam.Register("/google", "echo.free.beeceptor.com")
	beam.Register("/google", "localhost")
	beam.Register("/index", "echo.free.beeceptor.com")

	return nil
}
