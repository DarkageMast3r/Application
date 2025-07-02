package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthenticateKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "API key required",
				"message": "Please provide X-API-Key header",
			})
			c.Abort()
			return
		}

		// Pak API keys van environment variabele
		validKeys := getValidKeys()

		if !isValidKey(apiKey, validKeys) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid API key",
				"message": "The provided API key is not valid",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func getValidKeys() []string {
	keysEnv := os.Getenv("API_KEYS")
	if keysEnv == "" {
		// Testkey voor de handigheid
		return []string{"dev-key-123"}
	}

	keys := strings.Split(keysEnv, ",")
	for i, key := range keys {
		keys[i] = strings.TrimSpace(key)
	}

	return keys
}

func isValidKey(providedKey string, validKeys []string) bool {
	for _, validKey := range validKeys {
		if providedKey == validKey {
			return true
		}
	}
	return false
}

func PrintKeys() {
	apiKeys := os.Getenv("API_KEYS")
	if apiKeys == "" {
		fmt.Println("=== API Auth ===")
		fmt.Println("Development key: dev-key-123")
		fmt.Println("Zet API_KEYS in je environment voor betere security")
		fmt.Println("================")
		fmt.Println()
	} else {
		fmt.Println("=== API Auth ===")
		fmt.Println("Keys gebruiken van API_KEYS")
		fmt.Println("================")
		fmt.Println()
	}
}
