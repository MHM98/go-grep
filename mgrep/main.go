package main

import (
	"fmt"
	"mgrep/worker"
	"mgrep/workerlist"
	"os"
	"path/filepath"
	"sync"

	"github.com/alexflint/go-arg"
)

func discoverDirs(wl *workerlist.Worklist, path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error while reading directory. err is: %v", err)
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			nextPath := filepath.Join(path, entry.Name())
			discoverDirs(wl, nextPath)
		} else {
			wl.Add(workerlist.NewJob(filepath.Join(path, entry.Name())))
		}
	}
}

var args struct {
	SearchTerm string `arg:"positional,required"`
	SearchPath string `arg:"positional"`
}

func main() {
	arg.MustParse(&args)

	var workersWg sync.WaitGroup

	wl := workerlist.New(100)

	results := make(chan worker.Result, 100)
	numWorkers := 10

	workersWg.Add(1)

	go func() {
		defer workersWg.Done()
		discoverDirs(&wl, args.SearchPath)
		wl.Finilize(numWorkers)
	}()
//thisistest
	for i := 0; i < numWorkers; i++ {
		workersWg.Add(1)
		go func() {
			defer workersWg.Done()
			for {
				workEntry := wl.Next()
				if workEntry.Path != "" {
					workerResult := worker.FindInFile(workEntry.Path, args.SearchTerm)
					if workerResult != nil {
						for _, v := range workerResult.Inner {
							results <- v
						}
					}
				} else {
					return
				}
			}
		}()
	}

	blockWorkerChan := make(chan struct{})

	go func() {
		workersWg.Wait()
		close(blockWorkerChan)
	}()

	var displayWg sync.WaitGroup
	displayWg.Add(1)

	go func() {
		for {
			select {
			case r := <-results:
				fmt.Printf("%v[%v]%v\n", r.Path, r.LineNum, r.Line)
			case <-blockWorkerChan:
				if len(results) < 1 {
					displayWg.Done()
					return
				}
			}
		}
	}()
	displayWg.Wait()
}
