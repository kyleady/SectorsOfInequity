package utilities

type ClientInterface interface {
  Open()
  Close()
  Fetch(obj SQLInterface, id int64)
  FetchAll(asInterface interface{}, whereClause string, whereValues ...interface{})
  Save(obj SQLInterface)
  SaveAll(asInterface interface{})
}
