package user

import (
	"net/http"

	"github.com/lemoba/chix"
)

type Request struct {
	Name string `json:"name"`
	Age  uint   `json:"age"`
}

func Health(c chix.Context) {
	c.Success("ok")
}

func SayHello(c chix.Context) {
	var req Request

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
}
