package asset

import (
  "screamingvortex/utilities"
  "screamingvortex/config"
  "math/rand"
  "fmt"
)

type AssetGrid struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  Nodes []*AssetNode
  Connections []*AssetConnection
  grid [][]*AssetNode
  Type string
}

func (assetGrid *AssetGrid) TableName(assetGridType string) string {
  return "plan_asset_grid"
}

func (assetGrid *AssetGrid) GetId() *int64 {
  return &assetGrid.Id
}

func (assetGrid *AssetGrid) SetName(name string) {
  assetGrid.Name = name
}

func (assetGrid *AssetGrid) GetType() string {
  return assetGrid.Type
}

func (assetGrid *AssetGrid) SaveTo(client utilities.ClientInterface) {
  client.Save(assetGrid, "")
  assetGrid.SaveChildren(client)
}

func (assetGrid *AssetGrid) SaveChildren(client utilities.ClientInterface) {
  for _, assetNode := range assetGrid.Nodes {
    assetNode.GridId = assetGrid.Id
    assetNode.SaveParents(client)
  }

  client.SaveAll(&assetGrid.Nodes, "")

  for _, assetNode := range assetGrid.Nodes {
    assetNode.SaveChildren(client)
  }

  for _, assetConnection := range assetGrid.Connections {
    assetConnection.StartId = assetConnection.sourceNode.Id
    assetConnection.EndId = assetConnection.targetNode.Id
    assetConnection.GridId = assetGrid.Id
    assetConnection.SaveParents(client)
  }

  client.SaveAll(&assetGrid.Connections, "")

  for _, assetConnection := range assetGrid.Connections {
    assetConnection.SaveChildren(client)
  }
}

func RollAssetGrids(perterbation *config.Perterbation, gridConfigs []*config.GridConfig, prefix string) []*AssetGrid {
  assetGrids := []*AssetGrid{}
  for _, gridConfig := range gridConfigs {
    numberOfGrids := config.RollAll(gridConfig.Count, perterbation)
    for i := 0; i < numberOfGrids; i++ {
      assetGrid := RollAssetGrid(perterbation, gridConfig, prefix, 1+len(assetGrids))
      assetGrids = append(assetGrids, assetGrid)
    }
  }

  return assetGrids
}

func RollAssetGrid(perterbation *config.Perterbation, gridConfig *config.GridConfig, prefix string, index int) *AssetGrid {
  assetGrid := new(AssetGrid)
  assetGrid.Type = gridConfig.Name
  newPrefix := SetNameAndGetPrefix(assetGrid, prefix, index)
  assetGrid.randomizeStructure(perterbation, gridConfig)
  assetGrid.randomizeAssets(perterbation, gridConfig, newPrefix)
  return assetGrid
}

func (assetGrid *AssetGrid) randomizeStructure(perterbation *config.Perterbation, gridConfig *config.GridConfig) {
  assetGrid.createGrid(gridConfig, perterbation)
  assetGrid.populateGrid(perterbation, gridConfig)
  assetGrid.connectNodes(gridConfig, perterbation)
  blobSizes := assetGrid.labelBlobsAndGetSizes()
  assetGrid.trimToLargestBlob(blobSizes)
  assetGrid.gridToList()
  assetGrid.smooth(gridConfig, perterbation)
}

func (assetGrid *AssetGrid) randomizeAssets(perterbation *config.Perterbation, gridConfig *config.GridConfig, prefix string) {
  regionNameToIds := make(map[string][]int64)
  for  _, weightedRegion := range gridConfig.WeightedRegions {
    regionNameToIds[weightedRegion.ValueName] = weightedRegion.Values
  }

  totalNodes := 0
  for i, assetNode := range assetGrid.Nodes {
    assetNode.Randomize(perterbation, regionNameToIds[assetNode.regionName], prefix, i+1)
    totalNodes = i
  }

  for i, assetConnection := range assetGrid.Connections {
    assetConnection.Randomize(perterbation, gridConfig.ConnectionTypeId, prefix, i+1+totalNodes)
  }
}

func (assetGrid *AssetGrid) createGrid(gridConfig *config.GridConfig, perterbation *config.Perterbation) {
  height := config.RollAll(gridConfig.Height, perterbation)
  width := config.RollAll(gridConfig.Width, perterbation)
  assetGrid.grid = make([][]*AssetNode, height)
  for i := range assetGrid.grid {
    assetGrid.grid[i] = make([]*AssetNode, width)
  }
}

func (assetGrid *AssetGrid) populateGrid(perterbation *config.Perterbation, gridConfig *config.GridConfig) {
  for i, row := range assetGrid.grid {
    for j := range row {
      assetGrid.createNode(i, j, perterbation, gridConfig)
    }
  }
}

func (assetGrid *AssetGrid) connectNodes(gridConfig *config.GridConfig, perterbation *config.Perterbation) {
  reach := config.RollAll(gridConfig.ConnectionRange, perterbation)
  connectionRate := gridConfig.GetConnectionFraction(perterbation)
  rangeRateMultiplier := gridConfig.GetRangeMultiplierFraction(perterbation)
  for i := range assetGrid.grid {
    for j := range assetGrid.grid[i] {

      // Attempt to connect the node @{i,j} to all nodes within reach
      for r_i := -reach; r_i <= reach; r_i++ {
        for r_j := -reach; r_j <= reach; r_j++ {
          if roll := perterbation.Rand.Float64(); roll < assetGrid.connectionChance(r_i, r_j, connectionRate, rangeRateMultiplier) {
            targetNode := assetGrid.getNode(i+r_i, j+r_j)
            assetGrid.grid[i][j].ConnectTo(targetNode)
          }
        }
      }

    }
  }
}

func (assetGrid *AssetGrid) labelBlobsAndGetSizes() []int {
  currentLabel := 0
  blobSizes := make([]int, 0)
  for i := range assetGrid.grid {
    for j := range assetGrid.grid[i] {
      blobSize := assetGrid.grid[i][j].LabelBlob(currentLabel)

      if blobSize != 0 {
        blobSizes = append(blobSizes, blobSize)
        currentLabel++
      }
    }
  }

  return blobSizes
}

func (assetGrid *AssetGrid) trimToLargestBlob(blobSizes []int) {
  largestBlobSize := -100
  largestBlobLabel := -100
  totalSize := 0
  for blobLabel, blobSize := range blobSizes {
    totalSize += blobSize
    if blobSize > largestBlobSize {
      largestBlobSize = blobSize
      largestBlobLabel = blobLabel
    }
  }

  for i := range assetGrid.grid {
    for j := range assetGrid.grid[i] {
      assetGrid.grid[i][j].VoidNonMatchingLabel(largestBlobLabel)
    }
  }
}

func (assetGrid *AssetGrid) createNode(i int, j int, perterbation *config.Perterbation,  gridConfig *config.GridConfig) {
  assetNode := new(AssetNode)
  populationRate := gridConfig.GetPopulationFraction(perterbation)
  if roll := perterbation.Rand.Float64(); roll < populationRate {
    assetNode.InitializeAt(i, j)
  } else {
    assetNode.SetToVoid()
  }

  assetGrid.grid[i][j] = assetNode
}



func (assetGrid *AssetGrid) connectionChance(r_i int, r_j int, connectionRate float64, rangeRateMultiplier float64) float64 {
  abs_i := r_i
  if r_i < 0 {
    abs_i = - r_i
  }

  abs_j := r_j
  if r_j < 0 {
    abs_j = - r_j
  }

  var reach int
  if abs_j > abs_i {
    reach = abs_j
  } else {
    reach = abs_i
  }

  chance := connectionRate
  for i := 0; i < reach - 1; i++ {
    chance *= rangeRateMultiplier
  }

  return chance
}

func (assetGrid *AssetGrid) getNode(i int, j int) *AssetNode {
  if i < 0 || i >= len(assetGrid.grid) {
    return nil
  }

  if j < 0 || j >= len(assetGrid.grid[i]) {
    return nil
  }

  return assetGrid.grid[i][j]
}

func (assetGrid *AssetGrid) gridToList() {
  for i := range assetGrid.grid {
    for j := range assetGrid.grid[i] {
      if assetGrid.grid[i][j].IsVoid() == false {
        assetGrid.grid[i][j].nodeIndex = len(assetGrid.Nodes)
        assetGrid.Nodes = append(assetGrid.Nodes, assetGrid.grid[i][j])
      }
    }
  }

  for _, node := range assetGrid.Nodes {
    for _, connection := range node.Connections {
      reverseConnectionAlreadyAdded := false
      for _, finalCutConnection := range assetGrid.Connections {
        if connection.sourceNode == finalCutConnection.targetNode && connection.targetNode == finalCutConnection.sourceNode {
          reverseConnectionAlreadyAdded = true
          break
        }
      }

      if reverseConnectionAlreadyAdded {
        continue
      }

      assetGrid.Connections = append(assetGrid.Connections, connection)
    }
  }
}

func (assetGrid *AssetGrid) getTwoDifferentNodes(listByRegion map[string][]int, perterbation *config.Perterbation) (*AssetNode, *AssetNode, int, int) {
  nodeListIndexA := perterbation.Rand.Intn(len(assetGrid.Nodes))
  var regionA string
  var nodeA *AssetNode
  var nodeIndexA int
  for regionName, nodeIndexList := range listByRegion {
    if nodeListIndexA < len(nodeIndexList) {
      nodeIndexA = nodeIndexList[nodeListIndexA]
      nodeA = assetGrid.Nodes[nodeIndexA]
      regionA = regionName
      break
    } else {
      nodeListIndexA -= len(nodeIndexList)
    }
  }

  nodeListIndexB := rand.Intn(len(assetGrid.Nodes) - len(listByRegion[regionA]))
  var nodeB *AssetNode
  var nodeIndexB int
  for regionName, nodeIndexList := range listByRegion {
    if regionName == regionA {
      continue
    }

    if nodeListIndexB < len(nodeIndexList) {
      nodeIndexB = nodeIndexList[nodeListIndexB]
      nodeB = assetGrid.Nodes[nodeIndexB]
      break
    } else {
      nodeListIndexB -= len(nodeIndexList)
    }
  }

  return nodeA, nodeB, nodeListIndexA, nodeListIndexB
}

func (assetGrid *AssetGrid) smooth(gridConfig *config.GridConfig, perterbation *config.Perterbation) {
  smoothingFactor := gridConfig.GetSmoothingFraction(perterbation)
  listByRegion := make(map[string][]int)
  if smoothingFactor >= 1 {
    listByRegion = assetGrid.genClumpedRegionNames(gridConfig, perterbation)
    smoothingFactor--
  } else {
    listByRegion = assetGrid.genScatteredRegionNames(gridConfig, perterbation)
  }

  assetGrid.smoothRegions(listByRegion, smoothingFactor, perterbation)
}

func (assetGrid AssetGrid) genScatteredRegionNames(gridConfig *config.GridConfig, perterbation *config.Perterbation) map[string][]int {
  listByRegion := make(map[string][]int)
  for nodeIndex, node := range assetGrid.Nodes {
    randRegion := config.RollWeightedValues(gridConfig.WeightedRegions, perterbation, []*config.Roll{})
    regionName := randRegion.ValueName
    node.regionName = regionName
    listByRegion[regionName] = append(listByRegion[regionName], nodeIndex)
  }

  return listByRegion
}

func (assetGrid AssetGrid) smoothRegions(listByRegion map[string][]int, smoothingFactor float64, perterbation *config.Perterbation) {
  if len(assetGrid.Nodes) <= 2 {
    return
  }

  if len(listByRegion) <= 1 {
    return
  }

  maxSwitches := int(float64(len(assetGrid.Nodes)) * smoothingFactor)
  for i := 0; i < maxSwitches; i++ {
    nodeA, nodeB, nodeListIndexA, nodeListIndexB := assetGrid.getTwoDifferentNodes(listByRegion, perterbation)
    regionA := nodeA.regionName
    regionB := nodeB.regionName
    smoothScoreBefore := 0
    smoothScoreAfter := 0
    for _, connection := range nodeA.Connections {
      if connection.targetNode.regionName == regionA {
        smoothScoreBefore += 1
      } else if connection.targetNode.regionName == regionB {
        if connection.targetNode != nodeB {
          smoothScoreAfter += 1
        }
      }
    }

    for _, connection := range nodeB.Connections {
      if connection.targetNode.regionName == regionB {
        smoothScoreBefore += 1
      } else if connection.targetNode.regionName == regionA {
        if connection.targetNode != nodeA {
          smoothScoreAfter += 1
        }
      }
    }

    if smoothScoreAfter >= smoothScoreBefore {
      placeHolder := listByRegion[regionA][nodeListIndexA]
      listByRegion[regionA][nodeListIndexA] = listByRegion[regionB][nodeListIndexB]
      listByRegion[regionB][nodeListIndexB] = placeHolder
      nodeA.regionName = regionB
      nodeB.regionName = regionA
    }
  }
}

func (assetGrid AssetGrid) getRandomUnsetNode(nodesSet int, perterbation *config.Perterbation) *AssetNode {
  randomIndex := perterbation.Rand.Intn(len(assetGrid.Nodes) - nodesSet)

  if nodesSet == 0 {
    return assetGrid.Nodes[randomIndex]
  } else {
    for _, node := range assetGrid.Nodes {
      if node.regionName == "" {
        if randomIndex == 0 {
          return node
        } else {
          randomIndex--
        }
      }
    }
  }

  panic(fmt.Sprintf("One node should always be returned. {nodesSet=%d, len(assetGrid.Nodes)=%d}", nodesSet, len(assetGrid.Nodes)))
}

func (assetGrid AssetGrid) genClumpedRegionNames(gridConfig *config.GridConfig, perterbation *config.Perterbation) map[string][]int {
  listByRegion := make(map[string][]int)

  //determine the number of nodes in each region
  regionFrequency := make(map[string]int)
  for range assetGrid.Nodes {
    randRegion := config.RollWeightedValues(gridConfig.WeightedRegions, perterbation, []*config.Roll{})
    regionFrequency[randRegion.ValueName]++
  }

  //create ordered list of regions, from most populous to least
  regionFrequencyCopy := make(map[string]int)
  for regionName, freq := range regionFrequency {
    regionFrequencyCopy[regionName] = freq
  }

  orderedRegionNamesByFrequency := make([]string, len(regionFrequency))
  for i, _ := range orderedRegionNamesByFrequency {
    maxFreq := -1
    maxRegionName := ""
    for regionName, freq := range regionFrequencyCopy {
      if freq > maxFreq {
        maxFreq = freq
        maxRegionName = regionName
      }
    }

    orderedRegionNamesByFrequency[i] = maxRegionName
    regionFrequencyCopy[maxRegionName] = -1
  }

  //loop through each region, from most populous to least
  nodesSet := 0
  for _, regionName := range orderedRegionNamesByFrequency {
    boundryNodes := []*AssetNode{}

    //randomly choose starting node
    selectedNode := assetGrid.getRandomUnsetNode(nodesSet, perterbation)
    selectedIndex := -1

    //mark a node as part of this region, for each node in the region
    for i := 0; i < regionFrequency[regionName]; i++ {
      if selectedNode.regionName != "" {
        fmt.Printf("Selected Node not unset!\n%+v\n\n", selectedNode)
      }

      //mark the currently selected node as in the current region
      selectedNode.regionName = regionName
      listByRegion[regionName] = append(listByRegion[regionName], selectedNode.nodeIndex)
      nodesSet++

      //remove all occurences of the previously selected node from list of boundry nodes
      if selectedIndex >= 0 {
        for removeIndex := 0; removeIndex < len(boundryNodes); removeIndex++ {
          if boundryNodes[removeIndex] == selectedNode {
            boundryNodes = append(boundryNodes[:removeIndex], boundryNodes[removeIndex+1:]...)
            removeIndex--
          }
        }
      }

      //recalculate boundry nodes
      for _, connection := range selectedNode.Connections {
        boundryNode := connection.targetNode
        if boundryNode.regionName == "" {
          boundryNodes = append(boundryNodes, boundryNode)
        }
      }

      //check if you have set all nodes
      if nodesSet == len(assetGrid.Nodes) {
        break
      }

      //select a random node
      if len(boundryNodes) == 0 {
        //no boundry nodes are left
        //randomly choose new starting node somewhere else
        selectedNode = assetGrid.getRandomUnsetNode(nodesSet, perterbation)
        selectedIndex = -1
      } else {
        //randomly choose a boundry node
        selectedIndex = perterbation.Rand.Intn(len(boundryNodes))
        selectedNode = boundryNodes[selectedIndex]
      }
    }
  }

  return listByRegion
}
