package main

import (
	"fmt"

	"server"
)

func main() {
	doneCh, err := server.StartServer()
	if err != nil {
		fmt.Println(err)
	}

	<-doneCh
}
