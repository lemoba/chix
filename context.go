package chix

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"

	chixMw "github.com/lemoba/chix/internal/middlewares"
	"github.com/lemoba/chix/pkg/errors"
)

type Response struct {
	RequestID string      `json:"request_id"`
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data,omitempty"`
}

// Context represents the context of the current HTTP request.
type Context interface {
	Request() *http.Request

	Response() http.ResponseWriter

	RealIP() string

	Path() string

	SetPath(p string)

	QueryParam(name string) string

	QueryParams() url.Values

	QueryString() string

	FormValue(name string) string

	FormParams() (url.Values, error)

	FormFile(name string) (*multipart.FileHeader, error)

	MultipartForm() (*multipart.Form, error)

	Get(key string) interface{}

	Set(key string, val interface{})

	BindJSON(i interface{}) error

	Blob(code int, contentType string, b []byte) error

	String(code int, s string) error

	JSON(code int, i interface{}) error

	JSONBlob(code int, b []byte) error

	Stream(code int, contentType string, r io.Reader) error

	Success(data interface{})

	Handler() HandlerFunc

	Logger() Logger

	SetLogger(l Logger)

	GetRequestID() string
}

type context struct {
	logger   Logger
	response http.ResponseWriter
	request  *http.Request
	query    url.Values
	path     string

	lock  sync.RWMutex
	store Map

	chix *Chix

	handler HandlerFunc
}

var _ Context = (*context)(nil)

const (
	defaultMemory = 32 << 20 // 32MB
	defaultIndent = "  "
)

func NewContext(w http.ResponseWriter, r *http.Request, chix *Chix) *context {
	return &context{
		response: w,
		request:  r,
		chix:     chix,
		logger:   defaultLogger,
		store:    make(Map),
	}
}

func (c *context) writeContextType(val string) {
	header := c.response.Header()
	if header.Get(HeaderContentType) == "" {
		header.Set(HeaderContentType, val)
	}
}

func (c *context) Request() *http.Request {
	return c.request
}

func (c *context) Response() http.ResponseWriter {
	return c.response
}

func (c *context) RealIP() string {
	if ip := c.request.Header.Get(HeaderXForwardedFor); ip != "" {
		ips := strings.Split(ip, ",")

		if len(ips) > 0 {
			xffip := strings.TrimSpace(ips[0])
			xffip = strings.TrimPrefix(xffip, "[")
			xffip = strings.TrimSuffix(xffip, "]")
			return xffip
		}

		return ip
	}

	if ip := c.request.Header.Get(HeaderXRealIP); ip != "" {
		ip = strings.TrimPrefix(ip, "[")
		ip = strings.TrimSuffix(ip, "]")
		return ip
	}

	host, _, err := net.SplitHostPort(c.request.RemoteAddr)

	if err == nil {
		return host
	}

	return c.request.RemoteAddr
}

func (c *context) Path() string {
	return c.path
}

func (c *context) SetPath(p string) {
	c.path = p
}

func (c *context) QueryParam(name string) string {
	if c.query == nil {
		c.query = c.request.URL.Query()
	}
	return c.query.Get(name)
}

func (c *context) QueryParams() url.Values {
	if c.query == nil {
		c.query = c.request.URL.Query()
	}
	return c.query
}

func (c *context) QueryString() string {
	return c.request.URL.RawQuery
}

func (c *context) FormValue(name string) string {
	return c.request.FormValue(name)
}

func (c *context) FormParams() (url.Values, error) {
	ct := c.request.Header.Get(HeaderContentType)

	var err error
	if strings.HasPrefix(ct, MIMEMultipartForm) || strings.Contains(ct, MIMEMultipartForm) {
		err = c.request.ParseMultipartForm(defaultMemory)
	} else {
		err = c.request.ParseForm()
	}

	if err != nil {
		return nil, err
	}
	return c.request.Form, nil
}

func (c *context) FormFile(name string) (*multipart.FileHeader, error) {
	f, fh, err := c.request.FormFile(name)
	if err != nil {
		return nil, err
	}

	f.Close()

	return fh, nil
}

func (c *context) MultipartForm() (*multipart.Form, error) {
	err := c.request.ParseMultipartForm(defaultMemory)
	return c.request.MultipartForm, err
}

func (c *context) Get(key string) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.store[key]
}

func (c *context) Set(key string, val interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.store == nil {
		c.store = make(Map)
	}

	c.store[key] = val
}

func (c *context) BindJSON(i interface{}) error {
	decoder := json.NewDecoder(c.request.Body)
	if err := decoder.Decode(i); err != nil {
		return err
	}
	return nil
}

func (c *context) Blob(code int, contentType string, b []byte) (err error) {
	c.writeContextType(contentType)
	c.response.WriteHeader(code)

	_, err = c.response.Write(b)

	return
}

func (c *context) String(code int, s string) (err error) {
	return c.Blob(code, MIMETextPlain, []byte(s))
}

func (c *context) json(code int, i interface{}, indent string) (err error) {
	c.writeContextType(MIMEApplicationJSON)
	c.response.WriteHeader(code)
	return c.chix.JSONSerializer.Serializer(c, i, indent)
}

func (c *context) JSON(code int, i interface{}) error {
	return c.json(code, i, defaultIndent)
}

func (c *context) JSONBlob(code int, b []byte) error {
	return c.Blob(code, MIMEApplicationJSON, b)
}

func (c *context) Stream(code int, contentType string, r io.Reader) (err error) {
	c.writeContextType(contentType)
	c.response.WriteHeader(code)
	_, err = io.Copy(c.response, r)
	return
}

func (c *context) Handler() HandlerFunc {
	return c.handler
}

func (c *context) Logger() Logger {
	if c.logger != nil {
		return c.logger
	}

	return defaultLogger
}

func (c *context) SetLogger(l Logger) {
	c.logger = l
}

func (c *context) GetRequestID() string {
	if reqID, ok := c.request.Context().Value(chixMw.RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}

func (c *context) Success(data interface{}) {
	c.JSON(http.StatusOK, Response{
		RequestID: c.GetRequestID(),
		Code:      errors.ErrSuccess.Code,
		Msg:       errors.ErrSuccess.Msg,
		Data:      data,
	})
}

func (c *context) Failure(err *errors.Error, data ...interface{}) {
	resp := Response{
		RequestID: c.GetRequestID(),
		Code:      err.Code,
		Msg:       err.Msg,
	}

	if len(data) > 0 {
		resp.Data = data[0]
	}

	c.JSON(http.StatusOK, resp)
}

func (c *context) Message(msg ...string) {
	resp := Response{
		RequestID: c.GetRequestID(),
		Code:      errors.ErrSuccess.Code,
		Msg:       errors.ErrSuccess.Msg,
	}

	if len(msg) > 0 {
		resp.Msg = msg[0]
	}

	c.JSON(http.StatusOK, resp)
}

func (c *context) Throw(err error) error {
	panic(err)
}

func (c *context) Validate(v interface{}) {
	if err := validator.New().Struct(v); err != nil {
		c.Throw(errors.ErrInvalidRequestParam)
	}
}
