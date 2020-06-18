package asset

import "screamingvortex/config"
import "screamingvortex/grid"
import "screamingvortex/utilities"

type System struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  GridId int64
  Features []*Detail
  StarClusters []*StarCluster
  Routes []*Route
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

  for _, route := range system.Routes {
    route.SaveParents(client)
  }
  client.SaveAll(&system.Routes, "")
  for _, route := range system.Routes {
    client.Save(&utilities.SystemToRouteLink{ParentId: system.Id, ChildId: route.Id}, "")
  }
}

func RandomSystem(perterbation *config.Perterbation, prefix string, index int, gridSystem *grid.System) *System {
  systemConfig := perterbation.SystemConfig

  system := new(System)
  newPrefix := SetNameAndGetPrefix(system, prefix, index)

  system.GridId = gridSystem.Id

  system.Features, perterbation = RollDetails(systemConfig.SystemFeaturesRolls, systemConfig.WeightedInspirations, systemConfig.ExtraInspirationIds, perterbation)

  numberOfStarClusters := config.RollAll(systemConfig.SystemStarClustersRolls, perterbation.Rand)
  for i := 1; i <= numberOfStarClusters; i++ {
    starCluster := RandomStarCluster(perterbation, newPrefix, i)
    system.StarClusters = append(system.StarClusters, starCluster)
  }

  for i, gridRoute := range gridSystem.Routes {
    route := RandomRoute(perterbation, newPrefix, i+1)
    route.TargetId = gridRoute.EndId
    system.Routes = append(system.Routes, route)
  }

  return system
}
