package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"agentic/commerce/config"
	"agentic/commerce/internal/app"
	"agentic/commerce/internal/domains"
	internalhttp "agentic/commerce/internal/interfaces/http"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var serveCMD = &cobra.Command{
	Use:   "serve",
	Short: "Serve the application",
	Long:  `Serve the HTTP server of the application`,
	Run:   serve,
}

func serve(_ *cobra.Command, _ []string) {
	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		panic(err)
	}

	bootstrap := fx.New(
		fx.Supply(cfg),
		config.Module,
		app.LoggerModule,
		app.DatabaseModule,
		domains.Modules,
		internalhttp.Module,
	)

	if err := bootstrap.Start(context.Background()); err != nil {
		fmt.Println("Failed to start:", err)
		os.Exit(1)
	}

	<-signalCtx.Done()

	stopCtx, cancelStop := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelStop()

	if err := bootstrap.Stop(stopCtx); err != nil {
		fmt.Println("Failed to stop:", err)
		os.Exit(1)
	}

	fmt.Println("Application shut down successfully")
}
