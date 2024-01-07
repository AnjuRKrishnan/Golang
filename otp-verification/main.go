package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"otp-verification-go/users"
)

type Users struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	PhoneNumber       string    `json:"phone_number"`
	OTP               string    `json:"otp"`
	OTPExpirationTime time.Time `json:"otp_expiration_time"`
}

func main() {
	router := gin.Default()

	// Create a PostgreSQL connection pool
	dbPool, err := pgxpool.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close()

	// Routes
	router.POST("/api/users", createUser(dbPool))
	router.POST("/api/users/generateotp", generateOTP(dbPool))
	router.POST("/api/users/verifyotp", verifyOTP(dbPool))

	// Run the server
	if err := router.Run(":8080"); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
}

func createUser(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		// Check if phone number already exists
		var existingUser User
		err := db.QueryRow(context.Background(), "SELECT id FROM users WHERE phone_number = $1", user.PhoneNumber).
			Scan(&existingUser.ID)
		if err != pgx.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already exists"})
			return
		}

		// Insert the new user into the database
		err = db.QueryRow(context.Background(), "INSERT INTO users(name, phone_number) VALUES($1, $2) RETURNING id",
			user.Name, user.PhoneNumber).Scan(&user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": user.ID})
	}
}

func generateOTP(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		// Retrieve user from the database
		err := db.QueryRow(context.Background(), "SELECT id FROM users WHERE phone_number = $1",
			user.PhoneNumber).Scan(&user.ID)
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
			return
		}

		// Generate OTP
		otp := generateRandomOTP()
		expirationTime := time.Now().Add(time.Minute)

		// Update user OTP in the database
		_, err = db.Exec(context.Background(), "UPDATE users SET otp = $1, otp_expiration_time = $2 WHERE id = $3",
			otp, expirationTime, user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"otp": otp})
	}
}

func verifyOTP(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		// Retrieve user from the database
		err := db.QueryRow(context.Background(), "SELECT otp, otp_expiration_time FROM users WHERE phone_number = $1",
			user.PhoneNumber).Scan(&user.OTP, &user.OTPExpirationTime)
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify OTP"})
			return
		}

		// Check if OTP is correct and not expired
		if user.OTP == "" || time.Now().After(user.OTPExpirationTime) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP expired or incorrect"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
	}
}

func generateRandomOTP() string {
	// Implement logic to generate random OTP (e.g., using crypto/rand)
	return "1234" // Dummy OTP for example purposes
}
