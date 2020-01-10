package config

import "math/rand"

import "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"

type WeightedValue struct {
  Weight int `sql:"weight"`
  Value int64 `sql:"value_id"`
}

func (weightedValue *WeightedValue) TableName(weightedType string) string {
  switch weightedType {
  case WeightedRegionConfigTag():
    return "plan_weighted_config_region"
  case WeightedSystemInspirationTag():
    return "plan_weighted_inspiration_system"

  default:
    panic("Unexpected regionType.")
  }
}

func (weightedValue *WeightedValue) GetId() *int64 {
  panic("GetId() not implemented. Config should not be editted.")
}

func (weightedValue *WeightedValue) Clone() *WeightedValue {
  clonedWeightedValue := new(WeightedValue)
  clonedWeightedValue.Weight = weightedValue.Weight
  clonedWeightedValue.Value  = weightedValue.Value
  return clonedWeightedValue
}

func RollWeightedValues(weightedValues []*WeightedValue, rRand *rand.Rand) int64 {
  totalWeight := 0
  for _, weightedValue := range weightedValues {
    if weightedValue.Weight > 0 {
      totalWeight += weightedValue.Weight
    }
  }

  if totalWeight <= 0 {
    return -1
  }

  weightedRoll := rRand.Intn(totalWeight)
  for _, weightedValue := range weightedValues {
    if weightedRoll < weightedValue.Weight {
      return weightedValue.Value
    } else {
      weightedRoll -= weightedValue.Weight
    }
  }

  panic("RollWeightedValues should always return a value!")
}

func StackWeightedValues(firstWeightedValues []*WeightedValue, secondWeightedValues []*WeightedValue) []*WeightedValue {
  newWeightedValues := make([]*WeightedValue, len(firstWeightedValues))
  for i, firstWeightedValue := range firstWeightedValues {
    newWeightedValues[i] = firstWeightedValue.Clone()
  }

  for _, secondWeightedValue := range secondWeightedValues {
    for _, newWeightedValue := range newWeightedValues {
      if newWeightedValue.Value == secondWeightedValue.Value {
        newWeightedValue.Weight += secondWeightedValue.Weight
        continue
      }
    }

    newWeightedValues = append(newWeightedValues, secondWeightedValue.Clone())
  }

  return newWeightedValues
}

func FetchAllWeightedValues(client utilities.ClientInterface, weightedValues *[]*WeightedValue, weightedType string, parentId int64) {
  client.FetchAll(weightedValues, weightedType, "parent_id = ?", parentId)
}
