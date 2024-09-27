package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"

	"github.com/authgear/authgear-sms-gateway/pkg/handler"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/logger"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms"
)

func main() {
	err := godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("failed to load .env file: %s", err)
	}
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		panic(err)
	}

	logger := logger.NewLogger()

	smsServiceProviderConfigPath, err := filepath.Abs(cfg.SMSServiceProviderConfigPath)
	smsProviderConfigYAML, err := os.ReadFile(smsServiceProviderConfigPath)
	if err != nil {
		panic(err)
	}
	smsProviderConfig, err := config.ParseSMSProviderConfigFromYAML([]byte(smsProviderConfigYAML))
	if err != nil {
		panic(err)
	}
	smsClientMap := sms.NewSMSClientMap(smsProviderConfig, logger)
	smsService := &sms.SMSService{
		Logger:            logger,
		SMSProviderConfig: smsProviderConfig,
		SMSClientMap:      smsClientMap,
	}

	http.Handle("/healthz", &handler.HealthzHandler{})
	http.Handle("/send", &handler.SendHandler{
		Logger:     logger,
		SMSService: smsService,
	})

	server := &http.Server{
		Addr:              cfg.ListenAddr,
		ReadHeaderTimeout: 3 * time.Second,
	}

	logger.Info(fmt.Sprintf("Server running at %v", cfg.ListenAddr))
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
