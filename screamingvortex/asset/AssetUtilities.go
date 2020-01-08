package asset

import "strconv"

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
