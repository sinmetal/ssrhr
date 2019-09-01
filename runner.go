package main

import (
	"context"
	"fmt"
	"sync"
)

func goGetSmallSize(sss SmallSizeStore, goroutine int, rowNumber int, endCh chan<- error) {
	go func() {
		for {
			fmt.Printf("rowNumber=%d\n", rowNumber)
			var wg sync.WaitGroup
			for i := 0; i < goroutine; i++ {
				i := i
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					fmt.Printf("goroutine=%d,rowNumber=%d\n", i, rowNumber)
					ctx := context.Background()
					ctx, span := startSpan(ctx, "/go/goGetSmallSize")
					defer span.End()

					_, err := sss.Get(ctx, fmt.Sprintf("small%v", rowNumber))
					if err != nil {
						endCh <- err
					}
				}(i)
			}
			wg.Wait()
		}
	}()
}
