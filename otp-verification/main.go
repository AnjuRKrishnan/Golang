package main

import (
	"context"
	"fmt"
	"os"
	"net/http"
	"time"
	"math/rand"


	"github.com/jackc/pgx/v4"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	PhoneNumber       string    `json:"phone_number"`
	OTP               string    `json:"otp"`
	OTPExpirationTime time.Time `json:"otp_expiration_time"`
}

var conn *pgx.Conn

func main() {

	router := gin.Default()
	var err error
	
	dbPool, err := pgxpool.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close()

// Routes
router.POST("/api/users", createUser(dbPool))
router.POST("/api/users/generateotp", generateOTP(dbPool))
router.POST("/api/users/verifyotp", verifyOTP(dbPool))

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
	fmt.Println(user.PhoneNumber);
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
	expirationTime := time.Now().UTC().Add(time.Minute)

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
	var existingUser User
	// Retrieve user from the database
	err := db.QueryRow(context.Background(), "SELECT otp, otp_expiration_time FROM users WHERE phone_number = $1",
		user.PhoneNumber).Scan(&existingUser.OTP, &existingUser.OTPExpirationTime)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify OTP"})
		return
	}

	// Check if OTP is correct and not expired
	if time.Now().Before(existingUser.OTPExpirationTime) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP expired"})
		return
	}
	
	if user.OTP != existingUser.OTP {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP incorrect"})
		return
	}
	fmt.Println(existingUser.OTPExpirationTime)
	fmt.Println(time.Now().UTC())
	fmt.Println(time.Now().UTC().After(existingUser.OTPExpirationTime))
	
	

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}
}

func generateRandomOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := fmt.Sprintf("%04d", rand.Intn(10000))
	return otp
}
