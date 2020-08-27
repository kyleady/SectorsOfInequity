package asset

import "testing"
import "math/rand"
import "math"
import "strconv"

func createTestNodeObjs(depth int, branches int) ([][]*AssetNode) {
  nodes := make([][]*AssetNode, depth)

  baseNode := new(AssetNode)
  width := 1
  for i := range nodes {
    nodes[i] = make([]*AssetNode, width)
    base_j := 0
    relative_j := 0
    for j := range nodes[i] {
      nodes[i][j] = new(AssetNode)
      nodes[i][j].InitializeAt(i, j)
      if baseNode != nil {
        nodes[i][j].ConnectTo(baseNode)
      }

      relative_j += 1
      if i == 0 || relative_j >= branches {
        relative_j = 0
        base_j += 1
        if i > 0 && base_j < len(nodes[i-1]) {
          baseNode = nodes[i-1][base_j]
        } else {
          baseNode = nodes[i][0]
        }
      }
    }

    width *= branches
  }

  return nodes
}

func showBlob(blob [][]*AssetNode) string {
  output := ""

  for _, row := range blob {
    for _, assetNode := range row {
      output += strconv.Itoa(len(assetNode.Connections)) + ","
    }

    output += "\n\n"
  }

  return output
}

func TestNodeInitializeAt(t *testing.T) {
  node := new(AssetNode)
  randX := rand.Intn(100)
  randY := rand.Intn(100)

  node.InitializeAt(randX, randY)

  if node.Label() != node.TheUnsetLabel() {
    t.Errorf("Node's label was not initialized to the UnsetLabel. Label: %d, UnsetLabel: %d", node.Label(), node.TheUnsetLabel())
  }

  if node.X != randX || node.Y != randY {
    t.Errorf("Node's location does not match the initilization coords. Node: {%d,%d}, Location: {%d, %d}", node.X, node.Y, randX, randY)
  }

  if len(node.Connections) != 0 {
    t.Errorf("Node has more than zero connections. Connections: %d", len(node.Connections))
  }
}

func TestNodeSetToVoid(t *testing.T) {
  node := new(AssetNode)
  node.SetToVoid()

  if !node.IsVoid() {
    t.Errorf("Node was set to void, but is not recognized as void.")
  }
}

func TestNodeConnectTo(t *testing.T) {
  nodeA := new(AssetNode)
  nodeB := new(AssetNode)
  xA := rand.Intn(100)
  yA := rand.Intn(100)
  xB := rand.Intn(100)
  yB := rand.Intn(100)

  nodeA.InitializeAt(xA, yA)
  nodeA.ConnectTo(nodeA)
  if len(nodeA.Connections) != 0 {
    t.Errorf("A node should not be allowed to connect to itself. Connections: %d", len(nodeA.Connections))
  }

  nodeA.InitializeAt(xA, yA)
  nodeA.ConnectTo(nil)
  if len(nodeA.Connections) != 0 {
    t.Errorf("A node should not be allowed to connect to itself. Connections: %d", len(nodeA.Connections))
  }

  nodeA.SetToVoid()
  nodeB.InitializeAt(xB, yB)
  nodeA.ConnectTo(nodeB)
  if len(nodeA.Connections) != 0 || len(nodeB.Connections) != 0 {
    t.Errorf("Void nodes should not be allowed to connect. Void Connections: %d, NodeB Connections: %d", len(nodeA.Connections), len(nodeB.Connections))
  }

  nodeA.InitializeAt(xA, yA)
  nodeB.SetToVoid()
  nodeA.ConnectTo(nodeB)
  if len(nodeA.Connections) != 0 || len(nodeB.Connections) != 0 {
    t.Errorf("A node should not be allowed to connect to void nodes. NodeA Connections: %d, Void Connections: %d", len(nodeA.Connections), len(nodeB.Connections))
  }

  nodeA.InitializeAt(xA, yA)
  nodeB.InitializeAt(xB, yB)
  nodeA.ConnectTo(nodeB)
  if len(nodeA.Connections) != 1 || len(nodeB.Connections) != 1 {
    t.Errorf("Connecting two nodes should add a connection to both nodes. NodeA Connections: %d, NodeB Connections: %d", len(nodeA.Connections), len(nodeB.Connections))
  }

  if nodeA.Connections[0].X != xB || nodeB.Connections[0].X != xA || nodeA.Connections[0].Y != yB || nodeB.Connections[0].Y != yA {
    t.Errorf("The connections between two nodes should point to each other. NodeA TargetCoords: {%d,%d}, NodeB TargetCoords: {%d,%d}, NodeA Location: {%d,%d}, NodeB Location{%d,%d}", nodeA.Connections[0].X, nodeA.Connections[0].Y, nodeB.Connections[0].X, nodeB.Connections[0].Y, nodeA.X, nodeA.Y, nodeB.X, nodeB.Y)
  }

  nodeA.InitializeAt(xA, yA)
  nodeB.InitializeAt(xB, yB)
  nodeA.ConnectTo(nodeB)
  nodeB.ConnectTo(nodeA)
  nodeA.ConnectTo(nodeB)
  if len(nodeA.Connections) != 1 || len(nodeB.Connections) != 1 {
    t.Errorf("Connecting two nodes repeatedly should not add additional connections. NodeA Connections: %d, NodeB Connections: %d", len(nodeA.Connections), len(nodeB.Connections))
  }
}

func TestNodeLabelBlob(t *testing.T) {
  numberOfBlobs := 10
  maxDepth := 5
  maxBranches := 5

  blobs := make([][][]*AssetNode, numberOfBlobs)
  labels := make([]int, numberOfBlobs)
  depths := make([]int, numberOfBlobs)
  branches := make([]int, numberOfBlobs)
  expectedSizes := make([]int, numberOfBlobs)
  sizesByLength := make([]int, numberOfBlobs)
  randomNode := make([]*AssetNode, numberOfBlobs)

  for i := range blobs {
    labels[i] = i
    depths[i] = rand.Intn(maxDepth-2) + 2
    branches[i] = rand.Intn(maxBranches-2) + 2
    blobs[i] = createTestNodeObjs(depths[i], branches[i])
    expectedSizes[i] = int( ( math.Pow(float64(branches[i]), float64(depths[i])) - 1 ) / ( float64(branches[i]) - 1 ) )
    sizesByLength[i] = 0
    for j := range blobs[i] {
      sizesByLength[i] += len(blobs[i][j])
    }
  }

  for i := range blobs {
    randomNode[i] = blobs[i][rand.Intn(depths[i]-1)+1][rand.Intn(branches[i])]
    sizeOfBlob := randomNode[i].LabelBlob(labels[i])
    if sizeOfBlob != expectedSizes[i] {
      t.Errorf("The blob size was incorrectly calculated. Output: %d, Expected: %d, ByLength: %d\n" + showBlob(blobs[i]), sizeOfBlob, expectedSizes[i], sizesByLength[i])
    }

    errorFound := false
    for _, row := range blobs[i] {
      for _, node := range row {
        label := node.Label()
        if label != labels[i] {
          t.Errorf("The node does not have the proper label. Node: %d, Expected: %d", label, labels[i])
          errorFound = true
          break
        }
      }

      if errorFound {
        break
      }
    }
  }
}

func TestNodeVoidNonMatchingLabel(t *testing.T) {
  maxDepth := 5
  maxBranches := 5
  depth := rand.Intn(maxDepth - 3) + 3
  branches := rand.Intn(maxBranches - 3) + 3
  blob := createTestNodeObjs(depth, branches)
  label := rand.Intn(10)
  blob[rand.Intn(depth-1)+1][rand.Intn(branches-1)+1].LabelBlob(label)
  nodeThatShouldVoid := blob[rand.Intn(depth-1)+1][rand.Intn(branches-1)+1]
  nodeThatShouldNotVoid := blob[rand.Intn(depth-1)+1][rand.Intn(branches-1)+1]
  for nodeThatShouldVoid == nodeThatShouldNotVoid {
    nodeThatShouldNotVoid = blob[rand.Intn(depth-1)+1][rand.Intn(branches-1)+1]
  }

  connectionsBeforeRemoval := len(nodeThatShouldVoid.Connections)

  nodeThatShouldVoid.VoidNonMatchingLabel(label + 1)
  nodeThatShouldNotVoid.VoidNonMatchingLabel(label)

  if !nodeThatShouldVoid.IsVoid() {
    t.Errorf("The node's label is not void. Node: %d, Expected: %d", nodeThatShouldVoid.Label(), nodeThatShouldVoid.TheVoidLabel())
  }

  if len(nodeThatShouldVoid.Connections) != 0 {
    t.Errorf("The node's connections were not removed. Node: %d, Expected: 0", len(nodeThatShouldVoid.Connections))
  }

  if nodeThatShouldNotVoid.IsVoid() {
    t.Errorf("The node's label is void. Node: %d, Expected: %d", nodeThatShouldNotVoid.Label(), label)
  }

  if len(nodeThatShouldNotVoid.Connections) == 0 {
    t.Errorf("The node's connections were removed. Node: 0, Expected: %d", connectionsBeforeRemoval)
  }
}
