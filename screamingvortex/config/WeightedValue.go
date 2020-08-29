package config

import "strconv"
import "sort"

type WeightedValue struct {
  Id int64 `sql:"id"`
  Weights []*Roll
  Weight int
  Order []*Roll
  TmpOrder int
  Value int64 `sql:"value_id"`
  ValueName string
  Values []int64
}

func (weightedValue *WeightedValue) TableName(weightedType string) string {
  switch weightedType {
  case WeightedRegionTag():
    return "plan_weighted_region"
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
  clonedWeightedValue.TmpOrder = weightedValue.TmpOrder
  clonedWeightedValue.Value  = weightedValue.Value
  clonedWeightedValue.ValueName = weightedValue.ValueName
  clonedWeightedValue.Values = make([]int64, len(weightedValue.Values))
  copy(clonedWeightedValue.Values, weightedValue.Values)
  clonedWeightedValue.Weights = make([]*Roll, len(weightedValue.Weights))
  copy(clonedWeightedValue.Weights, weightedValue.Weights)
  clonedWeightedValue.Order = make([]*Roll, len(weightedValue.Order))
  copy(clonedWeightedValue.Order, weightedValue.Order)
  return clonedWeightedValue
}

func (weightedValue *WeightedValue) rollWeight(perterbation *Perterbation) {
  weightedValue.Weight = RollAll(weightedValue.Weights, perterbation)
}

func SortWeightedValues(weightedValues []*WeightedValue, perterbation *Perterbation) {
  for _, weightedValue := range weightedValues {
    weightedValue.TmpOrder = RollAll(weightedValue.Order, perterbation)
  }

  sort.Slice(weightedValues, func(i, j int) bool {
		return weightedValues[i].TmpOrder < weightedValues[j].TmpOrder
	})
}

func RollWeightedValues(weightedValues []*WeightedValue, perterbation *Perterbation, modifiers []*Roll) *WeightedValue {
  totalWeight := 0
  rRand := perterbation.Rand
  for _, weightedValue := range weightedValues {
    weightedValue.rollWeight(perterbation)
    if weightedValue.Weight > 0 {
      totalWeight += weightedValue.Weight
    }
  }

  if totalWeight <= 0 {
    return nil
  }

  weightedRoll := rRand.Intn(totalWeight) + RollAll(modifiers, perterbation)
  SortWeightedValues(weightedValues, perterbation)
  weightedValue := new(WeightedValue)
  for _, weightedValue = range weightedValues {
    if weightedValue.Weight <= 0 {
      continue
    }

    if weightedRoll < weightedValue.Weight {
      return weightedValue
    } else {
      weightedRoll -= weightedValue.Weight
    }
  }

  return weightedValue
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
        newWeightedValue.Weights = append(newWeightedValue.Weights, secondWeightedValue.Weights...)
        for _, value := range secondWeightedValue.Values {
          valueAlreadyInNewValues := false
          for _, newValue := range newWeightedValue.Values {
              if value == newValue {
                valueAlreadyInNewValues = true
                break
              }
          }

          if !valueAlreadyInNewValues {
            newWeightedValue.Values = append(newWeightedValue.Values, value)
          }
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

func fetchManyWeightedValues(manager *ConfigManager, parentId int64, tableName string, valueName string, weightedTag string) []*WeightedValue {
  weightedValues := make([]*WeightedValue, 0)
  weightTableName := new(WeightedValue).TableName(weightedTag)
  manager.Client.FetchMany(&weightedValues, parentId, tableName, weightTableName, valueName, weightedTag, false)
  for _, weightedValue := range weightedValues {
    weightedValue.Weights = FetchManyRolls(manager, weightedValue.Id, weightedValue.TableName(weightedTag), "weights")
    weightedValue.Order = FetchManyRolls(manager, weightedValue.Id, weightedValue.TableName(weightedTag), "order")
    weightedValue.Values = append(weightedValue.Values, weightedValue.Value)
  }

  return weightedValues
}

func FetchManyWeightedInspirations(manager *ConfigManager, parentId int64, tableName string, valueName string) []*WeightedValue {
  weightedValues := fetchManyWeightedValues(manager, parentId, tableName, valueName, WeightedInspirationTag())
  for _, weightedValue := range weightedValues {
    value := new(Inspiration)
    manager.Client.Fetch(value, "", weightedValue.Value)
    weightedValue.ValueName = value.Name
  }

  return weightedValues
}

func FetchManyWeightedRegions(manager *ConfigManager, parentId int64, tableName string, valueName string) []*WeightedValue {
  weightedValues := fetchManyWeightedValues(manager, parentId, tableName, valueName, WeightedRegionTag())
  for _, weightedValue := range weightedValues {
    value := new(RegionConfig)
    manager.Client.Fetch(value, "", weightedValue.Value)
    weightedValue.ValueName = strconv.FormatInt(value.TypeId, 10) + ":" + value.Name
  }

  return weightedValues
}

func WeightedRegionTag() string { return "weighted region" }
func WeightedInspirationTag() string { return "weighted inspiration" }
