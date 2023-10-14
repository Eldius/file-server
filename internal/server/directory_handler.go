package server

import (
	"embed"
	"github.com/eldius/file-server/internal/logger"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

var (
	//go:embed all:templates
	templatesDir embed.FS
	tmpl         *template.Template
)

type DirInfo struct {
	ParentPath  string
	CurrentPath string
	Entries     []Entry
}

type Entry struct {
	Name string
	Path string
}

func init() {
	tmpl = template.Must(template.ParseFS(templatesDir, "templates/*"))
}

func DirectoryHandler(w http.ResponseWriter, diskPath, reqPath string) {

	log := logger.GetLogger().With("disk_path", diskPath, "req_path", reqPath)

	var entries []Entry

	w.Header().Add("content-type", "text/html")

	if dirEntries, err := os.ReadDir(diskPath); err != nil {
		log.With("error", err).Error("failed to read entries")
	} else {
		for _, e := range dirEntries {
			entries = append(entries, Entry{
				Name: e.Name(),
				Path: filepath.Join(reqPath, e.Name()),
			})
		}
	}

	parentPath, _ := filepath.Abs(reqPath + "/..")

	resp := DirInfo{
		ParentPath:  parentPath,
		CurrentPath: reqPath,
		Entries:     entries,
	}

	log.With("resp", resp).Debug("Finished Handling")

	//if err := tmpl.Lookup("templates/dir_template.html").Execute(w, &resp); err != nil {
	if err := tmpl.ExecuteTemplate(w, "dir_template.html", &resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
}
