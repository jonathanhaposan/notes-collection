package main

import (
	"context"
	"errors"
	"log"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := make(chan error, 1)
	result := make(chan string, 1)

	go func() {
		r, c := logging(ctx)
		result <- r
		err <- c
	}()

	select {
	case <-ctx.Done():
		log.Println("log", ctx.Err())
	case <-err:
		log.Println("err")
	}
	ss := <-result

	log.Println(ss)
}

func logging(ctx context.Context) (str string, err error) {
	time.Sleep(2 * time.Second)
	return "asd", errors.New("errr")
}
