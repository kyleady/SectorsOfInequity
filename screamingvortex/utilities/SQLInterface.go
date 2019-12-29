package utilities

type SQLInterface interface {
  GetId() *int64
  TableName(string) string
}
