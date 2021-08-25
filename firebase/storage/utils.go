package storage

import "sync"

// Handle Goroutine Errors, Wait until either WaitGroup is done or an error is received through the channel
func HandleGoroutineErrors(wg *sync.WaitGroup, errChan chan error) error {
	wgDone := make(chan bool)
	go func() {
		wg.Wait()
		close(wgDone)
	}()

	// Wait until either WaitGroup is done or an error is received through the channel
	select {
	case <-wgDone:
		return nil
	case err := <-errChan:
		close(errChan)
		return err
	}
}

