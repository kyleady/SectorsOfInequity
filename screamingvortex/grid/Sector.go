package grid

import (
  "github.com/kyleady/SectorsOfInequity/screamingvortex/config"
  "math/rand"
  "time"
)

type Sector struct {
  Systems []System
  grid [][]System
  config *config.GridConfig
  rand *rand.Rand
}

func (sector *Sector) Randomize(gridConfig *config.GridConfig) {
  sector.config = gridConfig

  source := rand.NewSource(time.Now().UnixNano())
  sector.rand = rand.New(source)

  sector.createGrid()
  sector.populateGrid()
  sector.connectSystems()
  blobSizes := sector.labelBlobsAndGetSizes()
  sector.trimToLargestBlob(blobSizes)
  sector.saveSystems()
}

func (sector *Sector) createGrid() {
  sector.grid = make([][]System, sector.config.Height)
  for i := range sector.grid {
    sector.grid[i] = make([]System, sector.config.Width)
  }
}

func (sector *Sector) populateGrid() {
  for i, row := range sector.grid {
    for j := range row {
      sector.createSystem(i, j)
    }
  }
}

func (sector *Sector) createSystem(i int, j int) {
  system := new(System)
  if roll := sector.rand.Float64(); roll < sector.config.PopulationRate {
    system.InitializeAt(i, j)
  } else {
    system.SetToVoidSpace()
  }

  sector.grid[i][j] = *system
}

func (sector *Sector) connectSystems() {
  reach := sector.config.ConnectionRange

  for i := range sector.grid {
    for j := range sector.grid[i] {

      // Attempt to connect the system @{i,j} to all systems within reach
      for r_i := -reach; r_i <= reach; r_i++ {
        for r_j := -reach; r_j <= reach; r_j++ {
          if roll := sector.rand.Float64(); roll < sector.connectionChance(r_i, r_j) {
            targetSystem := sector.getSystem(i+r_i, j+r_j)
            sector.grid[i][j].ConnectTo(targetSystem)
          }
        }
      }

    }
  }
}

func (sector *Sector) connectionChance(r_i int, r_j int) float64 {
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

  chance := sector.config.ConnectionRate
  for i := 0; i < reach - 1; i++ {
    chance *= sector.config.RangeRateMultiplier
  }

  return chance
}

func (sector *Sector) getSystem(i int, j int) *System {
  if i < 0 || i >= sector.config.Height {
    return nil
  }

  if j < 0 || j >= sector.config.Width {
    return nil
  }

  return &sector.grid[i][j]
}

func (sector *Sector) labelBlobsAndGetSizes() []int {
  currentLabel := 0
  blobSizes := make([]int, 0)
  for i := range sector.grid {
    for j := range sector.grid[i] {
      blobSize := sector.grid[i][j].LabelBlob(currentLabel)

      if blobSize != 0 {
        blobSizes = append(blobSizes, blobSize)
        currentLabel++
      }
    }
  }

  return blobSizes
}

func (sector *Sector) trimToLargestBlob(blobSizes []int) {
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

  for i := range sector.grid {
    for j := range sector.grid[i] {
      sector.grid[i][j].VoidNonMatchingLabel(largestBlobLabel)
    }
  }
}

func (sector *Sector) saveSystems() {
  for i := range sector.grid {
    for j := range sector.grid {
      if sector.grid[i][j].IsVoidSpace() == false {
        sector.Systems = append(sector.Systems, sector.grid[i][j])
      }
    }
  }
}
