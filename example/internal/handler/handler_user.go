package handler

import "github.com/lemoba/chix"

func SayHello(c *chix.Context) {
	c.Success("Hello, World!")
}
