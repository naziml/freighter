package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"zombiezen.com/go/log"
)

var (
	rootDir    = flag.String("root", "", "Root directory to walk")
	numReaders = flag.Int("readers", 1, "Number of concurrent readers")
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
	flag.Parse()

	if *rootDir == "" {
		fmt.Println("Please specify a root directory to walk")
		return
	}

	root := *rootDir
	totalBytes := uint64(0)
	startTime := time.Now()
	fmt.Println("Listing files in", root)
	filePaths := make(map[int][]string, 0)
	index := 0
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		index += 1
		readerIndex := index % *numReaders
		if m, ok := filePaths[readerIndex]; !ok {
			filePaths[readerIndex] = []string{path}
		} else {
			filePaths[readerIndex] = append(m, path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", root, err)
		return
	}

	wg := sync.WaitGroup{}
	totalBytesReaders := make([]uint64, *numReaders)

	ctx := context.Background()

	for i := 0; i < *numReaders; i++ {
		totalBytesReaders[i] = 0
		wg.Add(1)
		go func(readerIndex int) {
			filesToRead := filePaths[readerIndex]
			log.Infof(ctx, "Reader %d reading %d files", readerIndex, len(filesToRead))
			defer wg.Done()
			for _, path := range filesToRead {
				if info, err := os.Stat(path); err != nil {
					fmt.Printf("error getting file info %q: %v\n", path, err)
				} else {
					if info.IsDir() {
						continue
					}
					if bytes, err := readAllFile(path); err != nil {
						fmt.Printf("error reading file %q: %v\n", path, err)
					} else {
						totalBytesReaders[readerIndex] += bytes
					}
				}
			}
		}(i)
	}

	/*	fmt.Printf("visited file or dir: %q\n", path)
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
		}*/

	wg.Wait()
	for i, bytes := range totalBytesReaders {
		fmt.Printf("Reader %d read %d bytes\n", i, bytes)
		totalBytes += bytes
	}
	endTime := time.Now()
	timeDelta := endTime.Sub(startTime)

	speed := float64(totalBytes) / timeDelta.Seconds()
	fmt.Printf("Time to walk: %v\n", timeDelta)
	fmt.Printf("Total bytes: %d\n", totalBytes)
	fmt.Printf("Speed: %f bytes/sec\n", speed)
}
