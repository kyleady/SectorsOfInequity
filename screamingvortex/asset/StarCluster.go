package asset

import "github.com/kyleady/SectorsOfInequity/screamingvortex/config"
import "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"

type StarCluster struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  ParentId int64 `sql:"parent_id"`
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
}


func RandomStarCluster(perterbation *config.Perterbation, prefix string, index int) *StarCluster {
  starCluster := new(StarCluster)
  //newPrefix :=
  SetNameAndGetPrefix(starCluster, prefix, index)

  return starCluster
}
