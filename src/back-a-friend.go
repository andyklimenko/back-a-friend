package main

import (
	"fmt"
	"path/filepath"
	"os"

	"server"
)

func main() {
	currDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		return
	}

	doneCh, err := server.StartServer(currDir)
	if err != nil {
		fmt.Println(err)
	}

	<-doneCh
}
