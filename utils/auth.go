package utils

import (
	crand "crypto/rand"
	"fmt"
	"io"
	mrand "math/rand"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	// CookieName is the name of the jwt session cookie
	CookieName = "session_jwt"
	// characters for random password generator
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

// User holds the user information
type User struct {
	Name     string
	Password []byte
}

// GetUser will return the user data
func GetUser() (user User, err error) {
	err = Storm.Get("auth", "user", &user)
	return
}

// InitUser will set the user data
func InitUser(name string) (err error) {
	password, hash, err := RandomPassword()
	if err != nil {
		return
	}

	// Print out the info to the console
	fmt.Printf("User Generated\nName: %s\nPassword: %s\n", name, password)

	return Storm.Set("auth", "user", &User{
		Name:     name,
		Password: hash,
	})
}

// CreateCookie will make a cookie for the JWT
func CreateCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Expires:  time.Now().Add(90 * 24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	}
}

// DeleteCookie will delete the JWT cookie
func DeleteCookie() *http.Cookie {
	return &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Expires:  time.Now().AddDate(-1, 0, 0),
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	}
}

// NewSecret will generate a 256 bit HMAC key
func NewSecret() []byte {
	key := make([]byte, 32)
	_, err := io.ReadFull(crand.Reader, key[:])
	if err != nil {
		panic(err)
	}
	return key
}

// GetSecret will return the secret key for JWT auth
func GetSecret() (secret []byte, err error) {
	err = Storm.Get("auth", "secret", &secret)
	return
}

// InitSecret will set the JWT key
func InitSecret() (err error) {
	fmt.Println("Secret Generated")
	return Storm.Set("auth", "secret", NewSecret())
}

// RandomPassword will generate a random password for password resets
func RandomPassword() (password string, hash []byte, err error) {
	password = generateRandomPassword(20)
	hash, err = HashPassword(password)
	return
}

// will generate a password with random characters
func generateRandomPassword(n int) string {
	// random source
	src := mrand.NewSource(time.Now().UnixNano())

	// byte slice to hold password
	b := make([]byte, n)

	// range over byte slice and fill with random letters
	for i := range b {
		b[i] = letterBytes[src.Int63()%int64(len(letterBytes))]
	}

	return string(b)
}

// HashPassword generates a bcrypt hash of the password using work factor 14.
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// CheckPassword securely compares a bcrypt hashed password with its possible
// plaintext equivalent. Returns nil on success, or an error on failure.
func CheckPassword(password string) (err error) {
	// get the saved user password hash
	user, err := GetUser()
	if err != nil {
		return
	}

	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}
