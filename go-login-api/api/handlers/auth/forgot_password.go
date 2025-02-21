package auth

import (
	"crypto/tls"
	"fmt"
	"go-login-api/config"
	"go-login-api/models"
	"math/rand"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// Fungsi untuk generate OTP 6 digit
func generateOTP() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // Generator lokal
	return fmt.Sprintf("%06d", rng.Intn(1000000))
}

// Fungsi untuk mengirim email OTP
// Fungsi untuk mengirim email OTP dengan TLS
// Fungsi untuk mengirim email OTP dengan STARTTLS
func sendEmail(to string, otp string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// Debugging: Cek apakah nilai SMTP dari .env sudah benar
	//fmt.Println("SMTP Debugging:")
	//fmt.Println("SMTP_EMAIL:", from)
	//fmt.Println("SMTP_PASSWORD:", password)
	//fmt.Println("SMTP_HOST:", smtpHost)
	//fmt.Println("SMTP_PORT:", smtpPort)

	// Buat koneksi ke server SMTP tanpa TLS awal
	conn, err := net.Dial("tcp", smtpHost+":"+smtpPort)
	if err != nil {
		fmt.Println("Connection Error:", err)
		return err
	}

	// Buat client SMTP
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		fmt.Println("SMTP Client Error:", err)
		return err
	}

	// Mulai STARTTLS
	tlsConfig := &tls.Config{ServerName: smtpHost}
	if err = client.StartTLS(tlsConfig); err != nil {
		fmt.Println("STARTTLS Error:", err)
		return err
	}

	// Setup autentikasi
	auth := smtp.PlainAuth("", from, password, smtpHost)
	if err = client.Auth(auth); err != nil {
		fmt.Println("SMTP Auth Error:", err)
		return err
	}

	// Setup sender dan receiver
	if err = client.Mail(from); err != nil {
		fmt.Println("Mail Error:", err)
		return err
	}
	if err = client.Rcpt(to); err != nil {
		fmt.Println("Recipient Error:", err)
		return err
	}

	// Kirim email
	wc, err := client.Data()
	if err != nil {
		fmt.Println("Data Error:", err)
		return err
	}
	message := []byte("To: " + to + "\r\n" +
		"Subject: Password Reset OTP\r\n\r\n" +
		"Your OTP for password reset is: " + otp + "\r\n")
	_, err = wc.Write(message)
	if err != nil {
		fmt.Println("Write Message Error:", err)
		return err
	}
	err = wc.Close()
	if err != nil {
		fmt.Println("Close Write Error:", err)
		return err
	}

	// Tutup koneksi
	client.Quit()
	fmt.Println("Email sent successfully!")
	return nil
}

// Handler untuk Forgot Password
func ForgotPasswordHandler(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Cek apakah user ada di database
	var user models.User
	result := config.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Hapus OTP lama jika ada
	config.DB.Unscoped().Where("email = ?", req.Email).Delete(&models.PasswordReset{})

	// Generate OTP baru dan simpan ke database
	otp := generateOTP()
	passwordReset := models.PasswordReset{
		Email:     req.Email,
		OTP:       otp,
		ExpiresAt: time.Now().Add(10 * time.Minute), // OTP berlaku 10 menit
	}
	if err := config.DB.Create(&passwordReset).Error; err != nil {
		fmt.Println("Error saving OTP to database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store OTP"})
		return
	}

	// Kirim email OTP ke user
	err := sendEmail(req.Email, otp)
	if err != nil {
		fmt.Println("Error sending email:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to email"})
}
