package config

import "math/rand"

type WeightedValue struct {
  Weight int `sql:"weight"`
  Value int64 `sql:"value_id"`
  ValueName string
  Values []int64
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
  clonedWeightedValue.Values = make([]int64, len(weightedValue.Values))
  copy(clonedWeightedValue.Values, weightedValue.Values)
  return clonedWeightedValue
}

func RollWeightedValues(weightedValues []*WeightedValue, rRand *rand.Rand) []int64 {
  totalWeight := 0
  for _, weightedValue := range weightedValues {
    if weightedValue.Weight > 0 {
      totalWeight += weightedValue.Weight
    }
  }

  if totalWeight <= 0 {
    return []int64{}
  }

  weightedRoll := rRand.Intn(totalWeight)
  for _, weightedValue := range weightedValues {
    if weightedRoll < weightedValue.Weight {
      return weightedValue.Values
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
    weightedValueStacked := false
    for _, newWeightedValue := range newWeightedValues {
      if newWeightedValue.ValueName == secondWeightedValue.ValueName {
        weightedValueStacked = true
        newWeightedValue.Weight += secondWeightedValue.Weight
        if newWeightedValue.Value != secondWeightedValue.Value {
          newWeightedValue.Values = append(newWeightedValue.Values, secondWeightedValue.Value)
        }

        break
      }
    }

    if !weightedValueStacked {
      newWeightedValues = append(newWeightedValues, secondWeightedValue.Clone())
    }
  }

  return newWeightedValues
}

func FetchAllWeightedPerterbations(manager *ConfigManager, parentId int64) []*WeightedValue {
  weightedValues := make([]*WeightedValue, 0)
  manager.Client.FetchAll(&weightedValues, WeightedPerterbationTag(), "parent_id = ?", parentId)
  for _, weightedValue := range weightedValues {
    weightedValue.Values = append(weightedValue.Values, weightedValue.Value)
  }

  return weightedValues
}

func FetchManyWeightedInspirations(manager *ConfigManager, parentId int64, tableName string, valueName string) []*WeightedValue {
  weightedValues := make([]*WeightedValue, 0)
  weightTableName := new(WeightedValue).TableName(WeightedInspirationTag())
  manager.Client.FetchMany(&weightedValues, parentId, tableName, weightTableName, valueName, WeightedInspirationTag(), false)
  for _, weightedValue := range weightedValues {
    inspiration := manager.GetInspiration(weightedValue.Value)
    weightedValue.ValueName = inspiration.Name
    weightedValue.Values = append(weightedValue.Values, weightedValue.Value)
  }

  return weightedValues
}

func WeightedPerterbationTag() string { return "weighted perterbation" }
func WeightedInspirationTag() string { return "weighted inspiration" }
