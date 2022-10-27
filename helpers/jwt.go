package helpers

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var secretKey = "rangkuti"

func GenerateToken(id, age int, email, username string) string {
	claims := jwt.MapClaims{
		"id": id,
		"email": email,
		"username": username,
		"age": age,
	}

	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := parseToken.SignedString([]byte(secretKey))

	return signedToken
}

func VerifyToken(c *fiber.Ctx) (interface{}, error) {
	errReponse := errors.New("sign in to proceed")
	headerToken := c.Get("Authorization")
	bearer := strings.HasPrefix(headerToken, "Bearer")

  	if !bearer {
    	return nil, errReponse
  	}

  	stringToken := strings.Split(headerToken, " ")[1]
	
	token, _ := jwt.Parse(stringToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok{
			return nil, errReponse
		}
		return []byte(secretKey), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token. Valid {
		return nil, errReponse
	}

   	return token.Claims.(jwt.MapClaims), nil
}