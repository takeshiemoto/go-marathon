package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go runWorker(ctx, &wg, i)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Press 's' to stop workers...")
	input := make(chan string, 1)
	go func() {
		var s string
		fmt.Scanln(&s)
		input <- s
	}()

	select {
	case sig := <-signals:
		fmt.Printf("Signal receiverd: %v\n", sig)
		cancel()
	case s := <-input:
		if s == "S" {
			cancel()
		} else {
			fmt.Println("Invalid input:", s)
		}
	}

	wg.Wait()
	fmt.Println("All workers stopped")
}

func runWorker(ctx context.Context, wg *sync.WaitGroup, id int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Workder %d: stopped\n", id)
			return
		default:
			fmt.Printf("Workder %d: Running...\n", id)
			time.Sleep(500 * time.Millisecond)
		}
	}
}
