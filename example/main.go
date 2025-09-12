package main

import (
	"github.com/lemoba/chix"
	"github.com/lemoba/chix/example/internal/handler/user"
)

func main() {
	r := chix.New()

	r.Post("/name", user.SayHello)
	r.Get("/health", user.Health)

	r.Run(":3000")
}
