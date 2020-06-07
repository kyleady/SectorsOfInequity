package asset

import "github.com/kyleady/SectorsOfInequity/screamingvortex/config"
import "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"

type System struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  Features []*Detail
  StarClusters []*StarCluster
}

func (system *System) TableName(systemType string) string {
  return "plan_asset_system"
}

func (system *System) GetId() *int64 {
  return &system.Id
}

func (system *System) GetType() string {
  return "System"
}

func (system *System) SetName(name string) {
  system.Name = name
}

func (system *System) SaveTo(client utilities.ClientInterface) {
  client.Save(system, "")
  system.SaveChildren(client)
}

func (system *System) SaveChildren(client utilities.ClientInterface) {
  client.SaveAll(&system.Features, "")
  for _, feature := range system.Features {
    client.Save(&utilities.SystemToDetailLink{ParentId: system.Id, ChildId: feature.Id}, "")
    feature.SaveChildren(client)
  }

  client.SaveAll(&system.StarClusters, "")
  for _, starCluster := range system.StarClusters {
    client.Save(&utilities.SystemToStarClusterLink{ParentId: system.Id, ChildId: starCluster.Id}, "")
    starCluster.SaveChildren(client)
  }
}

func RandomSystem(perterbation *config.Perterbation, prefix string, index int) *System {
  systemConfig := perterbation.SystemConfig

  system := new(System)
  newPrefix := SetNameAndGetPrefix(system, prefix, index)

  system.Features, perterbation = RollDetails(systemConfig.SystemFeaturesRolls, systemConfig.WeightedInspirations, perterbation)
  numberOfStarClusters := config.RollAll(systemConfig.SystemStarClustersRolls, perterbation.Rand)
  for i := 1; i <= numberOfStarClusters; i++ {
    starCluster := RandomStarCluster(perterbation, newPrefix, i)
    system.StarClusters = append(system.StarClusters, starCluster)
  }

  return system
}
