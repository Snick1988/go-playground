package main

import (
    "fmt"
    "db"
)

func main() {
    item := db.LoadItem(1)
    fmt.Printf("%s\n", item)
}
