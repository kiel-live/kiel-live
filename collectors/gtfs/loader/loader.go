package loader

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/artonge/go-gtfs"
)

func downloadFile(url string, out *os.File) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func unzip(src string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "gtfs-collector-*")
	if err != nil {
		return "", err
	}

	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(tmpDir, f.Name)

		if f.FileInfo().IsDir() {
			err := os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				return "", err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return "", err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return "", err
		}

		rc, err := f.Open()
		if err != nil {
			return "", err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return "", err
		}
	}

	return tmpDir, nil
}

func loadRemoteGTFS(url string) (*gtfs.GTFS, error) {
	tmpFile, err := os.CreateTemp("", "gtfs-collector-*.zip")
	if err != nil {
		return nil, err
	}

	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}()

	// download file
	err = downloadFile(url, tmpFile)
	if err != nil {
		return nil, err
	}

	// unzip file
	gtfsFolder, err := unzip(tmpFile.Name())
	if err != nil {
		return nil, err
	}

	return gtfs.Load(gtfsFolder, nil)
}

func LoadGTFS(path string) (*gtfs.GTFS, error) {
	// if url => download, unzip and load
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return loadRemoteGTFS(path)
	}

	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// if directory => load
	if stat.IsDir() {
		return gtfs.Load(path, nil)
	}

	// if file => unzip
	gtfsFolder, err := unzip(path)
	if err != nil {
		return nil, err
	}

	// remove extracted folder after loading
	defer os.RemoveAll(gtfsFolder)

	return gtfs.Load(gtfsFolder, nil)
}
