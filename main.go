package main

import (
	"github.com/dyammarcano/income-sale-ui/internal/assets"
	"log/slog"
	"net/http"
	"os"
)

//go:generate ng build

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fileName := r.URL.Path
		if fileName == "/" {
			fileName = "index.html"
		}

		// Get the file from the assets cache
		asset, err := assets.GetAsset(fileName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the file to the response
		w.Header().Set("Content-Type", asset.ContentType)
		w.Write(asset.Content)
	})

	slog.Info("Server started on port 8080")

	http.ListenAndServe(":8080", nil)
}
