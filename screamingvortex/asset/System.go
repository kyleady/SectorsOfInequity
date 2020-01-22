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

func (system *System) SaveTo(client utilities.ClientInterface, parentId int64) {
  client.Save(system, "")
  system.SaveChildren(client, parentId)
}

func (system *System) SaveChildren(client utilities.ClientInterface, parentId int64) {
  client.SaveAll(&system.Features, "")
  for _, feature := range system.Features {
    client.Save(&utilities.SystemToDetailLink{ParentId: system.Id, ChildId: feature.Id}, "")
  }

  client.SaveAll(&system.StarClusters, "")
  for _, starCluster := range system.StarClusters {
    client.Save(&utilities.SystemToStarClusterLink{ParentId: system.Id, ChildId: starCluster.Id}, "")
    //starCluster.SaveChildren(client, system.Id)
  }
}

func RandomSystem(perterbation *config.Perterbation, prefix string, index int) *System {
  systemConfig := perterbation.SystemConfig

  system := new(System)
  newPrefix := SetNameAndGetPrefix(system, prefix, index)
  numberOfSystemFeatures := config.RollAll(systemConfig.SystemFeaturesRolls, perterbation.Rand)
  numberOfStarClusters := config.RollAll(systemConfig.SystemStarClustersRolls, perterbation.Rand)

  for i := 1; i <= numberOfSystemFeatures; i++ {
    inspirationId := config.RollWeightedValues(systemConfig.WeightedInspirations, perterbation.Rand)
    inspiration, newPerterbation := perterbation.AddInspiration(inspirationId)
    systemFeature := RandomDetail(inspiration, perterbation.Rand)

    system.Features = append(system.Features, systemFeature)
    perterbation = newPerterbation
  }

  for i := 1; i <= numberOfStarClusters; i++ {
    starCluster := RandomStarCluster(perterbation, newPrefix, i)

    system.StarClusters = append(system.StarClusters, starCluster)
  }

  return system
}
