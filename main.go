package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s [path to zips] (optional: target app slug)\n", os.Args[0])
	}

	apps, err := GetCatalog()
	if err != nil {
		panic(err)
	}

	individualSlug := ""
	if len(os.Args) == 3 {
		// Individual creation mode setup
		individualSlug = os.Args[2]
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(apps))
	semaphore := make(chan any, 7)
	for _, app := range apps {
		go func(_app App) {
			defer wg.Done()
			if individualSlug != _app.Slug && individualSlug != "" {
				return
			}

			semaphore <- 1
			fmt.Println(_app.Name)
			err = makeBanner(&_app, os.Args[1])
			if err != nil {
				panic(err)
			}

			<-semaphore
		}(app)
	}

	wg.Wait()
}
