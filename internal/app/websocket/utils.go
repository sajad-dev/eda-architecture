package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/sajad-dev/eda-architecture/internal/database/model"
)

func checkPrivateKey(secret string, public string) bool {
	ou := model.Get([]string{"public_key", "secret_key"}, "channels", []model.Where_st{
		{Key: "public_key", Value: public, After: "AND", Operator: "="},
		{Key: "secret_key", Value: secret, After: "", Operator: "="},
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
