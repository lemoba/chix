package chix

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"github.com/lemoba/chix/pkg/errors"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type Context struct {
	W http.ResponseWriter
	R *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{W: w, R: r}
}

func (c *Context) JSON(status int, v interface{}) {
	render.Status(c.R, status)
	render.JSON(c.W, c.R, v)
}

func (c *Context) Success(data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: errors.ErrSuccess.Code,
		Msg:  errors.ErrSuccess.Msg,
		Data: data,
	})
}

func (c *Context) Throw(err error) {
	panic(err)
}

func (c *Context) BindJSON(v interface{}) error {
	return json.NewDecoder(c.R.Body).Decode(v)
}

func (c *Context) Query(key string) string {
	return c.R.URL.Query().Get(key)
}

func (c *Context) PostForm(key string) string {
	_ = c.R.ParseForm()
	return c.R.FormValue(key)
}

func (c *Context) GetMethod() string {
	return c.R.Method
}

func (c *Context) GetPath() string {
	return c.R.URL.Path
}

func (c *Context) GetHeader(key string) string {
	return c.R.Header.Get(key)
}

func (c *Context) Validate(v interface{}) {
	if err := validator.New().Struct(v); err != nil {
		c.Throw(errors.ErrInvalidRequestParam)
	}
}
