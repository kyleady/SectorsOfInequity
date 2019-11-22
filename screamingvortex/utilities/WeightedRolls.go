package utilities

import "math/rand"
import "reflect"
import "fmt"

type WeightedObj interface {
  GetWeight() int
  GetValue() interface{}
}

func WeightedRoll(asInterface interface{}, rand *rand.Rand) interface{} {
  asSlice := reflect.ValueOf(asInterface).Elem()
  totalWeight := 0
  for i := 0; i < asSlice.Len(); i++ {
    obj := asSlice.Index(i).Interface().(WeightedObj)
    if obj.GetWeight() > 0 {
      totalWeight += obj.GetWeight()
    }
  }

  if totalWeight <= 0 {
    panic(fmt.Sprintf("WeightedRoll received a weight of %d.", totalWeight))
  }

  weightedRoll := rand.Intn(totalWeight)
  for i := 0; i < asSlice.Len(); i++ {
    obj := asSlice.Index(i).Interface().(WeightedObj)
    if weightedRoll < obj.GetWeight() {
      return obj.GetValue()
    } else {
      weightedRoll -= obj.GetWeight()
    }
  }

  panic("WeightedRoll should always return a value!")
}
