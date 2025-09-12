package chix

import (
	context2 "context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/lemoba/chix/internal/middlewares"
)

type Chix struct {
	router         *Router
	Logger         Logger
	JSONSerializer JSONSerializer
}

type HTTPError struct {
	Internal error  `json:"-"`
	Message  string `json:"message"`
	Code     int    `json:"-"`
}

type JSONSerializer interface {
	Serializer(c Context, i interface{}, indent string) error
	Deserialize(c Context, i interface{}) error
}

type HandlerFunc func(Context)

type Map map[string]interface{}

// Headers
const (
	HeaderContentType   = "Content-Type"
	HeaderXForwardedFor = "X-Forwarded-For"
	HeaderXRealIP       = "X-Real-Ip"
)

// MIME types
const (
	MIMEApplicationJSON      = "application/json"
	MIMEMultipartForm        = "multipart/form-data"
	MIMETextPlain            = "text/plain"
	MIMETextPlainCharsetUTF8 = MIMETextPlain + "; " + charsetUTF8
)

const (
	charsetUTF8 = "charset=UTF-8"
)

func NewHTTPError(code int, message ...string) *HTTPError {
	he := &HTTPError{Code: code, Message: http.StatusText(code)}
	if len(message) > 0 {
		he.Message = message[0]
	}
	return he
}

func (he *HTTPError) Error() string {
	if he.Internal != nil {
		return fmt.Sprintf("code=%d, message=%v", he.Code, he.Message)
	}
	return fmt.Sprintf("code=%d, message=%v, internal=%v", he.Code, he.Message, he.Internal)
}

func (he *HTTPError) SetInternal(err error) *HTTPError {
	he.Internal = err
	return he
}

func (he *HTTPError) WithInternal(err error) *HTTPError {
	return &HTTPError{
		Internal: err,
		Code:     he.Code,
		Message:  he.Message,
	}
}

func New() *Chix {
	c := &Chix{
		router:         NewRouter(),
		Logger:         defaultLogger,
		JSONSerializer: DefaultJSONSerializer{},
	}

	c.router.Use(middlewares.RequestID)

	return c
}

func (c *Chix) Run(addr string) error {
	c.Logger.Info("Starting server on " + addr)

	server := &http.Server{
		Addr:    addr,
		Handler: c.router,
	}

	serverCtx, serverStopCtx := context2.WithCancel(context2.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		c.Logger.Info("Shutdown signal received")
		shutdownCtx, cancel := context2.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context2.DeadlineExceeded {
				c.Logger.Error("Graceful shutdown timed out, forcing exit")
				os.Exit(1)
			}
		}()

		if err := server.Shutdown(shutdownCtx); err != nil {
			c.Logger.Error("Server shutdown failed", err)
		} else {
			c.Logger.Info("Server shutdown completed")
		}

		serverStopCtx()
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		c.Logger.Error("Server startup failed", err)
		return err
	}

	<-serverCtx.Done()
	c.Logger.Info("Server exiting")
	return nil
}

func (c *Chix) Router() *Router {
	return c.router
}

func (c *Chix) Group(pattern string, fn func(g *Chix)) {
	c.router.Route(pattern, func(r chi.Router) {
		g := &Chix{
			router:         NewRouterWith(r),
			Logger:         c.Logger,
			JSONSerializer: c.JSONSerializer,
		}
		fn(g)
	})
}

func (c *Chix) SetJSONSerializer(serializer JSONSerializer) {
	if serializer != nil {
		c.JSONSerializer = serializer
	}
}

func (c *Chix) SetLogger(logger Logger) {
	if logger != nil {
		c.Logger = logger
	}
}

func (c *Chix) Connect(pattern string, handlerFunc func(Context)) {
	c.router.Connect(pattern, c.wrapHandler(handlerFunc))
}

func (c *Chix) Delete(pattern string, handlerFunc func(Context)) {
	c.router.Delete(pattern, c.wrapHandler(handlerFunc))
}

func (c *Chix) Get(pattern string, handlerFunc func(Context)) {
	c.router.Get(pattern, c.wrapHandler(handlerFunc))
}

func (c *Chix) Head(pattern string, handlerFunc func(Context)) {
	c.router.Head(pattern, c.wrapHandler(handlerFunc))
}

func (c *Chix) Options(pattern string, handlerFunc func(Context)) {
	c.router.Options(pattern, c.wrapHandler(handlerFunc))
}

func (c *Chix) Patch(pattern string, handlerFunc func(Context)) {
	c.router.Patch(pattern, c.wrapHandler(handlerFunc))
}

func (c *Chix) Post(pattern string, handlerFunc func(Context)) {
	c.router.Post(pattern, c.wrapHandler(handlerFunc))
}

func (c *Chix) Put(pattern string, handlerFunc func(Context)) {
	c.router.Put(pattern, c.wrapHandler(handlerFunc))
}

func (c *Chix) Trace(pattern string, handlerFunc func(Context)) {
	c.router.Trace(pattern, c.wrapHandler(handlerFunc))
}
