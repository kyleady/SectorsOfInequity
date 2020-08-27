package asset

import "testing"

func createConnectionTestObjs() (*AssetConnection, *AssetNode, *AssetNode) {
  xA := 3
  yA := 4
  xB := 0
  yB := 8

  nodeA := new(AssetNode)
  nodeA.InitializeAt(xA, yA)
  nodeB := new(AssetNode)
  nodeB.InitializeAt(xB, yB)

  connection := CreateAssetConnection(nodeA, nodeB)

  return connection, nodeA, nodeB
}

func TestConnectionInitializeFromNodes(t *testing.T) {
  _, nodeA, nodeB := createConnectionTestObjs()

  connection := CreateAssetConnection(nodeA, nodeB)

  if connection.X != nodeB.X || connection.Y != nodeB.Y {
    t.Errorf("Connection Target Coords did not match Target Node Coords. Connection: {%d,%d}, Node: {%d, %d}",
      connection.X,
      connection.Y,
      nodeB.X,
      nodeB.Y,
    )
  }
}

func TestConnectionSourceNode(t *testing.T) {
  connection, nodeA, _ := createConnectionTestObjs()

  sourceNode := connection.sourceNode

  if sourceNode != nodeA {
    t.Errorf("Connection sourceNode did not return the private Source Node.")
  }
}

func TestConnectionTargetNode(t *testing.T) {
  connection, _, nodeB := createConnectionTestObjs()

  targetNode := connection.targetNode

  if targetNode != nodeB {
    t.Errorf("Connection targetNode did not return the private Target Node.")
  }
}

func TestConnectionCreateReverse(t *testing.T) {
  connection, _, _ := createConnectionTestObjs()

  inverseConnection := connection.CreateReverse()

  if inverseConnection.sourceNode.X != connection.X || inverseConnection.sourceNode.Y != connection.Y {
    t.Errorf("Inverse Source Coords did not match original Target Coords. Inverse: {%d,%d}, Original: {%d, %d}",
      inverseConnection.sourceNode.X,
      inverseConnection.sourceNode.Y,
      connection.X,
      connection.Y,
    )
  }

  if inverseConnection.sourceNode != connection.targetNode {
    t.Errorf("Inverse sourceNode did not equal the original targetNode.")
  }

  if inverseConnection.targetNode != connection.sourceNode {
    t.Errorf("Inverse targetNode did not equal the original sourceNode.")
  }
}
