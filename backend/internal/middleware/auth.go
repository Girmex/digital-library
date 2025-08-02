package middleware
import (
	"net/http"
	"strings"
	"digital-library/backend/internal/utils"
	"github.com/gin-gonic/gin"
	 jwt "github.com/golang-jwt/jwt/v5"
)
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context){
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		tokenString:= strings.Split(authHeader, "Bearer ")[1]

		if tokenString ==""{
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}
		token, err := utils.ParseJWT(tokenString)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		claims:=token.Claims.(jwt.MapClaims)
		c.Set("userID", claims["user_id"])
		c.Set("roles", claims["roles"])
		c.Next()



	

	}
}

func HasPermission(permission string) gin.HandlerFunc{
	return  func(c *gin.Context){
		roles, exists := c.Get("roles")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No roles found"})
			return
		}
		// Check if the user's roles have this permission
		for _, role := range roles.([]string) {
			if role == permission {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
	}
}