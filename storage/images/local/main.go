package local

import (
	"fmt"
	"image"
	"image/png"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/dairycart/dairycart/storage/images"

	"github.com/go-chi/chi"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	portKey                     = "image_storage.port"
	baseURLKey                  = "image_storage.base_url"
	storageDirKey               = "image_storage.storage_dir"
	routePrefixKey              = "image_storage.route_prefix"
	LocalProductImagesDirectory = "product_images"
)

type localImageStorer struct {
	BaseURL     string
	StorageDir  string
	RoutePrefix string
}

var _ images.ImageStorer = (*localImageStorer)(nil)

func NewLocalImageStorer() *localImageStorer {
	return &localImageStorer{
		BaseURL:     "http://localhost:4321",
		StorageDir:  LocalProductImagesDirectory,
		RoutePrefix: LocalProductImagesDirectory,
	}
}

func saveImage(in image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "error creating local file")
	}
	return png.Encode(f, in)
}

func ensureDomainHasNoPort(domain string) (string, error) {
	u, err := url.Parse(domain)
	if err != nil {
		return "", err
	}
	u.Host, _, _ = net.SplitHostPort(u.Host)

	return u.String(), nil
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	// path := fmt.Sprintf("/%s/", local.LocalProductImagesDirectory)
	// root := http.Dir(local.LocalProductImagesDirectory)

	if strings.ContainsAny(path, "{}*") {
		panic("fileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func (l *localImageStorer) Init(cfg *viper.Viper, router chi.Router) error {
	port := cfg.GetInt(portKey)

	domain, err := ensureDomainHasNoPort(cfg.GetString(baseURLKey))
	if err != nil {
		return errors.Wrap(err, "error parsing provided image storage url")
	}

	if port == 0 {
		l.BaseURL = domain
	} else {
		l.BaseURL = fmt.Sprintf("%s:%d", domain, port)
	}

	storageDir := cfg.GetString(storageDirKey)
	if storageDir != "" {
		l.StorageDir = storageDir
	}

	if _, err := os.Stat(l.StorageDir); os.IsNotExist(err) {
		return os.MkdirAll(l.StorageDir, os.ModePerm)
	}

	routePrefix := cfg.GetString(routePrefixKey)
	if routePrefix == "" {
		l.RoutePrefix = storageDir
	}

	path := fmt.Sprintf("/%s/", l.RoutePrefix)
	fileServer(router, path, http.Dir(l.StorageDir))
	return nil
}

func (l *localImageStorer) CreateThumbnails(in image.Image) images.ProductImageSet {
	return images.ProductImageSet{
		Thumbnail: resize.Thumbnail(100, 100, in, resize.NearestNeighbor),
		Main:      resize.Thumbnail(500, 500, in, resize.NearestNeighbor),
		Original:  in,
	}
}

func (l *localImageStorer) StoreImages(in images.ProductImageSet, sku string, id uint) (*images.ProductImageLocations, error) {
	photoDir := fmt.Sprintf("%s/%s/%d", l.StorageDir, sku, id)

	var err error
	if _, err = os.Stat(photoDir); os.IsNotExist(err) {
		err = os.MkdirAll(photoDir, os.ModePerm)
		if err != nil {
			return nil, errors.Wrap(err, "error creating necessary folders")
		}
	}
	out := &images.ProductImageLocations{}

	thumbnailPath := fmt.Sprintf("%s/thumbnail.png", photoDir)
	err = saveImage(in.Thumbnail, thumbnailPath)
	if err != nil {
		return nil, err
	}
	out.Thumbnail = fmt.Sprintf("%s/%s", l.BaseURL, thumbnailPath)

	mainPath := fmt.Sprintf("%s/main.png", photoDir)
	err = saveImage(in.Main, mainPath)
	if err != nil {
		return out, err
	}
	out.Main = fmt.Sprintf("%s/%s", l.BaseURL, mainPath)

	originalPath := fmt.Sprintf("%s/original.png", photoDir)
	err = saveImage(in.Original, originalPath)
	if err != nil {
		return out, err
	}
	out.Original = fmt.Sprintf("%s/%s", l.BaseURL, originalPath)

	return out, nil
}
