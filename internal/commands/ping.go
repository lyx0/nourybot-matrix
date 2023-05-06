package commands

import (
	"fmt"

	"github.com/lyx0/nourybot-matrix/internal/common"
)

func Ping() (string, error) {
	n := common.GetCommandsUsed()
	resp := fmt.Sprintf("Pong! Commands used: %v", n)
	return resp, nil
}
