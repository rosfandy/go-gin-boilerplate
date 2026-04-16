package command

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"owner-api-proxy/cmd/bootstrap"
	"owner-api-proxy/internal/config"
	"owner-api-proxy/pkg/cli/server"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

func ServerCommand() *cobra.Command {
	runServer := func(configPath string) error {
		if _, err := config.LoadConfig(&configPath); err != nil {
			return err
		}

		infoLog := config.NewLogrusWithCategory("info")
		sqlClient := config.NewClientSql("mysql")
		if err := sqlClient.Open(); err != nil {
			return err
		}
		defer sqlClient.Close()

		log := config.NewLogrusWithCategory("server")
		ginCfg := config.LoadGinConfig()
		engine := config.NewGinEngine(ginCfg)
		bootstrap.BuildHTTPEngine(engine, sqlClient.Client())

		srv := &http.Server{
			Addr:    ginCfg.Address,
			Handler: engine,
		}

		errCh := make(chan error, 1)
		go func() {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
			close(errCh)
		}()

		infoLog.Info("database connected")
		log.WithField("address", ginCfg.Address).Info("server started")

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigCh)

		select {
		case err := <-errCh:
			if err != nil {
				return err
			}
			return nil
		case sig := <-sigCh:
			log.WithField("signal", sig.String()).Info("server shutdown requested")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return err
		}

		log.Info("server stopped")

		return nil
	}

	return server.ServerCommand(runServer)
}
