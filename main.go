package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const assetsDir = "dist/project/browser"

//go:embed dist/project/browser/*
var assets embed.FS

type AssetsCache struct {
	ctx    context.Context
	assets map[string]*Asset
}

type Asset struct {
	Name        string
	Content     []byte
	ContentType string
	TTL         int
}

func newAssetsCache(ctx context.Context) *AssetsCache {
	a := &AssetsCache{
		ctx:    ctx,
		assets: make(map[string]*Asset),
	}

	a.checkTTL()

	return a
}
func (ac *AssetsCache) checkTTL() {
	go func() {
		for {
			select {
			case <-ac.ctx.Done():
				return
			default:
				for k, v := range ac.assets {
					if v.TTL < time.Now().Second() {
						delete(ac.assets, k)
					}
				}
				time.Sleep(1 * time.Minute)
			}
		}
	}()
}

func (ac *AssetsCache) GetAsset(name string) (*Asset, error) {
	if strings.Contains(name, "..") {
		return nil, fmt.Errorf("Invalid file path")
	}

	if strings.Contains(name, "/") {
		name = strings.TrimPrefix(name, "/")
	}

	path, ok := ac.assets[name]
	if ok {
		return path, nil
	}

	content, err := assets.ReadFile(fmt.Sprintf("%s/%s", assetsDir, name))
	if err != nil {
		return nil, err
	}

	ac.assets[name] = &Asset{
		Name:        name,
		Content:     content,
		TTL:         time.Now().Add(12 * time.Hour).Second(),
		ContentType: detectContentType(name),
	}

	return ac.assets[name], nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var cache = newAssetsCache(ctx)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fileName := r.URL.Path
		if fileName == "/" {
			fileName = "index.html"
		}

		// Get the file from the assets cache
		asset, err := cache.GetAsset(fileName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the file to the response
		w.Header().Set("Content-Type", asset.ContentType)
		w.Write(asset.Content)
	})

	http.ListenAndServe(":8080", nil)
}

func detectContentType(filename string) string {
	switch {
	case strings.HasSuffix(filename, ".html"):
		return "text/html"
	case strings.HasSuffix(filename, ".css"):
		return "text/css"
	case strings.HasSuffix(filename, ".js"):
		return "application/javascript"
	case strings.HasSuffix(filename, ".json"):
		return "application/json"
	case strings.HasSuffix(filename, ".xml"):
		return "application/xml"
	case strings.HasSuffix(filename, ".png"):
		return "image/png"
	case strings.HasSuffix(filename, ".jpg"):
		return "image/jpeg"
	case strings.HasSuffix(filename, ".gif"):
		return "image/gif"
	case strings.HasSuffix(filename, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(filename, ".ico"):
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}
}
