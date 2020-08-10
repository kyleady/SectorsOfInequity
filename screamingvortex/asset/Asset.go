package asset

import "strconv"

import "screamingvortex/config"
import "screamingvortex/utilities"

type Asset struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  TypeId int64 `sql:"type_id"`
  Details []*Detail
  AssetGroups []*AssetGroup
  //Grids []*Grid
  Type *config.ConfigType
}

func (asset *Asset) TableName(assetType string) string {
  return "plan_asset"
}

func (asset *Asset) GetId() *int64 {
  return &asset.Id
}

func (asset *Asset) GetType() string {
  return asset.Type.Name
}

func (asset *Asset) SetName(name string) {
  asset.Name = name
}

func (asset *Asset) SetNameAndGetPrefix(prefix string, index int) string {
  var idNumber string
  indexAsAString := strconv.Itoa(index)
  if prefix != "" {
    idNumber = prefix + "-" + indexAsAString
  } else {
    idNumber = indexAsAString
  }

  asset.SetName(asset.GetType() + " " + idNumber)
  return idNumber
}

func (asset *Asset) SaveTo(client utilities.ClientInterface) {
  client.Save(asset, "")
  asset.SaveChildren(client)
}

func (asset *Asset) SaveChildren(client utilities.ClientInterface) {
  client.SaveAll(&asset.Details, "")
  client.SaveAll(&asset.AssetGroups, "")
  //client.SaveAll(&asset.Grids, "")
  client.SaveMany2ManyLinks(asset, &asset.Details, "", "", "details", false)
  client.SaveMany2ManyLinks(asset, &asset.AssetGroups, "", "", "asset_groups", false)
  //client.SaveMany2ManyLinks(asset, &asset.Grids, "", "", "grids", false)
  for _, detail := range asset.Details {
    detail.SaveChildren(client)
  }

  for _, assetGroup := range asset.AssetGroups {
    assetGroup.SaveChildren(client)
  }

  //for _, grid := range asset.Grids {
  //  grid.SaveChildren(client)
  //}
}

func RollAssets(perterbation *config.Perterbation, typeId int64, prefix string, countRolls []*config.Roll, extraInspirations []*config.InspirationExtra) []*Asset {
  assets := []*Asset{}
  assetsToAdd := config.RollAll(countRolls, perterbation)
  assetsPreviouslyAdded := 0
  i := 1
  for i = i; i <= assetsPreviouslyAdded + assetsToAdd; i++ {
    assets = append(assets, RollAsset(perterbation, typeId, prefix, i))
  }

  for _, extraInspiration := range extraInspirations {
    assetsToAdd := config.RollAll(extraInspiration.CountRolls, perterbation)
    assetsPreviouslyAdded := len(assets)
    if assetsToAdd <= 0 {
      continue
    }

    for i = i; i <= assetsPreviouslyAdded + assetsToAdd; i++ {
      assets = append(assets, ExtraAsset(perterbation, typeId, prefix, i, extraInspiration.InspirationTables))
    }
  }

  return assets
}

func RollAsset(perterbation *config.Perterbation, typeId int64, prefix string, index int) *Asset {
  assetConfig := perterbation.GetConfig(typeId)
  return ExtraAsset(perterbation, typeId, prefix, index, assetConfig.InspirationTables)
}

func ExtraAsset(perterbation *config.Perterbation, typeId int64, prefix string, index int, extraInspirationTables []*config.InspirationTable) *Asset {
  return newAsset(perterbation, typeId, prefix, index, extraInspirationTables)
}

func newAsset(perterbation *config.Perterbation, typeId int64, prefix string, index int, inspirationTables []*config.InspirationTable) *Asset {
  asset := new(Asset)
  details, newPerterbation := RollDetails(inspirationTables, perterbation)
  assetConfig := newPerterbation.GetConfig(typeId)
  asset.Type = newPerterbation.Manager.GetConfigType(typeId)
  asset.TypeId = typeId
  newPrefix := asset.SetNameAndGetPrefix(prefix, index)
  asset.Details = details
  asset.AssetGroups = RollAssetGroups(assetConfig.GroupConfigs, newPrefix, newPerterbation)
  //asset.Grids = RollGrids(assetConfig.GridConfigs, prefix, newPerterbation)
  return asset
}
