package utilities

type SectorToSystemLink struct {
  Id int64 `sql:"id"`
  ParentId int64 `sql:"asset_sector_id"`
  ChildId int64 `sql:"asset_system_id"`
}

func (link *SectorToSystemLink) TableName(linkType string) string {
  return "plan_asset_sector_systems"
}

func (link *SectorToSystemLink) GetId() *int64 {
  return &link.Id
}

//

type StarClusterToDetailLink struct {
  Id int64 `sql:"id"`
  ParentId int64 `sql:"asset_star_cluster_id"`
  ChildId int64 `sql:"detail_id"`
}

func (link *StarClusterToDetailLink) TableName(linkType string) string {
  return "plan_asset_star_cluster_stars"
}

func (link *StarClusterToDetailLink) GetId() *int64 {
  return &link.Id
}

//

type SystemToDetailLink struct {
  Id int64 `sql:"id"`
  ParentId int64 `sql:"asset_system_id"`
  ChildId int64 `sql:"detail_id"`
}

func (link *SystemToDetailLink) TableName(linkType string) string {
  return "plan_asset_system_details"
}

func (link *SystemToDetailLink) GetId() *int64 {
  return &link.Id
}

//

type SystemToStarClusterLink struct {
  Id int64 `sql:"id"`
  ParentId int64 `sql:"asset_system_id"`
  ChildId int64 `sql:"asset_star_cluster_id"`
}

func (link *SystemToStarClusterLink) TableName(linkType string) string {
  return "plan_asset_system_star_clusters"
}

func (link *SystemToStarClusterLink) GetId() *int64 {
  return &link.Id
}

//

type SystemToRouteLink struct {
  Id int64 `sql:"id"`
  ParentId int64 `sql:"asset_system_id"`
  ChildId int64 `sql:"asset_route_id"`
}

func (link *SystemToRouteLink) TableName(linkType string) string {
  return "plan_asset_system_routes"
}

func (link *SystemToRouteLink) GetId() *int64 {
  return &link.Id
}

//

type RouteToTargetSystemLink struct {
  Id int64 `sql:"id"`
  ParentId int64 `sql:"asset_route_id"`
  ChildId int64 `sql:"asset_system_id"`
}

func (link *RouteToTargetSystemLink) TableName(linkType string) string {
  return "plan_asset_route_target_systems"
}

func (link *RouteToTargetSystemLink) GetId() *int64 {
  return &link.Id
}
