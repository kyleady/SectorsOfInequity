package asset

import "testing"
import "time"
import "math/rand"
import "screamingvortex/config"
import "screamingvortex/utilities"

func createTestGridConfigObj() *config.GridConfig {
  weightedRegions := []*config.WeightedValue{
    &config.WeightedValue{
      Weights: []*config.Roll{config.CreateConstantRoll(3)},
      Values: []int64{1},
      ValueName: "RImperium",
    },
    &config.WeightedValue{
      Weights: []*config.Roll{config.CreateConstantRoll(2)},
      Values: []int64{2},
      ValueName: "RWarp",
    },
    &config.WeightedValue{
      Weights: []*config.Roll{config.CreateConstantRoll(1)},
      Values: []int64{3},
      ValueName: "RTyranids",
    },
  }

  weightedConnecitonTypes := []*config.WeightedValue{
    &config.WeightedValue{
      Weights: []*config.Roll{config.CreateConstantRoll(1)},
      Values: []int64{1},
      ValueName: "1",
    },
  }

  return &config.GridConfig{
    Id: 2,
    Name: "Test Grid",
    WeightedRegions: weightedRegions,
    ConnectionTypes: weightedConnecitonTypes,
    Count: []*config.Roll{config.CreateConstantRoll(1)},
    Height: []*config.Roll{config.CreateConstantRoll(20)},
    Width: []*config.Roll{config.CreateConstantRoll(20)},
    ConnectionRange: []*config.Roll{config.CreateConstantRoll(5)},
    PopulationPercent: []*config.Roll{config.CreateConstantRoll(50)},
    ConnectionPercent: []*config.Roll{config.CreateConstantRoll(50)},
    RangeMultiplierPercent: []*config.Roll{config.CreateConstantRoll(25)},
    SmoothingPercent: []*config.Roll{config.CreateConstantRoll(1000)},
    PopulationDenominator: 100,
    ConnectionDenominator: 100,
    RangeMultiplierDenominator: 100,
    SmoothingDenominator: 100,
  }
}

func TestGridRandomize(t *testing.T) {
  client := &utilities.ClientMock{}
  client.Open()
  defer client.Close()

  gridConfig := createTestGridConfigObj()

  assetGrid := new(AssetGrid)
  timeNow := time.Now()
  randSource := rand.NewSource(timeNow.UnixNano())
  rRand := rand.New(randSource)

  perterbation := config.CreateEmptyPerterbation(client, rRand)
  assetGrid.randomizeStructure(perterbation, gridConfig)

  theSurvivingLabel := -100
  nodeCount := 0
  for _, row := range assetGrid.grid {
    for _, node := range row {
      if node.Label() == node.TheUnsetLabel() {
        t.Errorf("Grid Node {%d, %d} has not been properly labeled.", node.X, node.Y)
      }

      if node.Label() >= 0 {
        nodeCount++
        if theSurvivingLabel >= 0 && theSurvivingLabel != node.Label() {
          t.Errorf("The assetGrid has more than one surviving label. Label: %d, Expected: %d", node.Label(), theSurvivingLabel)
        } else if theSurvivingLabel < 0 {
          theSurvivingLabel = node.Label()
        }
      }

      if len(node.Connections) == 0 && !node.IsVoid() {
        t.Errorf("Grid Node is not connected to other nodes. This is possible but improbable if there is only one node in the blob. Label: %d", theSurvivingLabel)
      }

      if len(node.Connections) > 0 && node.IsVoid() {
        t.Errorf("Void is still connected to other nodes.")
      }

      for _, connection := range node.Connections {
        if connection.targetNode.Label() != theSurvivingLabel {
          t.Errorf("Connection is connected to node outside the main label. TargetNode: %+v", connection.targetNode)
        }
      }

      if node.regionName == "" && !node.IsVoid() {
        t.Errorf("Grid Node {%d, %d} has not been properly given a RegionName.\n%+v", node.X, node.Y, node)
      }
    }
  }

  if theSurvivingLabel < 0 {
    t.Errorf("No blobs were found. Possible but improbable. Each and every node would need to be empty space. Label: %d", theSurvivingLabel)
  }

  if nodeCount != len(assetGrid.Nodes) {
    t.Errorf("The number of surviving nodes on the grid does not equal the recorded surviving nodes in the assetGrid. %d != %d", nodeCount, len(assetGrid.Nodes))
  }
}
