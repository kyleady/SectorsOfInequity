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
  client.SaveMany2ManyLinks(system, &system.Features, "", "", "details", false)
  for _, feature := range system.Features {
    feature.SaveChildren(client)
  }

  client.SaveAll(&system.StarClusters, "")
  client.SaveMany2ManyLinks(system, &system.StarClusters, "", "", "star_clusters", false)
  for _, starCluster := range system.StarClusters {
    starCluster.SaveChildren(client)
  }

  for _, route := range system.Routes {
    route.SaveParents(client)
  }
  client.SaveAll(&system.Routes, "")
  client.SaveMany2ManyLinks(system, &system.Routes, "", "", "routes", false)
}

func RandomSystem(perterbation *config.Perterbation, prefix string, index int, gridSystem *grid.System) *System {
  systemConfig := perterbation.SystemConfig

  system := new(System)
  newPrefix := SetNameAndGetPrefix(system, prefix, index)

  system.GridId = gridSystem.Id

  system.Features, perterbation = RollDetails(systemConfig.SystemFeaturesRolls, systemConfig.WeightedInspirations, systemConfig.ExtraInspirations, perterbation)

  numberOfStarClusters := config.RollAll(systemConfig.SystemStarClustersRolls, perterbation)
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
