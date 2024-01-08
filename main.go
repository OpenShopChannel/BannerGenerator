package main

import (
	"fmt"
	"sync"
)

func main() {
	apps, err := GetCatalog()
	if err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(apps))
	semaphore := make(chan any, 7)
	for _, app := range apps {
		go func(_app App) {
			defer wg.Done()
			semaphore <- 1

			fmt.Println(_app.Name)
			err = makeBanner(&_app)
			if err != nil {
				panic(err)
			}

			<-semaphore
		}(app)
	}

	wg.Wait()
}
