package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	qrcode "github.com/skip2/go-qrcode"
	"gopkg.in/gomail.v2"
)

// Email configurations
const (
	smtpServer   = "YOUR_SMTP_SERVER"
	smtpPort     = 587
	smtpUsername = "YOUR_SMTP_USERNAME"
	smtpPassword = "YOUR_SMTP_PASSWORD"
)

// generateCode generates a random 6-digit code.
func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(999999))
}

// sendEmail sends the 2FA code to the user's email.
func sendEmail(email, code string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", smtpUsername)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Two-Factor Authentication Code")
	message.SetBody("text/plain", "Your 2FA code is: "+code)

	dialer := gomail.NewDialer(smtpServer, smtpPort, smtpUsername, smtpPassword)

	if err := dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}

func TwoFAByMail() {
	// For demonstration purposes, let's assume the user's email is known.
	userEmail := "user@example.com"

	// Generate the 2FA code
	code := generateCode()

	// Send the code via email
	err := sendEmail(userEmail, code)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return
	}

	fmt.Println("2FA code sent to", userEmail)

	// Simulate user input (replace this with your actual user input handling)
	var userEnteredCode string
	fmt.Print("Enter the 2FA code: ")
	fmt.Scan(&userEnteredCode)

	// Validate the entered code
	if userEnteredCode == code {
		fmt.Println("Authentication successful!")
	} else {
		fmt.Println("Authentication failed. Invalid code.")
	}
}

//****************2FA By Google api*********************

// generateTOTP generates a Time-based One-Time Password (TOTP).
func generateTOTP(secret string) (string, error) {
	key, err := otp.NewKeyFromURL(secret)
	if err != nil {
		return "", err
	}

	return totp.GenerateCode(key.Secret(), time.Now())
}

func TwoFAByGoogleApi() {
	// Replace "YourAppName" with the name of your application.
	// This will be displayed in the user's 2FA app when they scan the QR code.
	appName := "YourAppName"

	// Generate a new TOTP secret for the user.
	secret, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      appName,
		AccountName: "user@example.com", // Replace with the user's email or username.
	})

	// Convert the TOTP secret to a URL that can be used to generate QR code.
	qrCodeURL := secret.URL()

	fmt.Println("Scan the following QR code using Google Authenticator or Authy:")
	fmt.Println(qrCodeURL)

	// Generate the QR code image and save it to a file.
	qrCodeImg, err := qrcode.Encode(qrCodeURL, qrcode.Medium, 256)
	if err != nil {
		fmt.Println("Error generating QR code:", err)
		os.Exit(1)
	}

	// Save the QR code image to a file (optional).
	file, err := os.Create("qrcode.png")
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
	defer file.Close()

	_, err = file.Write(qrCodeImg)
	if err != nil {
		fmt.Println("Error writing QR code to file:", err)
		os.Exit(1)
	}

	// Wait for the user to scan the QR code and configure their 2FA app.
	// Replace this with your actual user input handling.
	fmt.Print("Press Enter when ready...")
	fmt.Scanln()

	// Verify the entered code.
	if _, err := generateTOTP(secret.Secret()); err != nil {
		fmt.Println("Error generating TOTP code:", err)
		os.Exit(1)
	} else {
		// Simulate user input (replace this with your actual user input handling).
		var userEnteredCode string
		fmt.Print("Enter the 2FA code from your app: ")
		fmt.Scan(&userEnteredCode)

		// Verify the entered code.
		if totp.Validate(userEnteredCode, secret.Secret()) {
			fmt.Println("Authentication successful!")
		} else {
			fmt.Println("Authentication failed. Invalid code.")
		}
	}
}

func main() {
	//TwoFAByMail()
	TwoFAByGoogleApi()
}
