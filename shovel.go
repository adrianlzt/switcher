package main

import (
	"io"
	"log"
	"strings"
)

// proxy between two sockets
func Shovel(local, remote io.ReadWriteCloser) error {
	errch := make(chan error, 1)

	defer func() {
		go local.Close()
	}()

	defer func() {
		go remote.Close()
	}()

	go chanCopy(errch, local, remote)
	go chanCopy(errch, remote, local)

	<-errch
	return nil
}

// copy between pipes, sending errors to channel
func chanCopy(e chan error, dst, src io.ReadWriter) {
	_, err := io.Copy(dst, src)
	if err != nil && !strings.HasSuffix(err.Error(), ": use of closed network connection") {
		log.Printf("[ERROR] chanCopy: %v\n", err)
	}
	e <- err
}
