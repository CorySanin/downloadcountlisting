package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CorySanin/downloadcountlisting/internal/config"
	"github.com/CorySanin/downloadcountlisting/internal/web"
	"github.com/CorySanin/downloadcountlisting/pkg/storage"
)

func main() {
	conf := config.Config()
	storage := storage.New(*conf.Storage)
	web.InitWeb(conf, storage)
	http.HandleFunc("/", web.Handler)
	fmt.Printf("Listening on port %d", *conf.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *conf.Port), nil))
}
