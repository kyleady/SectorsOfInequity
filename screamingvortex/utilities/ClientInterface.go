package utilities

type ClientInterface interface {
  Open()
  Close()
  Fetch(obj SQLInterface, tableType string, id int64)
  FetchAll(asInterface interface{}, tableType string, whereClause string, whereValues ...interface{})
  FetchMany(asInterface interface{}, parentId int64, parentTableName string, childTableName string, valueName string, childType string, reverseAccess bool)
  Save(obj SQLInterface, tableType string)
  SaveAll(asInterface interface{}, tableType string)
  Update(obj SQLInterface, tableType string)
  FetchManyToManyChildIds(ids *[]int64, parendId int64, parentTableName string, childTableName string, valueName string, childType string, reverseAccess bool)
}
