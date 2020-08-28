package asset

import "strconv"

type AssetInterface interface {
  SetName(string)
  GetType() string
}

func SetNameAndGetPrefix(assetInterface AssetInterface, prefix string, index int) string {
  newPrefix := ""
  indexAsAString := strconv.Itoa(index)
  if prefix != "" {
    newPrefix = prefix + "-" + indexAsAString
  } else {
    newPrefix = indexAsAString
  }

  assetInterface.SetName(assetInterface.GetType() + " " + newPrefix)
  return newPrefix
}
