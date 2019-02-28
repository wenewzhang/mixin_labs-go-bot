package main

import (
  "fmt"
	"github.com/MooooonStar/mixin-sdk-go/messenger"
)
func main() {
  uuid := messenger.UuidNewV4()
  fmt.Println(uuid)
}
