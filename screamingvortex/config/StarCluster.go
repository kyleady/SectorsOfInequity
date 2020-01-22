package config

import "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"

type StarCluster struct {
  StarsRolls []*Roll
}

func (system *StarCluster) TableName(systemType string) string {
  return "plan_config_star_cluster"
}

func (system *StarCluster) GetId() *int64 {
  panic("GetId() not implemented. Config should not be editted.")
}

func CreateEmptyStarClusterConfig() *StarCluster {
  starCluster := new(StarCluster)
  starCluster.StarsRolls = make([]*Roll, 0)
  return starCluster
}

func (starCluster *StarCluster) AddPerterbation(perterbation *StarCluster) *StarCluster {
  newStarCluster := new(StarCluster)
  newStarCluster.StarsRolls = append(starCluster.StarsRolls, perterbation.StarsRolls...)
  return newStarCluster
}

func LoadStarClusterConfigFrom(client utilities.ClientInterface, id int64) *StarCluster {
  starCluster := new(StarCluster)
  FetchAllRolls(client, &starCluster.StarsRolls, id, starCluster.TableName(""), "star_count")
  return starCluster
}
