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
