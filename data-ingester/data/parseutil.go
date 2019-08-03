package data

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/mholt/archiver"
)

// UnzipAndSave unzips file at the specified filepath and saves it to the databse
func UnzipAndSave(filepath string) error {
	z := archiver.NewZip()

	return z.Walk(filepath, archiveVisitor)
}

func archiveVisitor(f archiver.File) error {
	if f.IsDir() {
		return nil
	}

	if path.Ext(f.Name()) != ".zip" {
		err := parseAndSave(f)
		if err != nil {
			return err
		}

		return nil
	}

	z1 := archiver.NewZip()
	b, _ := ioutil.ReadAll(f.ReadCloser)

	readerAt := bytes.NewReader(b)
	err := z1.Open(readerAt, f.Size())
	defer z1.Close()
	if err != nil {
		return err
	}

	file, err := z1.Read()
	if err != nil {
		return err
	}

	err = parseAndSave(file)
	if err != nil {
		return err
	}

	return nil
}

func parseAndSave(f io.Reader) error {
	devices, err := Parse(f)
	if err != nil {
		return err
	}

	err = SaveDeviceData(*devices)
	if err != nil {
		return err
	}

	return nil
}

// DownloadFile downloads a file from a url and saves it to a temporary file
func DownloadFile(url string) (string, error) {
	tmp, err := ioutil.TempFile(os.TempDir(), "")

	if err != nil {
		return "", err
	}

	resp, err := http.Get(url)
	if err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return "", err
	}
	defer resp.Body.Close()

	_, err = io.Copy(tmp, resp.Body)
	if err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return "", err
	}

	tmp.Close()
	return tmp.Name(), nil
}
