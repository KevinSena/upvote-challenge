package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("My Uri: ", os.Getenv("MONGO_URI"))
}
