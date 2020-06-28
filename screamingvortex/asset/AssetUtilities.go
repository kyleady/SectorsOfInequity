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
  numberOfDetails := config.RollAll(rollableDetailCount, perterbation.Rand)
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
    detailsToAdd := config.RollAll(extraInspiration.Weights, perterbation.Rand)
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

func RollAssets(rollableAssetCount []*config.Roll, newPrefix string, perterbation *config.Perterbation, assetGenerator func(*config.Perterbation, string, int) interface{}) interface{} {
  assets := []interface{}{}
  numberOfAssets := config.RollAll(rollableAssetCount, perterbation.Rand)
  for i := 1; i <= numberOfAssets; i++ {
    asset := assetGenerator(perterbation, newPrefix, i)
    assets = append(assets, asset)
  }

  return assets
}
