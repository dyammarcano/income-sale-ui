//go:build test

package assets

type Asset struct {
	Name        string
	Content     []byte
	ContentType string
	TTL         int
}

func GetAsset(name string) (*Asset, error) {
	return nil, nil
}
