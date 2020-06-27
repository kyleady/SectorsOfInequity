package config

type StarCluster struct {
  StarsRolls []*Roll
  WeightedStarTypes []*WeightedValue
  ExtraStarTypeIds []int64
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
  starCluster.WeightedStarTypes = make([]*WeightedValue, 0)
  starCluster.ExtraStarTypeIds = make([]int64, 0)
  return starCluster
}

func (starCluster *StarCluster) AddPerterbation(perterbation *StarCluster) *StarCluster {
  newStarCluster := new(StarCluster)
  newStarCluster.StarsRolls = append(starCluster.StarsRolls, perterbation.StarsRolls...)
  newStarCluster.WeightedStarTypes = StackWeightedValues(starCluster.WeightedStarTypes, perterbation.WeightedStarTypes)
  newStarCluster.ExtraStarTypeIds = append(starCluster.ExtraStarTypeIds, perterbation.ExtraStarTypeIds...)
  return newStarCluster
}

func FetchStarClusterConfig(manager *ConfigManager, id int64) *StarCluster {
  starCluster := new(StarCluster)
  starCluster.StarsRolls = FetchManyRolls(manager, id, starCluster.TableName(""), "star_count")
  starCluster.WeightedStarTypes = FetchManyWeightedInspirations(manager, id, starCluster.TableName(""), "star_inspirations")
  starCluster.ExtraStarTypeIds = FetchManyInspirationIds(manager, id, starCluster.TableName(""), "star_extra")
  return starCluster
}
