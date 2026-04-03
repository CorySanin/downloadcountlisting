package web

import (
	"fmt"
	"net/http"
	"path"

	"github.com/CorySanin/downloadcountlisting/internal/config"
)

var conf config.Conf

func InitWeb(cfg config.Conf) {
	conf = cfg
}

func Handler(w http.ResponseWriter, r *http.Request) {
	destination := path.Join(*conf.Directory, r.URL.Path[1:])

	fmt.Fprintf(w, "Looking for %s", destination)
}
