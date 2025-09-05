package user

import (
	"github.com/lemoba/chix"
	"github.com/lemoba/chix/example/response"
)

func SayHello(c *chix.Context) {
	c.Throw(response.ErrRecordsNotExist)
	c.Success("Hello, World!")
}
