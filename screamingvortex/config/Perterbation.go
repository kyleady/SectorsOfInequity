package config

import "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"

type Perterbation struct {
  SystemId int64 `sql:"system_id"`
  SystemConfig *System
}

func LoadPerterbationFrom(client utilities.ClientInterface, perterbationType string, id int64) *Perterbation {
  perterbation := new(Perterbation)
  client.Fetch(perterbation, perterbationType, id)
  client.Fetch(perterbation.SystemConfig, "", perterbation.SystemId)

  return perterbation
}

func (perterbation *Perterbation) TableName(perterbationType string) string {
  switch perterbationType {
  case PerterbationRegionTag():
    return "plan_config_region"
  case PerterbationSystemTag():
    return "plan_perterbation_system"
  default:
    panic("Unexpected perterbationType.")
  }
}

func (perterbation *Perterbation) GetId() *int64 {
  panic("GetId() not implemented. Config should not be editted.")
}

func (basePerterbation *Perterbation) AddPerterbation(perterbation *Perterbation) *Perterbation {
  newPerterbation := new(Perterbation)
  newPerterbation.SystemConfig = basePerterbation.SystemConfig.AddPerterbation(perterbation.SystemConfig)

  return newPerterbation
}

func PerterbationRegionTag() string { return "region" }
func PerterbationSystemTag() string { return "system" }
