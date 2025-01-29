package websocket

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/gorilla/websocket"
	"github.com/sajad-dev/eda-architecture/internal/database/model"
)

func generateHMACSHA256(secret, message string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func checkPrivateKey(sin string, public string) bool {
	ou := model.Get([]string{"public_key", "secret_key"}, "channels", []model.Where_st{
		{Key: "public_key", Value: public, After: "", Operator: "="},
	}, "id", true)

	return len(ou) > 0
}

func removeClient(connList []*websocket.Conn, connDel *websocket.Conn) []*websocket.Conn {
	for index, conn := range connList {
		if conn == connDel {
			return append(connList[:index], connList[index+1:]...)
		}
	}
	return connList
}
