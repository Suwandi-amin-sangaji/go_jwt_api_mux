package authcontroller

import (
	"encoding/json"
	"go-jwt-api/config"
	"go-jwt-api/helpers"
	"go-jwt-api/models"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// Mengambil Inputan Json
	var userInput models.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helpers.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	defer r.Body.Close()

	// Mengambil Data User Berdasarkan Username
	var user models.User
	if err := models.DB.Where("username = ?", userInput.Username).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response := map[string]string{"message": "username Atau Password Salah"}
			helpers.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		default:
			response := map[string]string{"message": err.Error()}
			helpers.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
	}

	// Pengecekan Password Dengan Yang Ada di Database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		response := map[string]string{"message": "Username Atau Password Salah"}
		helpers.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	// Proses Pembuatan JWT
	expTime := time.Now().Add(time.Minute * 1)
	claims := &config.JWTClaim{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-jwt-api",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	// Mendekalasikan Algoritma yang akan di gunakan untuk login

	tokenAlgoritma := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenAlgoritma.SignedString(config.JWT_KEY)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Jika Behasil Set Token Ke Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Token",
		Path:     "/",
		Value:    token, //mengambil variabel Token Di atas
		HttpOnly: true,
	})

	response := map[string]string{"message": "Login Berhasil"}
	helpers.ResponseJSON(w, http.StatusOK, response)
}

func Register(w http.ResponseWriter, r *http.Request) {
	// Mengambil Inputan Json
	var userInput models.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helpers.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	defer r.Body.Close()

	// Hash PAssword Menggunakan bcCrypt
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	userInput.Password = string(hashPassword)

	// insert To Database
	if err := models.DB.Create(&userInput).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}
	response := map[string]string{"message": "Register Berhasil"}
	helpers.ResponseJSON(w, http.StatusOK, response)
}

func Logout(w http.ResponseWriter, r *http.Request) {

	// Hapus  Token di Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Token",
		Path:     "/",
		Value:    "", //mengambil variabel Token Di atas
		HttpOnly: true,
		MaxAge:   -1,
	})

	response := map[string]string{"message": "Logout Berhasil"}
	helpers.ResponseJSON(w, http.StatusOK, response)
}
