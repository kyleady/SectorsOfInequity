package config

import "math/rand"

import "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"

type WeightedValue struct {
  Weight int `sql:"weight"`
  Value int64 `sql:"value_id"`
}

func (weightedValue *WeightedValue) TableName(weightedType string) string {
  switch weightedType {
  case WeightedPerterbationTag():
    return "plan_weighted_perterbation"
  case WeightedInspirationTag():
    return "plan_weighted_inspiration"

  default:
    panic("Unexpected weightedType: " + weightedType)
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

func FetchAllWeightedPerterbations(client utilities.ClientInterface, weightedValues *[]*WeightedValue, parentId int64) {
  client.FetchAll(weightedValues, WeightedPerterbationTag(), "parent_id = ?", parentId)
}

func FetchAllWeightedInspirations(client utilities.ClientInterface, weightedValues *[]*WeightedValue, parentId int64, tableName string, valueName string) {
  weightTableName := new(WeightedValue).TableName(WeightedInspirationTag())
  client.FetchMany(weightedValues, parentId, tableName, weightTableName, valueName, WeightedInspirationTag(), false)
}


func WeightedPerterbationTag() string { return "weighted perterbation" }
func WeightedInspirationTag() string { return "weighted inspiration" }
