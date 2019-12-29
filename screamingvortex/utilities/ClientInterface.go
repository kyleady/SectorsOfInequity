package utilities

type ClientInterface interface {
  Open()
  Close()
  Fetch(obj SQLInterface, tableType string, id int64)
  FetchAll(asInterface interface{}, tableType string, whereClause string, whereValues ...interface{})
  Save(obj SQLInterface, tableType string)
  SaveAll(asInterface interface{}, tableType string)
}
