package ui

import (
	"embed"
	"fmt"
	"strings"
)

//go:generate ng build

const assetsDir = "project/browser"

//go:embed project/browser/*
var afs embed.FS
var assets = make(map[string]*Asset)

type Asset struct {
	Name        string
	Content     []byte
	ContentType string
	TTL         int
}

func GetAsset(name string) (*Asset, error) {
	if strings.Contains(name, "/") {
		name = strings.TrimPrefix(name, "/")
	}

	path, ok := assets[name]
	if ok {
		return path, nil
	}

	content, err := afs.ReadFile(fmt.Sprintf("%s/%s", assetsDir, name))
	if err != nil {
		return nil, err
	}

	assets[name] = &Asset{
		Name:        name,
		Content:     content,
		ContentType: contentType(name),
	}

	return assets[name], nil
}

func contentType(filename string) string {
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
