package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates JWT and sets user context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid auth header format"})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				return nil, errors.New("missing JWT_SECRET environment variable")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Safe extraction helpers
		getInt := func(key string) int {
			if val, ok := claims[key]; ok {
				switch v := val.(type) {
				case float64:
					return int(v)
				case int:
					return v
				}
			}
			return 0
		}

		getString := func(key string) string {
			if val, ok := claims[key]; ok {
				if s, ok := val.(string); ok {
					return s
				}
			}
			return ""
		}

		userID := getInt("user_id")
		employeeID := getInt("employee_id")
		role := getString("role")
		exp := int64(getInt("exp"))

		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}

		if exp != 0 && time.Now().Unix() > exp {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		// Store for next handlers
		c.Set("user_id", userID)
		c.Set("employee_id", employeeID)
		c.Set("role", role)

		c.Next()
	}
}

// RoleRequired ensures only specified roles can access an endpoint
func RoleRequired(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No role found in context"})
			c.Abort()
			return
		}

		role, ok := roleVal.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid role type"})
			c.Abort()
			return
		}

		for _, allowed := range allowedRoles {
			if strings.EqualFold(role, allowed) {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Length, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
