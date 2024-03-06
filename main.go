package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/rramiachraf/dumb/handlers"
	"github.com/sirupsen/logrus"
)

func main() {
	var logger = logrus.New()

	server := &http.Server{
		Handler:      handlers.New(logger),
		WriteTimeout: 25 * time.Second,
		ReadTimeout:  25 * time.Second,
	}

	port, _ := strconv.Atoi(os.Getenv("PORT"))

	if port == 0 {
		port = 5555
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatalln(err)
	}

	logger.Infof("server is listening on port %d\n", port)

	logger.Fatalln(server.Serve(l))
}
