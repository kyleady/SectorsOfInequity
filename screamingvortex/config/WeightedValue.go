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

func WeightedRegionConfigTag() string { return "region config" }
func WeightedSystemInspirationTag() string { return "system inspiration" }

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

func FetchAllWeightedValues(client utilities.ClientInterface, weightedValues *[]*WeightedValue, weightedType string, parentId int64) {
  client.FetchAll(weightedValues, weightedType, "parent_id = ?", parentId)
}
