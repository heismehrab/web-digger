package cmd

import (
	"context"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"web-digger/internal/app"
	"web-digger/internal/config"
	"web-digger/pkg/logger"
)

func RunDigger() {
	ctx, cancel := context.WithCancel(context.Background())

	// Read .env file.
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	viper.SetConfigName("web-digger")
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AddConfigPath("/etc/web-digger")
	viper.AddConfigPath("$HOME/.web-digger")
	viper.AddConfigPath(".")

	// Read config.
	err := viper.ReadInConfig()

	if err != nil {
		cancel()
		log.Fatalf("failed to read config file: %s", err)
	}

	var cfg config.Config
	err = viper.Unmarshal(&cfg)

	if err != nil {
		cancel()
		log.Fatalf("failed to unmarshal config: %s", err)
	}

	// Create Logger instance.
	loggerInstance := logger.CreateLogger(cfg.Logger)

	hostname, err := os.Hostname()

	if err != nil {
		cancel()
		loggerInstance.FatalF("failed to get container name with error: %v", err)
	}

	loggerInstance.InfoF("[OK] Hostname acquired :%s", hostname)

	// Provide Graceful shutdown functionality.
	done := make(chan bool, 1)
	quiteSignal := make(chan os.Signal, 1)

	signal.Notify(quiteSignal, syscall.SIGINT, syscall.SIGTERM)

	application := app.NewInstance(cfg)

	err = application.Bootstrap(ctx, loggerInstance)

	if err != nil {
		cancel()
		loggerInstance.FatalF("failed to bootstrap Heimdall instance with error: %v", err)
	}

	go application.GracefulShutdown(quiteSignal, done)

	// Start HTTP server (blocking).
	application.Start(ctx)

	// Wait for HTTP Server to be killed gracefully.
	<-done

	// Killing other background jobs.
	cancel()
	loggerInstance.Info("Waiting for background jobs to finish their works")

	// Wait for all other background jobs to finish their works.
	application.Wait()

	loggerInstance.Info("Master App was shutdown successfully")
}
