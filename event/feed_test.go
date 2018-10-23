package event

import (
	"testing"
	"sync"
	"fmt"
)

func TestFeed_Subscribe(t *testing.T) {
	type someEvent struct{ I int }
	var feed Feed
	var wg sync.WaitGroup

	ch := make(chan someEvent)
	sub := feed.Subscribe(ch)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for event := range ch {
			fmt.Printf("Received: %#v\n", event.I)
		}
		sub.Unsubscribe()
		fmt.Println("done")
	}()

	feed.Send(someEvent{1})
	feed.Send(someEvent{2})
	feed.Send(someEvent{3})
	feed.Send(someEvent{4})
	close(ch)

	wg.Wait()
}