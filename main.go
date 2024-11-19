package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zelenin/go-tdlib/client"
)

var (
	BotToken string
)

func main() {
	// Load environment variables
	BotToken = os.Getenv("BOT_TOKEN", "7426075639:AAE854r2874ZJAVat6zVUeSR4IYBnsW-y-w") // Telegram Bot Token
	if BotToken == "" {
		log.Fatal("Error: BOT_TOKEN environment variable not set.")
	}

	// Initialize Telegram Bot
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot: %v", err)
	}
	bot.Debug = true
	fmt.Printf("Authorized on bot: %s\n", bot.Self.UserName)

	// Goroutine to handle Telegram bot updates
	go func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates := bot.GetUpdatesChan(u)

		for update := range updates {
			if update.Message != nil {
				handleMessage(bot, update.Message)
			} else if update.CallbackQuery != nil {
				handleCallbackQuery(bot, update.CallbackQuery)
			}
		}
	}()

	// Start HTTP server
	fmt.Println("Starting HTTP server on port 8080...")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bot is running!"))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID

	switch msg.Text {
	case "/start":
		// Send welcome message with buttons
		welcomeMsg := "Welcome to the String Session Generator Bot! Select a session type below to begin."
		buttons := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Generate v1 Session", "generate_v1"),
				tgbotapi.NewInlineKeyboardButtonData("Generate v2 Session", "generate_v2"),
			),
		)
		msg := tgbotapi.NewMessage(chatID, welcomeMsg)
		msg.ReplyMarkup = buttons
		bot.Send(msg)

	default:
		// Handle unknown commands or messages
		bot.Send(tgbotapi.NewMessage(chatID, "Unknown command. Please use /start to begin."))
	}
}

func handleCallbackQuery(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) {
	chatID := query.Message.Chat.ID
	data := query.Data

	// Acknowledge the callback query
	bot.AnswerCallbackQuery(tgbotapi.NewCallback(query.ID, ""))

	switch data {
	case "generate_v1", "generate_v2":
		// Ask for API ID and API Hash
		sessionType := "v1"
		if data == "generate_v2" {
			sessionType = "v2"
		}
		message := fmt.Sprintf("You selected %s session.\n\nPlease send your API ID:", sessionType)
		bot.Send(tgbotapi.NewMessage(chatID, message))

		// Handle user inputs for API ID, API Hash, and phone number in a separate flow
		go waitForUserInputs(bot, chatID, sessionType)
	}
}

func waitForUserInputs(bot *tgbotapi.BotAPI, chatID int64, sessionType string) {
	// Simulate user input handling
	reader := strings.NewReader("") // Replace with actual input mechanism

	// Get API ID
	fmt.Print("Enter your Telegram API ID: ")
	apiIDInput := readLine(reader)
	apiID, err := strconv.Atoi(strings.TrimSpace(apiIDInput))
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Invalid API ID. Please try again."))
		return
	}

	// Get API Hash
	fmt.Print("Enter your Telegram API Hash: ")
	apiHashInput := readLine(reader)
	apiHash := strings.TrimSpace(apiHashInput)

	// Get Phone Number
	fmt.Print("Enter your phone number (with country code): ")
	phoneInput := readLine(reader)
	phone := strings.TrimSpace(phoneInput)

	// Generate the string session
	stringSession, err := generateStringSession(apiID, apiHash, phone, sessionType)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Failed to generate string session. Please try again."))
		return
	}

	// Send the string session back to the user
	response := fmt.Sprintf("Your %s string session:\n\n`%s`\n\nSave it securely!", sessionType, stringSession)
	bot.Send(tgbotapi.NewMessage(chatID, response))
}

func generateStringSession(apiID int, apiHash, phone, sessionType string) (string, error) {
	// TDLib configuration
	tdlibConfig := &client.Config{
		APIID:              apiID,
		APIHash:            apiHash,
		SystemLanguageCode: "en",
		DeviceModel:        "Golang Bot",
		ApplicationVersion: "1.0",
		DatabaseDirectory:  "./tdlib-db",
		FilesDirectory:     "./tdlib-files",
	}

	// Create TDLib client
	tdlibClient, err := client.NewClient(tdlibConfig)
	if err != nil {
		return "", fmt.Errorf("Failed to create TDLib client: %v", err)
	}

	// Send phone number
	_, err = tdlibClient.AuthSendPhoneNumber(phone)
	if err != nil {
		return "", fmt.Errorf("Error sending phone number: %v", err)
	}

	// Get OTP
	fmt.Print("Enter the OTP sent to your phone: ")
	otp := readLine(strings.NewReader(""))
	_, err = tdlibClient.AuthSendCode(strings.TrimSpace(otp))
	if err != nil {
		return "", fmt.Errorf("Error verifying OTP: %v", err)
	}

	// Two-step password handling
	fmt.Print("Enter your two-step verification password (leave blank if not set): ")
	password := readLine(strings.NewReader(""))
	if password != "" {
		_, err = tdlibClient.AuthSendPassword(strings.TrimSpace(password))
		if err != nil {
			return "", fmt.Errorf("Error verifying two-step password: %v", err)
		}
	}

	// Generate string session
	session, err := tdlibClient.GetStringSession()
	if err != nil {
		return "", fmt.Errorf("Failed to generate string session: %v", err)
	}

	// Cleanup
	tdlibClient.Close()

	return session, nil
}

func readLine(reader *strings.Reader) string {
	var line string
	fmt.Fscanln(reader, &line)
	return line
}
