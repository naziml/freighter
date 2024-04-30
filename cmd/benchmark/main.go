package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

func prepareTestDirTree(tree string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", fmt.Errorf("error creating temp directory: %v\n", err)
	}

	err = os.MkdirAll(filepath.Join(tmpDir, tree), 0755)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", err
	}

	return tmpDir, nil
}

func readAllFile(path string) (uint64, error) {
	bytes := uint64(0)
	if f, err := os.Open(path); err != nil {
		return 0, err
	} else {
		defer f.Close()
		data := make([]byte, 4096)
		for {
			data = data[:cap(data)]
			n, err := f.Read(data)
			if err != nil {
				if err == io.EOF {
					break
				}
				return 0, err
			}
			data = data[:n]
			bytes += uint64(n)
		}
	}
	return bytes, nil
}

func main() {
	root := "/Users/johnewart/Temp/container"
	totalBytes := uint64(0)
	startTime := time.Now()
	fmt.Println("Listing files in", root)
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		fmt.Printf("visited file or dir: %q\n", path)
		if bytes, err := readAllFile(path); err != nil {
			fmt.Printf("error reading file %q: %v\n", path, err)
		} else {
			totalBytes += bytes
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", root, err)
		return
	}

	endTime := time.Now()
	timeDelta := endTime.Sub(startTime)
	speed := float64(totalBytes) / timeDelta.Seconds()
	fmt.Printf("Time to walk: %v\n", timeDelta)
	fmt.Printf("Total bytes: %d\n", totalBytes)
	fmt.Printf("Speed: %f bytes/sec\n", speed)
}
