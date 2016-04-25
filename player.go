package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Name  string
	IsDir bool
	Mode  os.FileMode
}

const (
	filePrefix = "/music/"
	root       = "./music"
	css        = "/css/"
	css_root   = "./css"
	js         = "/js/"
	js_root    = "./js"
	font       = "/font/"
	font_root  = "./font"
)

func main() {
	http.HandleFunc("/", playerMainFrame)
	http.HandleFunc(filePrefix, File)
	http.HandleFunc(css, CSS)
	http.HandleFunc(js, JS)
//	http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil)
	http.HandleFunc(font, Font)
	http.ListenAndServe(":80", nil)

}

func redir(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "https://localhost"+req.RequestURI, http.StatusMovedPermanently)
}

func playerMainFrame(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./player.html")
}

func File(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(root, r.URL.Path[len(filePrefix):])
	stat, err := os.Stat(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if stat.IsDir() {
		serveDir(w, r, path)
		return
	}
	http.ServeFile(w, r, path)
}
func CSS(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(css_root, r.URL.Path[len(css):])
	stat, err := os.Stat(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if stat.IsDir() {
		serveDir(w, r, path)
		return
	}
	http.ServeFile(w, r, path)
}
func Font(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(font_root, r.URL.Path[len(font):])
	stat, err := os.Stat(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if stat.IsDir() {
		serveDir(w, r, path)
		return
	}
	http.ServeFile(w, r, path)
}
func JS(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(js_root, r.URL.Path[len(js):])
	stat, err := os.Stat(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if stat.IsDir() {
		serveDir(w, r, path)
		return
	}
	http.ServeFile(w, r, path)
}
func serveDir(w http.ResponseWriter, r *http.Request, path string) {
	defer func() {
		if err, ok := recover().(error); ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	files, err := file.Readdir(-1)
	if err != nil {
		panic(err)
	}

	fileinfos := make([]FileInfo, len(files), len(files))

	for i := range files {
		fileinfos[i].Name = files[i].Name()
		fileinfos[i].IsDir = files[i].IsDir()
		fileinfos[i].Mode = files[i].Mode()
	}

	j := json.NewEncoder(w)

	if err := j.Encode(&fileinfos); err != nil {
		panic(err)
	}
}
