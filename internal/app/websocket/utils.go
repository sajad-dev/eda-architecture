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

func generateHMACSHA256(secret, message string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func checkPrivateKey(queryParams url.Values, path string, method string) bool {
	auth_signature := queryParams.Get("auth_signature")
	delete(queryParams, "auth_signature")

	var keys []string
	for key := range queryParams {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var stringToSign strings.Builder
	stringToSign.WriteString(method + "\n" + path + "\n")

	for i, key := range keys {
		if i > 0 {
			stringToSign.WriteString("&")
		}
		stringToSign.WriteString(key + "=" + queryParams.Get(key))
	}

	ou := model.Get([]string{"public_key", "secret_key"}, "channels", []model.Where_st{
		{Key: "public_key", Value: queryParams.Get("auth_key"), After: "", Operator: "="},
	}, "id", true)

	expectedSignature := generateHMACSHA256(ou[0]["secret_key"], stringToSign.String())

	if expectedSignature != auth_signature {
		return false
	}
	return true
}

func removeClient(connList []*websocket.Conn, connDel *websocket.Conn) []*websocket.Conn {
	for index, conn := range connList {
		if conn == connDel {
			return append(connList[:index], connList[index+1:]...)
		}
	}
	return connList
}
