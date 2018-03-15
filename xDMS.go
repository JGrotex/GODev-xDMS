package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var rootpath string
var port string
var ref string

func main() {

	//defaults
	port = "4180"         // default Listen Port
	rootpath = "c://temp" // default root Folder

	//check commandline Parameter
	for index := 1; index < len(os.Args); index++ {
		if strings.Contains(os.Args[index], "-rootpath=") {
			rootpath = strings.Replace(os.Args[index], "-rootpath=", "", 1)
		}
		if strings.Contains(os.Args[index], "-port=") {
			port = strings.Replace(os.Args[index], "-port=", "", 1)
		}
	}

	// Create new ServeMux for asset files
	assetServer := http.NewServeMux()
	assetServer.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./web"))))

	http.HandleFunc("/", serveGZIP(func(res http.ResponseWriter, req *http.Request) {

		res.Header().Set("Cache-Control", "public ,max-age=0")

		if strings.Contains(req.URL.Path, ".") {
			assetServer.ServeHTTP(res, req)

		} else if req.URL.Path == "/" {
			ref = req.URL.Query().Get("ref")

			//some security checks to block some requests of folder browsing
			ref = strings.Replace(ref, ".", "unknown", -1)
			if ref == "" {
				ref = "unknown"
			}
			// ***

			readCopyServe("./web/xDMS.html", res, req)

		} else if req.URL.Path == "/upload" {
			readCopyServe("./web/xDMSupload.html", res, req)
		} else {
			errorCodeHandler(res, req, http.StatusNotFound)
		}

	}))

	http.HandleFunc("/list/", serveGZIP(func(res http.ResponseWriter, req *http.Request) {
		Dir := strings.Replace(req.URL.Query().Get("dir"), "~!", " ", -1)

		Dir = rootpath + "/" + ref + "/" + Dir
		//fmt.Println("list Dir: " + Dir)

		if exists(Dir) {
			res.Write(getDirInfoList(Dir))
		} else {
			res.Write([]byte("Invalid Path"))
		}

	}))

	http.HandleFunc("/file/", serveGZIP(func(res http.ResponseWriter, req *http.Request) {
		Dir := strings.Replace(req.URL.Query().Get("dir"), "~!", " ", -1)
		Dir = rootpath + "/" + ref + "/" + Dir
		//fmt.Println("file Dir: " + Dir)

		if Dir != "" && existsAS(Dir) == "file" {
			file, err := os.Open(Dir)
			if err != nil {
				errorCodeHandler(res, req, http.StatusNotFound)
			}
			io.Copy(res, file)
			file.Close()
		} else {
			res.Write([]byte("Invalid :  " + existsAS(Dir)))
		}

	}))

	http.HandleFunc("/fupload", func(res http.ResponseWriter, req *http.Request) {
		Dir := strings.Replace(req.URL.Query().Get("dir"), "~!", " ", -1)
		Dir = rootpath + "/" + ref + "/" + Dir
		//fmt.Println("file Dir: " + Dir)

		req.ParseMultipartForm(32 << 20)
		file, handler, err := req.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		readCopyServe("./web/xDMS.html", res, req)

		filepath := strings.Split(handler.Filename, "\\") //Edge fix, as it provide full path
		filename := filepath[len(filepath)-1]

		f, err := os.OpenFile(Dir+filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	})

	fmt.Println("xDMS is running on Port: " + port + ", and rootpath is: " + rootpath)

	// Listen on Port
	errHTTP := http.ListenAndServe(":"+port, nil)
	check(errHTTP)
}

func readCopyServe(filename string, res http.ResponseWriter, req *http.Request) {
	file, err := os.Open(filename)
	if err != nil {
		errorCodeHandler(res, req, http.StatusNotFound)
	}
	io.Copy(res, file)
	file.Close()
}

func getLocalpath() string {
	localpath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return localpath
}

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func existsAS(name string) string {
	Stat, err := os.Stat(name)
	if err == nil && !os.IsNotExist(err) {
		if Stat.Mode().IsDir() {
			return "directory"
		} else if Stat.Mode().IsRegular() {
			return "file"
		}
	}
	return "invalid"
}

// FileInfo Data
type FileInfo struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
}

func getDirInfoList(Directory string) []byte {
	dir, err := os.Open(Directory)
	check(err)
	entries, err := dir.Readdir(0)
	check(err)

	list := []FileInfo{}

	for _, entry := range entries {
		f := FileInfo{
			Name:    entry.Name(),
			Size:    entry.Size(),
			Mode:    entry.Mode(),
			ModTime: entry.ModTime(),
			IsDir:   entry.IsDir(),
		}
		list = append(list, f)
	}

	output, err := json.Marshal(list)
	check(err)

	DirInfoList := output
	return DirInfoList
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	if "" == w.Header().Get("Content-Type") {
		// no content type, apply un-gzipped body
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}
	return w.Writer.Write(b)
}

func serveGZIP(fn http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			fn(res, req)
			return
		}
		res.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(res)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: res}
		fn(gzr, req)
	})
}

func errorCodeHandler(res http.ResponseWriter, r *http.Request, status int) {
	res.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(res, "Error 404, Page not available")
	}
}
