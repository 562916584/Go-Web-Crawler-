package persist

import "log"

func ItemSaver() chan interface{} {
	out := make(chan interface{})
	go func() {
		itemCount := 1
		for {
			item := <-out
			log.Printf("Saver item :#%d:  %v\n", itemCount, item)
			itemCount++
		}
	}()
	return out
}
