package main

import (
	"github.com/pa-pe/luca-control/config"
	"github.com/pa-pe/luca-control/src"
	"github.com/pa-pe/luca-control/src/service"
	"github.com/pa-pe/luca-control/src/storage"
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

func main() {
	curDateTimeStr := time.Now().Format(config.LogFileFormat)
	logFileName := curDateTimeStr + ".log"

	// setup LogDir
	if _, err := os.Stat(config.LogDir); os.IsNotExist(err) {
		err := os.MkdirAll(config.LogDir, 0755)
		if err != nil {
			log.Fatalf("Error creating log directory %s: %v", config.LogDir, err)
		}
	}

	// setup LogFile
	logPath := filepath.Join(config.LogDir, logFileName)
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

	tgBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if tgBotToken == "" {
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

	tgBot := src.NewTelegramBot(tgBotToken)

	storages := storage.NewStorages(dbConn)
	services := service.NewServices(storages, tgBot)

	tgBot.Handler = services.ChatBotMsgRouter.Handle
	go tgBot.ListenAndServ()

	router := gin.Default()
	src.SetupRoutes(router, dbConn)

	log.Fatal(router.Run(":" + config.WebServerPort))
}
