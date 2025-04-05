package main

import (
	"embed"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/rramiachraf/dumb/handlers"
	"github.com/rramiachraf/dumb/utils"
)

//go:embed static
var staticFiles embed.FS

func main() {
	logger := utils.NewLogger(os.Stdout)

	server := &http.Server{
		Handler:      handlers.New(logger, staticFiles),
		WriteTimeout: 25 * time.Second,
		ReadTimeout:  25 * time.Second,
	}

	PROXY_ENV := os.Getenv("PROXY")
	if PROXY_ENV != "" {
		if _, err := url.ParseRequestURI(PROXY_ENV); err != nil {
			logger.Fatal("invalid proxy")
		}

		logger.Info("using a custom proxy for requests")
	}

	port, _ := strconv.Atoi(os.Getenv("PORT"))

	if port == 0 {
		port = 5555
		logger.Infof("using default port %d", port)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatal(err.Error())
	}

	logger.Infof("server is listening on port %d", port)

	if err := server.Serve(l); err != nil {
		logger.Fatal(err.Error())
	}
}
