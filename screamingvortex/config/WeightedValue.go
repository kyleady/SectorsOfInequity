package config

import "math/rand"

type WeightedValue struct {
  Weight int `sql:"weight"`
  Value int64 `sql:"value_id"`
  ValueName string
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
  clonedWeightedValue.ValueName = weightedValue.ValueName
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

func FetchAllWeightedPerterbations(manager *ConfigManager, parentId int64) []*WeightedValue {
  weightedValues := make([]*WeightedValue, 0)
  manager.Client.FetchAll(&weightedValues, WeightedPerterbationTag(), "parent_id = ?", parentId)
  return weightedValues
}

func FetchManyWeightedInspirations(manager *ConfigManager, parentId int64, tableName string, valueName string) []*WeightedValue {
  weightedValues := make([]*WeightedValue, 0)
  weightTableName := new(WeightedValue).TableName(WeightedInspirationTag())
  manager.Client.FetchMany(&weightedValues, parentId, tableName, weightTableName, valueName, WeightedInspirationTag(), false)
  return weightedValues
}

func WeightedPerterbationTag() string { return "weighted perterbation" }
func WeightedInspirationTag() string { return "weighted inspiration" }
