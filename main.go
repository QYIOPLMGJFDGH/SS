package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/zelenin/go-tdlib/client"
)

func main() {
	// Initialize TDLib logging
	client.SetLogVerbosityLevel(1)
	client.SetFilePath("./tdlib.log")

	// User inputs for API ID and API Hash
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your Telegram API ID: ")
	apiIDInput, _ := reader.ReadString('\n')
	apiID, err := strconv.Atoi(strings.TrimSpace(apiIDInput))
	if err != nil {
		log.Fatalf("Invalid API ID: %v", err)
	}

	fmt.Print("Enter your Telegram API Hash: ")
	apiHashInput, _ := reader.ReadString('\n')
	apiHash := strings.TrimSpace(apiHashInput)

	// TDLib client configuration
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
		log.Fatalf("Failed to create TDLib client: %v", err)
	}

	fmt.Println("Telegram client initialized.")

	// User inputs for authentication
	fmt.Print("Enter your phone number (with country code): ")
	phoneInput, _ := reader.ReadString('\n')
	phone := strings.TrimSpace(phoneInput)

	// Send phone number
	_, err = tdlibClient.AuthSendPhoneNumber(phone)
	if err != nil {
		log.Fatalf("Error sending phone number: %v", err)
	}

	// Enter OTP
	fmt.Print("Enter the OTP sent to your phone: ")
	otpInput, _ := reader.ReadString('\n')
	otp := strings.TrimSpace(otpInput)

	// Verify OTP
	_, err = tdlibClient.AuthSendCode(otp)
	if err != nil {
		log.Fatalf("Error verifying OTP: %v", err)
	}

	// Two-step password handling
	fmt.Print("Enter your two-step verification password (leave blank if not set): ")
	passwordInput, _ := reader.ReadString('\n')
	password := strings.TrimSpace(passwordInput)

	if password != "" {
		_, err = tdlibClient.AuthSendPassword(password)
		if err != nil {
			log.Fatalf("Error verifying two-step password: %v", err)
		}
	}

	// Generate string session
	session, err := tdlibClient.GetStringSession()
	if err != nil {
		log.Fatalf("Failed to generate string session: %v", err)
	}

	// Output the string session
	fmt.Println("\nYour Telegram String Session:")
	fmt.Println(session)
	fmt.Println("\nSave this string session securely. Do not share it with anyone!")

	// Cleanup
	tdlibClient.Close()
}
