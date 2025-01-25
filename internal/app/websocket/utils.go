package websocket

import "github.com/sajad-dev/eda-architecture/internal/database/model"

func checkPrivateKey(secret string, public string) bool {
	val := model.Get([]string{"public_key", "secret_key"}, "channel", []model.Where_st{
		{Key: "public_key", Value: public, After: "AND", Operator: "="},
		{Key: "secret_key", Value: secret, After: "AND", Operator: "="},
	}, "id", true)

	return len(val) > 0
}
