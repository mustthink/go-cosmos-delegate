package main

import (
	"context"

	app "github.com/mustthink/go-cosmos-delegate/internal"
)

func main() {
	application := app.New()

	ctx := context.Background()
	application.Run(ctx)
}
