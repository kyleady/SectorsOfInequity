package config

import "fmt"

type InspirationKey struct {
  Key string
  Index int64
  Type string
}

func LogAddress(address []*InspirationKey) {
  fmt.Println("Address: [")
  for _, key := range address {
    fmt.Printf("  %s: \"%s\" [%d],\n", key.Type, key.Key, key.Index)
  }
  fmt.Println("]")
}
