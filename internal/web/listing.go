package web

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/CorySanin/downloadcountlisting/internal/config"
)

type FileEntry struct {
	Filename string
	Size     int64
	Date     time.Time
	DL       *int
	DLTotal  *int
}

var conf config.Conf

func InitWeb(cfg config.Conf) {
	conf = cfg
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/v1/") {
		apiHandler(w, r)
		return
	}
	destination := path.Join(*conf.Directory, r.URL.Path[1:])
	if childDirs, _, err := getChildren(destination); err == nil {
		fmt.Printf("\ndirs: %s", strings.Join(childDirs, ", "))
		fmt.Fprintf(w, "found directory %s", r.URL.Path[1:])
		return
	} else if file, err := os.Open(destination); err == nil {
		defer file.Close()
		_, fileName := filepath.Split(destination)
		fileStat, err := file.Stat()
		if err != nil {
			http.Error(w, "Internal server error.", 500)
			return
		}
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Type", r.Header.Get("Content-Type")) // TODO: set content-type accordingly
		w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
		http.ServeContent(w, r, fileName, fileStat.ModTime(), file)
		return
	}
	http.Error(w, "404 file not found", 404)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "API not yet implemented", 404)
}

func getChildren(path string) ([]string, []FileEntry, error) {
	entires, err := os.ReadDir(path)
	if err != nil {
		return nil, nil, err
	}
	var dirCount = 0
	var fileCount = 0
	for _, v := range entires {
		if v.IsDir() {
			dirCount++
		} else {
			fileCount++
		}
	}

	childDirs := make([]string, dirCount)
	childFiles := make([]FileEntry, fileCount)
	dirCount = 0
	fileCount = 0

	for _, v := range entires {
		if v.IsDir() {
			childDirs[dirCount] = v.Name()
			dirCount++
		} else {
			var fEntry = FileEntry{
				Filename: v.Name(),
			}
			if info, err := v.Info(); err == nil {
				fEntry.Size = info.Size()
				fEntry.Date = info.ModTime()
			}
			childFiles[fileCount] = fEntry
			fileCount++
		}
	}
	return childDirs, childFiles, nil
}
