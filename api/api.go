package api

import (
	"encoding/json"
	"fmt"

	"github.com/gothello/logic"
)

func EncodeState(b *logic.Board) []byte {
	res, err := json.Marshal(b.Board)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	return res
}
