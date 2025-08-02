package utils

import(
	"time"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"crypto/rand"
	"encoding/hex"

)
var jwtSecret=[]byte("abcd")
func HashPassword(password string)(string,error){
	bytes,err:= bcrypt.GenerateFromPassword([] byte(password),14)
	return string(bytes),err
}

func CheckPasswordHash(password,hash string)bool{
	err:=bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
	return err==nil
}

func GenerateJWT(userID uint, roles [] string) (string, error){
	token:= jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"userId":userID,
		"roles":roles,
		"exp":time.Now().Add(time.Hour*24).Unix(),
	})
	return token.SignedString(jwtSecret)
}

func ParseJWT(tokenString string) (*jwt.Token, error) {
   return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	  if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		 return nil, jwt.ErrSignatureInvalid
	  }
	  return jwtSecret, nil
   })
}

func GenerateVerificationToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
