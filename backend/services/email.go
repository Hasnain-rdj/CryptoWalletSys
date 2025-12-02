package services

import (
	"fmt"
	"net/smtp"
	"os"
)

// EmailConfig holds SMTP configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// GetEmailConfig returns email configuration from environment
func GetEmailConfig() EmailConfig {
	return EmailConfig{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPUser:     os.Getenv("SMTP_USER"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromEmail:    os.Getenv("FROM_EMAIL"),
		FromName:     os.Getenv("FROM_NAME"),
	}
}

// SendOTPEmail sends OTP to user's email
func SendOTPEmail(toEmail, otp string) error {
	config := GetEmailConfig()

	// Check if email is configured
	if config.SMTPHost == "" || config.SMTPUser == "" || config.SMTPPassword == "" {
		// Fallback to console logging for development
		fmt.Printf("===========================================\n")
		fmt.Printf("Email Configuration Not Found - Development Mode\n")
		fmt.Printf("OTP for %s: %s\n", toEmail, otp)
		fmt.Printf("===========================================\n")
		return nil
	}

	// Setup authentication
	auth := smtp.PlainAuth("", config.SMTPUser, config.SMTPPassword, config.SMTPHost)

	// Compose email
	subject := "Your Blockchain Wallet OTP"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 10px 10px; }
        .otp-box { background: white; border: 2px dashed #667eea; padding: 20px; text-align: center; margin: 20px 0; border-radius: 8px; }
        .otp-code { font-size: 32px; font-weight: bold; color: #667eea; letter-spacing: 8px; }
        .footer { text-align: center; margin-top: 20px; color: #666; font-size: 12px; }
        .warning { background: #fff3cd; border-left: 4px solid #ffc107; padding: 12px; margin: 15px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔐 Blockchain Wallet</h1>
            <p>Email Verification</p>
        </div>
        <div class="content">
            <h2>Hello,</h2>
            <p>You requested an OTP to verify your email address for your Blockchain Wallet account.</p>
            
            <div class="otp-box">
                <p style="margin: 0; font-size: 14px; color: #666;">Your OTP Code</p>
                <div class="otp-code">%s</div>
                <p style="margin: 10px 0 0 0; font-size: 12px; color: #999;">Valid for 5 minutes</p>
            </div>

            <div class="warning">
                <strong>⚠️ Security Notice:</strong><br>
                Never share this code with anyone. Our team will never ask for your OTP.
            </div>

            <p>If you didn't request this code, please ignore this email or contact support if you have concerns.</p>
        </div>
        <div class="footer">
            <p>This is an automated email. Please do not reply.</p>
            <p>&copy; 2025 Blockchain Wallet. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, otp)

	// Email message
	message := []byte(
		"From: " + config.FromName + " <" + config.FromEmail + ">\r\n" +
			"To: " + toEmail + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"\r\n" +
			body + "\r\n")

	// Send email
	addr := config.SMTPHost + ":" + config.SMTPPort
	err := smtp.SendMail(addr, auth, config.FromEmail, []string{toEmail}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	fmt.Printf("OTP email sent successfully to %s\n", toEmail)
	return nil
}

// SendWelcomeEmail sends welcome email to new user
func SendWelcomeEmail(toEmail, fullName string) error {
	config := GetEmailConfig()

	if config.SMTPHost == "" || config.SMTPUser == "" || config.SMTPPassword == "" {
		fmt.Printf("Welcome email would be sent to: %s (%s)\n", toEmail, fullName)
		return nil
	}

	auth := smtp.PlainAuth("", config.SMTPUser, config.SMTPPassword, config.SMTPHost)

	subject := "Welcome to Blockchain Wallet!"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 10px 10px; }
        .feature { background: white; padding: 15px; margin: 10px 0; border-radius: 5px; border-left: 4px solid #667eea; }
        .footer { text-align: center; margin-top: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 Welcome to Blockchain Wallet!</h1>
        </div>
        <div class="content">
            <h2>Hello %s,</h2>
            <p>Congratulations! Your blockchain wallet has been created successfully.</p>
            
            <h3>What you can do:</h3>
            <div class="feature">
                <strong>💸 Send Money</strong><br>
                Transfer cryptocurrency to other wallet holders securely
            </div>
            <div class="feature">
                <strong>📊 Track Transactions</strong><br>
                View your complete transaction history on the blockchain
            </div>
            <div class="feature">
                <strong>🔍 Blockchain Explorer</strong><br>
                Explore blocks and verify transactions
            </div>
            <div class="feature">
                <strong>🕌 Automatic Zakat</strong><br>
                2.5%% monthly Zakat deduction managed automatically
            </div>

            <p><strong>⚠️ Important Security Reminders:</strong></p>
            <ul>
                <li>Keep your private key secure and never share it</li>
                <li>Enable two-factor authentication</li>
                <li>Use a strong, unique password</li>
                <li>Beware of phishing attempts</li>
            </ul>
        </div>
        <div class="footer">
            <p>&copy; 2025 Blockchain Wallet. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, fullName)

	message := []byte(
		"From: " + config.FromName + " <" + config.FromEmail + ">\r\n" +
			"To: " + toEmail + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"\r\n" +
			body + "\r\n")

	addr := config.SMTPHost + ":" + config.SMTPPort
	err := smtp.SendMail(addr, auth, config.FromEmail, []string{toEmail}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
