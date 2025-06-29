package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	// "github.com/ggarrido85/api-backend/usecases"
	// "github.com/ggarrido85/api-backend/utils"
)

type Option func(*options)

func WithLocalTest(localTest bool) Option {
	return func(o *options) {
		o.localTest = localTest
	}
}

type options struct {
	localTest bool
}

func applyOptions(opts []Option) *options {
	o := &options{
		localTest: false,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func NewServer(
	router *gin.Engine,
	conf Configuration,
	tokenHandler string /*TokenHandler*/,
	logger *slog.Logger,
	opts ...Option,
) *http.Server {
	o := applyOptions(opts)

	addRoutes(router, conf, tokenHandler, logger)

	var host string
	if o.localTest {
		host = "localhost"
	} else {
		host = "0.0.0.0"
	}

	// Add 5 seconds to the server timeout to gracefully handle the timeout in our code
	maxTimeout := conf.DefaultTimeout + 5*time.Second

	return &http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, conf.Port),
		WriteTimeout: maxTimeout,
		ReadTimeout:  maxTimeout,
		IdleTimeout:  maxTimeout,
		Handler:      h2c.NewHandler(router, &http2.Server{}),
	}
}
