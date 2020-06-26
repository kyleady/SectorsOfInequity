package utilities

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
