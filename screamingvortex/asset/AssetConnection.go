package asset

import (
  "screamingvortex/utilities"
  "screamingvortex/config"
)

type AssetConnection struct {
  Id int64 `sql:"id"`
  GridId int64 `sql:"grid_id"`
  AssetId int64 `sql:"asset_id"`
  StartId int64 `sql:"start_id"`
  EndId int64 `sql:"end_id"`
  sourceNode *AssetNode
  targetNode *AssetNode
  X int
  Y int
  asset *Asset
}

func (assetConnection *AssetConnection) TableName(assetConnectionType string) string {
  return "plan_asset_connection"
}

func (assetConnection *AssetConnection) GetId() *int64 {
  return &assetConnection.Id
}

func (assetConnection *AssetConnection) CreateReverse() *AssetConnection {
    return CreateAssetConnection(assetConnection.targetNode, assetConnection.sourceNode)
}

func CreateAssetConnection(sourceNode *AssetNode, targetNode *AssetNode) *AssetConnection {
  assetConnection := new(AssetConnection)
  assetConnection.sourceNode = sourceNode
  assetConnection.targetNode = targetNode
  assetConnection.X = targetNode.X
  assetConnection.Y = targetNode.Y

  return assetConnection
}

func (assetConnection *AssetConnection) Randomize(perterbation *config.Perterbation, typeId int64, prefix string, index int) {
  assetConnection.asset = RollAsset(perterbation, typeId, prefix, index)
}

func (assetConnection *AssetConnection) SaveParents(client utilities.ClientInterface) {
  client.Save(assetConnection.asset, "")
  assetConnection.AssetId = assetConnection.asset.Id
}

func (assetConnection *AssetConnection) SaveChildren(client utilities.ClientInterface) {
  assetConnection.asset.SaveChildren(client)
}
