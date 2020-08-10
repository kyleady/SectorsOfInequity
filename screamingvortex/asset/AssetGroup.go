package asset

import "screamingvortex/config"
import "screamingvortex/utilities"

type AssetGroup struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  TypeId int64 `sql:"type_id"`
  Assets []*Asset
}

func (assetGroup *AssetGroup) TableName(assetGroupType string) string {
  return "plan_asset_group"
}

func (assetGroup *AssetGroup) GetId() *int64 {
  return &assetGroup.Id
}

func (assetGroup *AssetGroup) SaveTo(client utilities.ClientInterface) {
  client.Save(assetGroup, "")
  assetGroup.SaveChildren(client)
}

func (assetGroup *AssetGroup) SaveChildren(client utilities.ClientInterface) {
  client.SaveAll(&assetGroup.Assets, "")
  client.SaveMany2ManyLinks(assetGroup, &assetGroup.Assets, "", "", "assets", false)
  for _, asset := range assetGroup.Assets {
    asset.SaveChildren(client)
  }
}

func RollAssetGroups(configGroups []*config.GroupConfig, prefix string, perterbation *config.Perterbation) []*AssetGroup {
  assetGroups := []*AssetGroup{}
  for _, configGroup := range configGroups {
    assetGroup := RollAssetGroup(configGroup, prefix, perterbation)
    assetGroups = append(assetGroups, assetGroup)
  }

  return assetGroups
}

func RollAssetGroup(configGroup *config.GroupConfig, prefix string, perterbation *config.Perterbation) *AssetGroup {
  assetGroup := new(AssetGroup)
  assetGroup.TypeId = configGroup.TypeId
  assetGroup.Name = configGroup.Name
  assetGroup.Assets = RollAssets(perterbation, assetGroup.TypeId, prefix, configGroup.Count, configGroup.Extras)
  return assetGroup
}
