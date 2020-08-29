package asset

import (
  "screamingvortex/utilities"
  "screamingvortex/config"
)

type AssetNode struct {
  Id int64 `sql:"id"`
  GridId int64 `sql:"grid_id"`
  AssetId int64 `sql:"asset_id"`
  X int `sql:"x"`
  Y int `sql:"y"`
  Connections []*AssetConnection
  blobLabel int
  nodeIndex int
  regionName string
  asset *Asset
}

func (assetNode *AssetNode) TableName(systemType string) string {
  return "plan_asset_node"
}

func (assetNode *AssetNode) GetId() *int64 {
  return &assetNode.Id
}

func (assetNode *AssetNode) TheVoidLabel() int {
  return -2
}

func (assetNode *AssetNode) TheUnsetLabel() int {
  return -1
}

func (assetNode *AssetNode) Label() int {
  return assetNode.blobLabel
}

func (assetNode *AssetNode) IsVoid() bool {
  return assetNode.blobLabel == assetNode.TheVoidLabel()
}

func (assetNode *AssetNode) IsUnset() bool {
  return assetNode.blobLabel == assetNode.TheUnsetLabel()
}

func (assetNode *AssetNode) InitializeAt(i int, j int) {
  assetNode.X = i
  assetNode.Y = j
  assetNode.Connections = make([]*AssetConnection, 0)
  assetNode.blobLabel = assetNode.TheUnsetLabel()
  assetNode.regionName = ""
}

func (assetNode *AssetNode) SetToVoid() {
  assetNode.Connections = make([]*AssetConnection, 0)
  assetNode.blobLabel = assetNode.TheVoidLabel()
}

func (assetNode *AssetNode) ConnectTo(targetNode *AssetNode) {
  if targetNode == nil || assetNode == targetNode {
    return
  }

  if assetNode.IsVoid() || targetNode.IsVoid() {
    return
  }

  for _, connection := range assetNode.Connections {
    if connection.targetNode == targetNode {
      return
    }
  }

  newConnection := CreateAssetConnection(assetNode, targetNode)
  assetNode.Connections = append(assetNode.Connections, newConnection)
  targetNode.Connections = append(targetNode.Connections, newConnection.CreateReverse())
}

func (assetNode *AssetNode) LabelBlob(label int) int {
  if assetNode.blobLabel == assetNode.TheVoidLabel() || assetNode.blobLabel >= 0 {
    return 0
  }

  assetNode.blobLabel = label
  blobSize := 1
  for _, connection := range assetNode.Connections {
    blobSize += connection.targetNode.LabelBlob(label)
  }

  return blobSize
}

func (assetNode *AssetNode) VoidNonMatchingLabel(label int) {
  if assetNode.blobLabel != label {
    assetNode.SetToVoid()
  }
}

func (assetNode *AssetNode) Randomize(perterbation *config.Perterbation, regionConfigIds []int64, prefix string, index int) {
  newPerterbation := perterbation
  weightedTypes := []*config.WeightedValue{}
  for _, regionConfigId := range regionConfigIds {
    regionConfig := config.FetchRegionConfig(newPerterbation.Manager, regionConfigId)
    weightedTypes = config.StackWeightedValues(weightedTypes, regionConfig.Types)
    for _, regionPerterbationId := range regionConfig.PerterbationIds {
      newPerterbation = newPerterbation.AddPerterbation(regionPerterbationId)
    }
  }

  typeId := config.RollWeightedValues(weightedTypes, perterbation, nil).Values[0]
  assetNode.asset = RollAsset(newPerterbation, typeId, prefix, index)
}

func (assetNode *AssetNode) SaveParents(client utilities.ClientInterface) {
  client.Save(assetNode.asset, "")
  assetNode.AssetId = assetNode.asset.Id
}

func (assetNode *AssetNode) SaveChildren(client utilities.ClientInterface) {
  assetNode.asset.SaveChildren(client)
}
