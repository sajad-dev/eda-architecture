package websocket

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/sajad-dev/eda-architecture/internal/database/model"
)

// generateHMACSHA256 generates an HMAC-SHA256 signature using a secret key
func generateHMACSHA256(secret string, message string) string {
	h := hmac.New(sha256.New, []byte(secret)) // Create a new HMAC instance
	h.Write([]byte(message)) // Write the message to the HMAC instance
	return hex.EncodeToString(h.Sum(nil)) // Return the hexadecimal representation of the signature
}

// checkPrivateKey verifies the authenticity of the request using HMAC-SHA256
func checkPrivateKey(queryParams url.Values, path string, method string) bool {
	auth_signature := queryParams.Get("auth_signature") // Extract the provided auth signature
	delete(queryParams, "auth_signature") // Remove the auth signature from query parameters
	var keys []string
	for key := range queryParams {
		keys = append(keys, key) // Collect all query parameter keys
	}
	sort.Strings(keys) // Sort the keys alphabetically

	var stringToSign strings.Builder
	stringToSign.WriteString(method + "\n" + path + "\n") // Construct the base string

	for i, key := range keys {
		if i > 0 {
			stringToSign.WriteString("&")
		}
		stringToSign.WriteString(key + "=" + queryParams.Get(key)) // Append query parameters
	}

	// Retrieve public and secret keys from the database
	ou := model.Get([]string{"public_key", "secret_key"}, "channels", []model.Where_st{
		{Key: "public_key", Value: queryParams.Get("auth_key"), After: "", Operator: "="},
	}, "id", true)
	// Generate expected signature using the secret key
	expectedSignature := generateHMACSHA256(ou[0]["secret_key"], stringToSign.String())

	// Compare expected signature with provided signature
	if expectedSignature != auth_signature {
		return false
	}
	return true
}

// removeClient removes a WebSocket connection from a list of connections
func removeClient(connList []*websocket.Conn, connDel *websocket.Conn) []*websocket.Conn {
	for index, conn := range connList {
		if conn == connDel {
			return append(connList[:index], connList[index+1:]...) // Remove the client connection
		}
	}
	return connList // Return the updated list
}
