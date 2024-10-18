package main

import (
	"github.com/pa-pe/luca-control/src"
	"github.com/pa-pe/luca-control/src/service"
	"github.com/pa-pe/luca-control/src/storage"
	"github.com/pa-pe/luca-control/src/web"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type config struct {
	logDir     string
	logFile    string
	tgBotToken string
}

func main() {
	curDateTimeStr := time.Now().Format("2006-01-02 15:04:05")
	cfg := config{logDir: "./logs", logFile: curDateTimeStr + ".log"}

	// setup LogDir
	if _, err := os.Stat(cfg.logDir); os.IsNotExist(err) {
		err := os.MkdirAll(cfg.logDir, 0755)
		if err != nil {
			log.Fatalf("Error creating log directory %s: %v", cfg.logDir, err)
		}
	}

	// setup LogFile
	logPath := filepath.Join(cfg.logDir, cfg.logFile)
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		log.Fatalf("Error opening log file %s: %v", logPath, err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println("--------------------------")
	log.Println("Starting...")

	// Channel for intercepting system signals
	// interception of signals is needed to log the termination of the application
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Waiting for a system signal in a goroutine
	go func() {
		sig := <-sigs
		log.Printf("Received signal: %v. Shutting down.", sig)
		os.Exit(0)
	}()

	cfg.tgBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	if cfg.tgBotToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is empty")
	}

	dbConn, err := gorm.Open(sqlite.Open("luca_control.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error db connection: %v", err)
	}

	db, err := dbConn.DB()
	if err != nil {
		log.Fatalf("Error getting SQL DB: %v", err)
	}

	storage.RunMigrations(db)

	storages := storage.NewStorages(dbConn)
	services := service.NewServices(storages)

	tgBot := src.NewTelegramBot(cfg.tgBotToken, services)
	go tgBot.ListenAndServ()

	router := gin.Default()
	web.SetupRoutes(router, dbConn)

	log.Fatal(router.Run(":8080"))
}
