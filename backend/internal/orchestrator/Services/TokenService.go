package Services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	//"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
)

func CreateToken(user_id uuid.UUID) (string, error) {
	// type Token struct {
	// 	token      string
	// 	expires_at int
	// }
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	//claims["nbf"] = time.Now().Add(time.Minute).Unix()
	claims["iat"] = time.Now().Unix()
	//claims["exp"] = time.Now().Add(5 * time.Minute).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("API_SECRET")))

}

func TokenValid(r *http.Request) error {
	tokenString := ExtractToken(r)
	//log.Println(tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		Pretty(claims)
	}
	return nil
}

func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}

func ExtractTokenID(r *http.Request) (*uuid.UUID, error) {
	var uid string
	tokenString := ExtractToken(r)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		str, _ := claims["user_id"].(string)
		uid = str
		log.Println(uid)
		if err != nil {
			return nil, err
		}
		result, _ := uuid.FromString(uid)
		return &result, nil
	}
	return nil, nil
}

func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(b))
}
