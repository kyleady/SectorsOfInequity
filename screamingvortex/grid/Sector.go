package grid

import (
  "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"
  "github.com/kyleady/SectorsOfInequity/screamingvortex/config"
  "math/rand"
  "time"
  "fmt"
)

type Sector struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  Systems []*System
  grid [][]*System
  config *config.Grid
  rand *rand.Rand
}

func (sector *Sector) TableName(sectorType string) string {
  return "plan_grid_sector"
}

func (sector *Sector) GetId() *int64 {
  return &sector.Id
}

func (sector *Sector) SaveTo(client utilities.ClientInterface) {
  client.Save(sector, "")
  sector.saveSystems(client)
  sector.saveSystemRoutes(client)
}

func (sector *Sector) saveSystems(client utilities.ClientInterface) {
  for _, system := range sector.Systems {
    system.SectorId = sector.Id
  }
  client.SaveAll(&sector.Systems, "")
}

func (sector *Sector) saveSystemRoutes(client utilities.ClientInterface) {
  for _, system := range sector.Systems {
    for _, route := range system.Routes {
      route.StartId = system.Id
      route.EndId = -1
      for _, endSystem := range sector.Systems {
        if route.targetSystem == endSystem {
          route.EndId = endSystem.Id
          break
        }
      }
    }
    client.SaveAll(&system.Routes, "")
  }
}

func LoadSectorFrom(client utilities.ClientInterface, id int64) *Sector {
  sector := &Sector{}
  client.Fetch(sector, "", id)
  client.FetchAll(&sector.Systems, "", "sector_id = ?", sector.Id)
  for system_i := range sector.Systems {
    system_ptr := sector.Systems[system_i]
    client.FetchAll(&system_ptr.Routes, "", "start_id = ?", system_ptr.Id)
    for route_i := range system_ptr.Routes {
      route_ptr := system_ptr.Routes[route_i]
      route_ptr.sourceSystem = system_ptr
      for endSystem_i := range sector.Systems {
        endSystem_ptr := sector.Systems[endSystem_i]
        if route_ptr.EndId == endSystem_ptr.Id {
          route_ptr.targetSystem = endSystem_ptr
          route_ptr.X = endSystem_ptr.X
          route_ptr.Y = endSystem_ptr.Y
          break
        }
      }
    }
  }
  return sector
}

func (sector *Sector) Randomize(gridConfig *config.Grid, client utilities.ClientInterface, job *utilities.Job) {
  sector.config = gridConfig
  t := time.Now()
  sector.Name = gridConfig.Name + t.Format("_20060102150405")

  source := rand.NewSource(t.UnixNano())
  sector.rand = rand.New(source)


  sector.createGrid()
  sector.populateGrid()
  sector.connectSystems()
  blobSizes := sector.labelBlobsAndGetSizes()
  sector.trimToLargestBlob(blobSizes)
  sector.gridToList()
  job.Step(client)
  smoothingFactor := sector.config.SmoothingFactor
  listByRegion := make(map[int64][]int)
  if smoothingFactor >= 1 {
    listByRegion = sector.genClumpedRegionIds()
    smoothingFactor--
  } else {
    listByRegion = sector.genScatteredRegionIds()
  }
  job.Step(client)

  sector.smoothRegionIds(listByRegion, smoothingFactor)
  job.Step(client)
}

func (sector *Sector) createGrid() {
  sector.grid = make([][]*System, sector.config.Height)
  for i := range sector.grid {
    sector.grid[i] = make([]*System, sector.config.Width)
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

  sector.grid[i][j] = system
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

  return sector.grid[i][j]
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

func (sector *Sector) gridToList() {
  sector.Systems = make([]*System, 0)
  for i := range sector.grid {
    for j := range sector.grid[i] {
      if sector.grid[i][j].IsVoidSpace() == false {
        sector.grid[i][j].systemIndex = len(sector.Systems)
        sector.Systems = append(sector.Systems, sector.grid[i][j])
      }
    }
  }
}

func (sector *Sector) getTwoDifferentSystems(listByRegion map[int64][]int) (*System, *System, int, int) {
  systemListIndexA := rand.Intn(len(sector.Systems))
  var regionA int64
  var systemA *System
  var systemIndexA int
  for regionId, systemIndexList := range listByRegion {
    if systemListIndexA < len(systemIndexList) {
      systemIndexA = systemIndexList[systemListIndexA]
      systemA = sector.Systems[systemIndexA]
      regionA = regionId
      break
    } else {
      systemListIndexA -= len(systemIndexList)
    }
  }

  systemListIndexB := rand.Intn(len(sector.Systems) - len(listByRegion[regionA]))
  var systemB *System
  var systemIndexB int
  for regionId, systemIndexList := range listByRegion {
    if regionId == regionA {
      continue
    }

    if systemListIndexB < len(systemIndexList) {
      systemIndexB = systemIndexList[systemListIndexB]
      systemB = sector.Systems[systemIndexB]
      break
    } else {
      systemListIndexB -= len(systemIndexList)
    }
  }

  return systemA, systemB, systemListIndexA, systemListIndexB
}

func (sector *Sector) genScatteredRegionIds() map[int64][]int {
  listByRegion := make(map[int64][]int)
  for systemIndex, system := range sector.Systems {
    randRegion := config.RollWeightedValues(sector.config.WeightedRegions, sector.rand)
    system.RegionId = randRegion
    listByRegion[randRegion] = append(listByRegion[randRegion], systemIndex)
  }

  return listByRegion
}

func (sector *Sector) smoothRegionIds(listByRegion map[int64][]int, smoothingFactor float64) {
  if len(sector.Systems) <= 2 {
    return
  }

  if len(listByRegion) <= 1 {
    return
  }

  maxSwitches := int(float64(len(sector.Systems)) * sector.config.SmoothingFactor)
  for i := 0; i < maxSwitches; i++ {
    systemA, systemB, systemListIndexA, systemListIndexB := sector.getTwoDifferentSystems(listByRegion)
    regionA := systemA.RegionId
    regionB := systemB.RegionId
    vBefore := 0
    vAfter := 0
    for _, route := range systemA.Routes {
      if route.TargetSystem().RegionId == regionA {
        vBefore += 1
      } else if route.TargetSystem().RegionId == regionB {
        if route.TargetSystem() != systemB {
          vAfter += 1
        }
      }
    }

    for _, route := range systemB.Routes {
      if route.TargetSystem().RegionId == regionB {
        vBefore += 1
      } else if route.TargetSystem().RegionId == regionA {
        if route.TargetSystem() != systemA {
          vAfter += 1
        }
      }
    }

    if vAfter >= vBefore {
      placeHolder := listByRegion[regionA][systemListIndexA]
      listByRegion[regionA][systemListIndexA] = listByRegion[regionB][systemListIndexB]
      listByRegion[regionB][systemListIndexB] = placeHolder
      systemA.RegionId = regionB
      systemB.RegionId = regionA
    }
  }
}

func (sector *Sector) getRandomUnsetSystem(systemsSet int) *System {
  randomIndex := sector.rand.Intn(len(sector.Systems) - systemsSet)

  if systemsSet == 0 {
    return sector.Systems[randomIndex]
  } else {
    for _, system := range sector.Systems {
      if system.RegionId == int64(system.TheUnsetLabel()) {
        if randomIndex == 0 {
          return system
        } else {
          randomIndex--
        }
      }
    }
  }

  panic(fmt.Sprintf("One system should always be returned. {systemsSet=%d, len(sector.Systems)=%d}", systemsSet, len(sector.Systems)))
}

func (sector *Sector) genClumpedRegionIds() map[int64][]int {
  listByRegion := make(map[int64][]int)
  //determine the number of systems in each region
  regionFrequency := make(map[int64]int)
  for range sector.Systems {
    randRegion := config.RollWeightedValues(sector.config.WeightedRegions, sector.rand)
    regionFrequency[randRegion]++
  }

  //create ordered list of regions, from most populous to least
  regionFrequencyCopy := make(map[int64]int)
  for regionId, freq := range regionFrequency {
    regionFrequencyCopy[regionId] = freq
  }


  orderedRegionIdByFrequency := make([]int64, len(regionFrequency), len(regionFrequency))
  currentIndex := 0
  for range regionFrequency {
    maxFreq := -1
    maxRegionId := int64(-1)
    for regionId, freq := range regionFrequencyCopy {
      if freq > maxFreq {
        maxFreq = freq
        maxRegionId = regionId
      }
    }

    orderedRegionIdByFrequency[currentIndex] = maxRegionId
    regionFrequencyCopy[maxRegionId] = -1
    currentIndex++
  }

  //loop through each region, from most populous to least
  systemsSet := 0
  adjacentSystems := []*System{}
  for _, regionId := range orderedRegionIdByFrequency {
    //reset adjacentSystems
    adjacentSystems = []*System{}

    //randomly choose starting system
    selectedSystem := sector.getRandomUnsetSystem(systemsSet)
    selectedIndex := -1

    //mark a system as part of this region, for each system in the region
    for i := 0; i < regionFrequency[regionId]; i++ {
      //mark that system as in the current region
      selectedSystem.RegionId = regionId
      listByRegion[regionId] = append(listByRegion[regionId], selectedSystem.systemIndex)
      systemsSet++

      //remove all occurences of that system from list of adjacent systems
      if selectedIndex >= 0 {
        for adjacentIndex := 0; adjacentIndex < len(adjacentSystems); adjacentIndex++ {
          if adjacentSystems[adjacentIndex] == selectedSystem {
            adjacentSystems = append(adjacentSystems[:adjacentIndex], adjacentSystems[adjacentIndex+1:]...)
            adjacentIndex--;
          }
        }
      }

      //add all newly adjacent systems that are not yet in a region
      for _, route := range selectedSystem.Routes {
        adjacentSystem := route.TargetSystem()
        if adjacentSystem.RegionId == int64(adjacentSystem.TheUnsetLabel()) {
          adjacentSystems = append(adjacentSystems, adjacentSystem)
        }
      }

      //check if you have run out of adjacent systems
      if systemsSet < len(sector.Systems) {
        if len(adjacentSystems) == 0 {
          //randomly choose new starting system
          selectedSystem = sector.getRandomUnsetSystem(systemsSet)
          selectedIndex = -1

          //mark that system as part of the current region
        } else {
          //randomly choose an adjacent system
          selectedIndex = sector.rand.Intn(len(adjacentSystems))
          selectedSystem = adjacentSystems[selectedIndex]
        }
      }
    }
  }

  return listByRegion
}
