package main

import (
	"context"
	"sync"
)

func goGetSmallSize(sss SmallSizeStore, goroutine int, endCh chan<- error) {
	go func() {
		for {
			var wg sync.WaitGroup
			for i := 0; i < goroutine; i++ {
				i := i
				wg.Add(1)
				go func(i int) {
					defer wg.Done()

					ctx := context.Background()
					ctx, span := startSpan(ctx, "/go/goGetTweet3Tables")
					defer span.End()

					_, err := sss.Get(ctx, "small1")
					if err != nil {
						endCh <- err
					}
				}(i)
			}
			wg.Wait()
		}
	}()
}
