package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CorySanin/downloadcountlisting/internal/config"
	"github.com/CorySanin/downloadcountlisting/internal/web"
)

func main() {
	conf := config.Config()
	web.InitWeb(conf)
	http.HandleFunc("/", web.Handler)
	fmt.Printf("Listening on port %d", *conf.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *conf.Port), nil))
}
