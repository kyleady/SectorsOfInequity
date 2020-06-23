package asset

import "screamingvortex/config"
import "screamingvortex/utilities"

type StarCluster struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  Stars []*Detail
  Zones []*Zone
}

func (starCluster *StarCluster) TableName(starClusterType string) string {
  return "plan_asset_star_cluster"
}

func (starCluster *StarCluster) GetId() *int64 {
  return &starCluster.Id
}

func (starCluster *StarCluster) GetType() string {
  return "Star Cluster"
}

func (starCluster *StarCluster) SetName(name string) {
  starCluster.Name = name
}

func (starCluster *StarCluster) SaveTo(client utilities.ClientInterface) {
  client.Save(starCluster, "")
  starCluster.SaveChildren(client)
}

func (starCluster *StarCluster) SaveChildren(client utilities.ClientInterface) {
  client.SaveAll(&starCluster.Stars, "")
  for _, star := range starCluster.Stars {
    client.Save(&utilities.StarClusterToDetailLink{ParentId: starCluster.Id, ChildId: star.Id}, "")
    star.SaveChildren(client)
  }

  client.SaveAll(&starCluster.Zones, "")
  for _, zone := range starCluster.Zones {
    client.Save(&utilities.StarClusterToZoneLink{ParentId: starCluster.Id, ChildId: zone.Id}, "")
    zone.SaveChildren(client)
  }
}

func RandomStarCluster(perterbation *config.Perterbation, prefix string, index int) *StarCluster {
  starClusterConfig := perterbation.StarClusterConfig

  starCluster := new(StarCluster)
  newPrefix := SetNameAndGetPrefix(starCluster, prefix, index)
  starCluster.Stars, perterbation = RollDetails(
    starClusterConfig.StarsRolls,
    starClusterConfig.WeightedStarTypes,
    starClusterConfig.ExtraStarTypeIds,
    perterbation,
  )

  starCluster.Zones = RandomZones(perterbation, newPrefix)

  return starCluster
}
