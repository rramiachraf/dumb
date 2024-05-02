package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/rramiachraf/dumb/handlers"
	"github.com/rramiachraf/dumb/utils"
)

func main() {
	logger := utils.NewLogger(os.Stdout)

	server := &http.Server{
		Handler:      handlers.New(logger),
		WriteTimeout: 25 * time.Second,
		ReadTimeout:  25 * time.Second,
	}

	port, _ := strconv.Atoi(os.Getenv("PORT"))

	if port == 0 {
		port = 5555
		logger.Info("using default port %d", port)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Error(err.Error())
	}

	logger.Info("server is listening on port %d", port)

	if err := server.Serve(l); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
