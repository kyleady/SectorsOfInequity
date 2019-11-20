package utilities

import "math/rand"
import "reflect"

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

  weightedRoll := rand.Intn(totalWeight)
  for i := 0; i < asSlice.Len(); i++ {
    obj := asSlice.Index(i).Interface().(WeightedObj)
    if weightedRoll <= 0 {
      return obj.GetValue()
    }

    if obj.GetWeight() > 0 {
      weightedRoll -= obj.GetWeight()
    }
  }

  return nil
}
