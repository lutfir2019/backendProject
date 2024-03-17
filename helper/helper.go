package helper

import (
	crand "crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
	"go.mod/database"
	"go.mod/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func HashAndSalt(pwd []byte) string {
	hash, _ := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	return string(hash)
}

func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	return err == nil
}

// Universal date the Session Will Expire
func SessionExpires() time.Time {
	return time.Now().Add(1 * 60 * time.Minute)
}

func SetPagination(totalItems int64, limit int, page int) fiber.Map {
	totalPages := int(totalItems) / limit
	if int(totalItems)%limit > 0 {
		totalPages++
	}
	return fiber.Map{
		"totalItems":  totalItems,
		"totalPages":  totalPages,
		"currentPage": page,
		"perPage":     limit,
	}
}

func GenerateCode(inputText string) string {
	// Mengubah string menjadi slice kata-kata
	words := strings.Fields(inputText)
	// Menggabungkan kata-kata menjadi satu string tanpa spasi
	charset := strings.Join(words, "")

	// Get the initials from the input string
	for {
		// Generate a random code
		code := make([]byte, len(charset))
		for i := range code {
			code[i] = charset[rand.Intn(len(charset))]
		}

		// Convert code to string
		newCode := strings.ToUpper(string(code[:3]))

		// Check if the code already exists
		count := 0
		for {
			count++

			if count == 3 {
				newCode = fmt.Sprintf("%s%d", newCode, count-2)
				if !codeExists(newCode) {
					return newCode
				}
			}

			if !codeExists(newCode) {
				return string(newCode)
			}
		}
	}
}

func codeExists(newCode string) bool {
	found := model.Shop{}
	query := model.Shop{Spcd: newCode}
	err := database.DB.First(&found, &query).Error
	return err != gorm.ErrRecordNotFound
}

func ResponsSuccess(c *fiber.Ctx, status int, msg string, data interface{}, totalItems int64, limit int, page int) error {
	return c.Status(status).JSON(fiber.Map{
		"message": msg,
		"data": fiber.Map{
			"data":       data,
			"pagination": SetPagination(totalItems, limit, page),
		},
	})
}

func ResponsError(c *fiber.Ctx, status int, msg string, err error) error {
	return c.Status(status).JSON(fiber.Map{
		"message": msg,
		"details": err,
	})
}

func ResponseBasic(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(fiber.Map{
		"message": msg,
	})
}

func ParseSessionId(c *fiber.Ctx, tokenString string) error {
	token, err := guuid.Parse(tokenString)
	if err != nil {
		return ResponsError(c, 400, "Invalid Session Id", err)
	}

	session := model.Session{}
	query := model.Session{Sessionid: token}
	err = database.DB.First(&session, &query).Error
	if err == fiber.ErrNotFound {
		return ResponsError(c, 401, "Please login before", err)
	}

	user := model.User{}
	queryUser := model.User{UID: session.UserRefer}
	err = database.DB.First(&user, &queryUser).Error
	if err == fiber.ErrNotFound {
		return ResponsError(c, 404, "User not register", err)
	}

	userData := map[string]interface{}{
		"username":  user.Unm,
		"role":      user.Rlcd,
		"shopcode":  user.Spcd,
		"sessionId": token,
	}
	c.Locals("user", userData)

	return c.Next()
}

// GetPrivateKey mengembalikan kunci privat RSA
func GetPrivateKey() *rsa.PrivateKey {
	// Misalnya, di sini Anda menghasilkan kunci privat baru
	privateKey, err := rsa.GenerateKey(crand.Reader, 2048)
	if err != nil {
		log.Fatal("Failed to generate RSA private key:", err)
	}
	return privateKey
}

func GetUserLocal(c *fiber.Ctx, get string) (string, error) {
	userLocal := c.Locals("user")
	if userLocal != nil {
		if user, ok := userLocal.(map[string]interface{}); ok {
			if value, found := user[get]; found {
				if strValue, ok := value.(string); ok {
					return strValue, nil // Mengembalikan nilai string jika ditemukan
				}
				return "", fmt.Errorf("value is not a string") // Mengembalikan error jika nilai bukan string
			}
		}
	}
	return "", fmt.Errorf("user or requested key not found") // Mengembalikan error jika user atau key tidak ditemukan
}
