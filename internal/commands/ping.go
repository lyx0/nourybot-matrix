package commands

import (
	"fmt"

	"github.com/lyx0/nourybot-matrix/internal/common"
)

func Ping() (string, error) {
	n := common.GetCommandsUsed()
	up := common.GetUptime()
	resp := fmt.Sprintf("Pong! Commands used: %v Last restart: %v", n, up)
	return resp, nil
}
