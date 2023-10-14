package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eldius/file-server/internal/config"
	"github.com/eldius/file-server/internal/logger"
	"github.com/gabriel-vasile/mimetype"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func HandleFileRequest(basePath string) func(w http.ResponseWriter, r *http.Request) {
	if home, err := filepath.Abs(basePath); err != nil {
		err = fmt.Errorf("parsing base path: %w", err)
		logger.GetLogger().With("error", err).Error("failed to parse base path")
		panic(err)
	} else {
		return func(w http.ResponseWriter, r *http.Request) {

			diskPath := filepath.Join(home, r.URL.Path)
			log := logger.GetLogger().With("req_path", r.URL.Path, "path_on_disk", diskPath)

			log.Debug("Request Received")

			if s, err := os.Stat(diskPath); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					w.WriteHeader(http.StatusNotFound)
					_ = json.NewEncoder(w).Encode(map[string]any{
						"error": err,
						"path":  r.URL.Path,
					})
					return
				} else if errors.Is(err, os.ErrInvalid) {
					w.WriteHeader(http.StatusBadRequest)
					_ = json.NewEncoder(w).Encode(map[string]any{
						"error": err,
						"path":  r.URL.Path,
					})
					return
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(map[string]any{
						"error": err,
						"path":  r.URL.Path,
					})
					return
				}
			} else if s.IsDir() { //Handle dir
				DirectoryHandler(w, diskPath, r.URL.Path)
				return
			} else { //Handle file
				if mt, err := mimetype.DetectFile(diskPath); err != nil {
					err = fmt.Errorf("detecting file mime type: %w", err)
					log.With("error", err).Error("failed to detect mime type")
				} else {
					log.With("mime-type", mt.String()).Debug("mime type detected")
					w.Header().Add("content-type", mt.String())
				}
				if f, err := os.Open(diskPath); err != nil {
					err = fmt.Errorf("opening file: %w", err)
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(err.Error()))
					return
				} else {
					w.WriteHeader(http.StatusOK)
					_, _ = io.Copy(w, f)
					return
				}
			}

			if config.GetDebugModeEnabled() {
				log.Debug("RequestReceived")
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		}
	}
}

func Start(port int, basePath string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HandleFileRequest(basePath))
	log := logger.GetLogger()

	log.With("debug", config.GetDebugModeEnabled(), "port", config.GetServerPort()).Info("Starting server")

	s := http.Server{
		Addr:                         fmt.Sprintf(":%d", port),
		Handler:                      mux,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
	}

	fmt.Printf("server will be listening on port %d", port)

	if err := s.ListenAndServe(); err != nil {
		err = fmt.Errorf("starting to listen: %w", err)
		return err
	}
	return nil
}
