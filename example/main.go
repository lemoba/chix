package main

import (
	"github.com/lemoba/chix"
	"github.com/lemoba/chix/example/pkg/handler/user"
)

func main() {
	c := chix.New()

	c.Post("/name", user.SayHello)
	c.Get("/health", user.Health)

	c.Run(":3000")
}
