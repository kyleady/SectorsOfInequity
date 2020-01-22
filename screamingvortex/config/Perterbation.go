package config

import "math/rand"

import "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"

type Perterbation struct {
  SystemId int64 `sql:"system_id"`
  StarClusterId int64 `sql:"star_cluster_id"`
  SystemConfig *System
  StarClusterConfig *StarCluster
  Manager *ConfigManager
  Rand *rand.Rand
}

func CreateEmptyPerterbation(client *utilities.Client, rRand *rand.Rand) *Perterbation {
  perterbation := new(Perterbation)
  perterbation.Manager = CreateEmptyManager(client)
  perterbation.Rand = rRand
  perterbation.SystemConfig = CreateEmptySystemConfig()
  perterbation.StarClusterConfig = CreateEmptyStarClusterConfig()
  return perterbation
}

func LoadPerterbationFrom(client utilities.ClientInterface, id int64) *Perterbation {
  perterbation := new(Perterbation)
  client.Fetch(perterbation, "", id)
  perterbation.SystemConfig = LoadSystemConfigFrom(client, perterbation.SystemId)
  perterbation.StarClusterConfig = LoadStarClusterConfigFrom(client, perterbation.StarClusterId)
  return perterbation
}

func (perterbation *Perterbation) TableName(perterbationType string) string {
  return "plan_perterbation"
}

func (perterbation *Perterbation) GetId() *int64 {
  panic("GetId() not implemented. Config should not be editted.")
}

func (basePerterbation *Perterbation) AddInspiration(inspirationId int64) (*Inspiration, *Perterbation) {
  inspiration := basePerterbation.Manager.GetInspiration(inspirationId)
  newPerterbation := basePerterbation.AddPerterbation(inspiration.PerterbationId)

  return inspiration, newPerterbation
}

func (basePerterbation *Perterbation) AddPerterbation(perterbationId int64) *Perterbation {
  newPerterbation := new(Perterbation)
  newPerterbation.Rand = basePerterbation.Rand
  newPerterbation.Manager = basePerterbation.Manager
  modifyingPerterbation := basePerterbation.Manager.GetPerterbation(perterbationId)

  newPerterbation.SystemConfig = basePerterbation.SystemConfig.AddPerterbation(modifyingPerterbation.SystemConfig)
  newPerterbation.StarClusterConfig = basePerterbation.StarClusterConfig.AddPerterbation(modifyingPerterbation.StarClusterConfig)

  return newPerterbation
}
