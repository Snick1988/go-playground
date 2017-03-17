package main

import (
    "fmt"
    "time"
    "math/rand"
)

type Worker struct {
    id int
}

func (w Worker) process(c chan int) {
    for {
        data := <-c
        fmt.Println(data)
        time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
    }
}

func main() {
    c := make(chan int, 100)

    for i := 0; i < 10; i++ {
        worker := &Worker{id: i}
        go worker.process(c)
    }

    for {
        select {
            case c <- rand.Int():    
            case t := <-time.After(time.Millisecond * 100):
                fmt.Println("Timed out at", t)
        }
        time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)))
    }
}
