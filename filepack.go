package filepack

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func GetFiles(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	files := make([]string, 0)
	for _, entry := range entries {
		fileName := fmt.Sprintf("%s/%s", path, entry.Name())
		if entry.IsDir() {
			_files, err := GetFiles(fileName)
			if err != nil {
				return nil, err
			}
			files = append(files, _files...)
			continue
		}
		files = append(files, fileName)
	}
	return files, nil
}

func Move(src string, target string) error {
	files, err := GetFiles(src)
	if err != nil {
		return err
	}
	for _, file := range files {
		_file := strings.Replace(file, src, target, 1)
		segments := strings.FieldsFunc(_file, func(r rune) bool {
			return r == '/' || r == '\\'
		})
		path := strings.Join(segments[:len(segments)-1], "/")
		_, err := os.ReadDir(path)
		if os.IsNotExist(err) {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return err
			}
		}
		err = os.Rename(file, _file)
		if err != nil {
			return err
		}
	}
	root := strings.FieldsFunc(strings.Replace(src, filepath.VolumeName(src), "", 1), func(r rune) bool {
		return r == '/' || r == '\\'
	})[0]
	err = os.RemoveAll(root)
	if err != nil {
		return err
	}
	return nil
}

func Copy(src string, target string) error {
	files, err := GetFiles(src)
	if err != nil {
		return err
	}
	for _, file := range files {
		_file := strings.Replace(file, src, target, 1)
		segments := strings.FieldsFunc(_file, func(r rune) bool {
			return r == '/' || r == '\\'
		})
		path := strings.Join(segments[:len(segments)-1], "/")
		_, err := os.ReadDir(path)
		if os.IsNotExist(err) {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return err
			}
		}
		r, err := os.Open(file)
		if err != nil {
			return err
		}
		w, err := os.Create(_file)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, r)
		if err != nil {
			return err
		}
	}
	return nil
}
