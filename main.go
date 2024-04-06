package main

import (
	"flag"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var port string

func init() {
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05.000"})
	flag.StringVar(&port, "port", "3030", "server port")
	flag.Parse()
}

func main() {
	http.HandleFunc("/", index)

	log.Infof("server listen on port %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		panic(err)
	}
}
