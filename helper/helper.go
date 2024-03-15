package helper

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

func ParseJwtToken(c *fiber.Ctx, token *jwt.Token) error {
	// Simpan informasi pengguna ke dalam fiber.Ctx.Locals()
	claims := token.Claims.(jwt.MapClaims)
	// username := claims["username"].(string)
	// role := claims["role"].(string)
	// shopCode := claims["spcd"].(string)

	fmt.Println(claims)
	c.Locals("user", claims)

	return c.Next()
}
