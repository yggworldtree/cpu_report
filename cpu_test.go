package main

import (
	"testing"
	"time"
)

func Test1(t *testing.T) {
	for i := 0; i < 1000; i++ {
		go func() {
			for {
				a := 1 + 1
				println(i, ":", a)
			}
		}()
	}
	time.Sleep(time.Second * 40)
}
