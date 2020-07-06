package asset

import "strconv"

import "screamingvortex/config"

type AssetInterface interface {
  SetName(name string)
  GetType() string
}

func SetNameAndGetPrefix(assetObj AssetInterface, prefix string, index int) string {
  var idNumber string
  indexAsAString := strconv.Itoa(index)
  if prefix != "" {
    idNumber = prefix + "-" + indexAsAString
  } else {
    idNumber = indexAsAString
  }

  assetObj.SetName(assetObj.GetType() + " " + idNumber)

  return idNumber
}

func RollDetails(rollableDetailCount []*config.Roll, weightedInspirations []*config.WeightedValue, extraInspirations []*config.WeightedValue, perterbation *config.Perterbation) ([]*Detail, *config.Perterbation) {
  numberOfDetails := config.RollAll(rollableDetailCount, perterbation)
  var details []*Detail
  for i := 1; i <= numberOfDetails; i++ {
    detail, newPerterbation := RollDetail(weightedInspirations, perterbation)

    if detail != nil {
      details = append(details, detail)
      perterbation = newPerterbation
    }
  }

  extraInspirations = config.ModifyExtraInspirations(extraInspirations, weightedInspirations)
  for _, extraInspiration := range extraInspirations {
    detailsToAdd := config.RollAll(extraInspiration.Weights, perterbation)
    for i := 0; i < detailsToAdd; i++ {
      detail, newPerterbation := NewDetail(extraInspiration.Values, perterbation)
      if detail != nil {
        details = append(details, detail)
        perterbation = newPerterbation
      }
    }
  }

  return details, perterbation
}

func RollAssetInspirations(rollableAssetCount []*config.Roll, extraAssetTypes []*config.WeightedValue, modifiers []*config.WeightedValue, perterbation *config.Perterbation) [][]int64 {
  assetInspirations := [][]int64{}
  rRand := perterbation.Rand
  numberOfRandomAssets := config.RollAll(rollableAssetCount, perterbation)
  shuffledExtraIds := config.ExtraInspirationsToShuffledExtraIds(extraAssetTypes, modifiers, perterbation)
  numberOfExtraAssets := len(shuffledExtraIds)
  numberOfAssets := numberOfRandomAssets + numberOfExtraAssets
  numberOfRandomAssetsCreated := 0
  numberOfExtraAssetsCreated := 0
  for i := 1; i <= numberOfAssets; i++ {
    if numberOfExtraAssets - numberOfExtraAssetsCreated > 0 && rRand.Intn(numberOfAssets - numberOfExtraAssetsCreated - numberOfRandomAssetsCreated) < numberOfExtraAssets - numberOfExtraAssetsCreated {
      extraInspirationIds := shuffledExtraIds[numberOfExtraAssetsCreated]
      assetInspirations = append(assetInspirations, extraInspirationIds)
      numberOfExtraAssetsCreated++
    } else {
      assetInspirations = append(assetInspirations, nil)
      numberOfRandomAssetsCreated++
    }
  }

  return assetInspirations
}
