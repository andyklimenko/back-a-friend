package server

import "time"

func StartServer() (chan struct{}, error) {
	doneCh := make(chan struct{})
	go func() {
		// init http-server
		time.Sleep(time.Second)
		close(doneCh)
	}()
	return doneCh, nil
}
