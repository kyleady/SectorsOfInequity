package utilities

import "math/rand"

type WeightedValue struct {
  Id int64 `sql:"id"`
  Weight int `sql:"weight"`
  Value int64 `sql:"value_id"`
}

func (weightedValue *WeightedValue) TableName(weightedType string) string {
  switch weightedType {
  case "region":
    return "plan_weighted_config_region"
  default:
    panic("Unexpected regionType.")
  }
}

func (weightedValue *WeightedValue) GetId() *int64 {
  return &weightedValue.Id
}

func RollWeightedValues(weightedValues []*WeightedValue, rand *rand.Rand) int64 {
  totalWeight := 0
  for _, weightedValue := range weightedValues {
    if weightedValue.Weight > 0 {
      totalWeight += weightedValue.Weight
    }
  }

  if totalWeight <= 0 {
    return -1
  }

  weightedRoll := rand.Intn(totalWeight)
  for _, weightedValue := range weightedValues {
    if weightedRoll < weightedValue.Weight {
      return weightedValue.Value
    } else {
      weightedRoll -= weightedValue.Weight
    }
  }

  panic("RollWeightedValues should always return a value!")
}
