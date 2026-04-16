package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	httpadapter "owner-api-proxy/internal/adapter/http"
	"owner-api-proxy/internal/config"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

func ServerCommand() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := config.LoadConfig(&configPath); err != nil {
				return err
			}

			sqlClient := config.NewClientSql("postgres")
			if err := sqlClient.Open(); err != nil {
				return err
			}
			defer sqlClient.Close()

			log := config.NewLogrusWithCategory("server")
			ginCfg := config.LoadGinConfig()
			engine := config.NewGinEngine(ginCfg)
			httpadapter.RegisterRoutes(engine)

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
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "app.yaml", "Config file path")

	return cmd
}
