package main

import (
	"go-mail/adapters"
	"go-mail/constants"
	"go-mail/controllers"
	"go-mail/db"
	"go-mail/repositories"
	"go-mail/services"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or failed to load")
	}

	baseLLMUrl := os.Getenv("baseLLMUrl")
	modelName := os.Getenv("model")
	if baseLLMUrl == "" || modelName == "" {
		log.Fatal("baseLLMUrl and model must be set in .env file")
	}

	db.InitPostgres()

	// adapters
	goAuthAdapter := adapters.NewGoAuthAdapter()
	telegramAdapter := adapters.NewTelegramAdapter()
	gmailAdapter := adapters.NewGmailAdapter()
	aiAdapter := adapters.NewAIAdapter(baseLLMUrl, modelName)
	// repositories
	userRepo := repositories.NewUserRepository(db.DB)
	preferenceRepo := repositories.NewPreferenceRepository(db.DB)

	// services
	priorimailService := services.NewPriorimailService(preferenceRepo, userRepo, telegramAdapter, gmailAdapter, aiAdapter, goAuthAdapter)

	// controllers
	priorimailController := controllers.NewPriorimailController(priorimailService)

	var arg string
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}

	if arg != "" {
		filename := "logs/" + time.Now().Format("2006-01-02T15-04-05") + ".log"
		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("Failed to create log file: %v", err)
		}
		defer file.Close()

		switch arg {
		case "processEmails":
			priorimailController.ProcessEmails(filename)
		default:
			log.Fatalf("Unknown batch command: %s", arg)
		}
		return
	}

	e := echo.New()

	port := os.Getenv(constants.PORT)
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(e.Start(":" + port))
}
