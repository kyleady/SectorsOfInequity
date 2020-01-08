package config

import "math/rand"

import "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"

type Perterbation struct {
  SystemId int64 `sql:"system_id"`
  SystemConfig *System
  Manager *ConfigManager
  Rand *rand.Rand
}

func CreateEmptyPerterbation(client *utilities.Client, rRand *rand.Rand) *Perterbation {
  perterbation := new(Perterbation)
  perterbation.Manager = CreateEmptyManager(client)
  perterbation.Rand = rRand
  return perterbation
}

func LoadPerterbationFrom(client utilities.ClientInterface, perterbationType string, id int64) *Perterbation {
  perterbation := new(Perterbation)
  client.Fetch(perterbation, perterbationType, id)
  perterbation.SystemConfig = LoadSystemConfigFrom(client, perterbation.SystemId)

  return perterbation
}

func (perterbation *Perterbation) TableName(perterbationType string) string {
  switch perterbationType {
  case PerterbationRegionTag():
    return "plan_config_region"
  case InspirationSystemFeatureTag():
    return "plan_perterbation_system"
  default:
    panic("Unexpected perterbationType.")
  }
}

func (perterbation *Perterbation) GetId() *int64 {
  panic("GetId() not implemented. Config should not be editted.")
}

func (basePerterbation *Perterbation) AddInspiration(inspirationType string, inspirationId int64) (*Inspiration, *Perterbation) {
  inspiration := basePerterbation.Manager.GetInspiration(inspirationType, inspirationId)
  newPerterbation := basePerterbation.AddPerterbation(inspirationType, inspiration.PerterbationId)

  return inspiration, newPerterbation
}

func (basePerterbation *Perterbation) AddPerterbation(perterbationType string, perterbationId int64) *Perterbation {
  newPerterbation := new(Perterbation)
  newPerterbation.Rand = basePerterbation.Rand
  modifyingPerterbation := basePerterbation.Manager.GetPerterbation(perterbationType, perterbationId)

  newPerterbation.SystemConfig = basePerterbation.SystemConfig.AddPerterbation(modifyingPerterbation.SystemConfig)

  return newPerterbation
}
