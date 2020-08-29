package asset

import "screamingvortex/config"
import "screamingvortex/utilities"

type AssetGroup struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
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

func RollAssetGroups(address []*config.InspirationKey, prefix string, perterbation *config.Perterbation) []*AssetGroup {
  assetGroups := []*AssetGroup{}
  configGroupKeys := perterbation.GetGroupConfigKeys(address)
  for _, configGroupKey := range configGroupKeys {
    assetGroup := RollAssetGroup(append(
        address,
        configGroupKey,
    ), prefix, perterbation)
    assetGroups = append(assetGroups, assetGroup)
  }

  return assetGroups
}

func RollAssetGroup(address []*config.InspirationKey, prefix string, perterbation *config.Perterbation) *AssetGroup {
  configGroup := perterbation.GetGroupConfig(address)
  newPerterbation := perterbation
  for _, perterbationId := range configGroup.PerterbationIds {
    newPerterbation = newPerterbation.AddPerterbation(perterbationId)
  }

  assetGroup := new(AssetGroup)
  assetGroup.Name = configGroup.Name
  assetGroup.Assets = RollAssets(newPerterbation, configGroup.Types, prefix, configGroup.Count, address)
  return assetGroup
}
