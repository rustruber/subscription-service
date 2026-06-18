// tools.go
//go:build tools

package tools

import (
	_ "github.com/swaggo/swag/cmd/swag"
)

//go:generate swag init -g cmd/server/main.go -o ./docs