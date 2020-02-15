package utilities

import "fmt"
import "reflect"
import "regexp"
import "strings"
import "database/sql"


type ClientMock struct {
  Environment string
  Region string
  Resource string
  Secret string

  isConnected bool
  tmpDB map[string][]SQLInterface
}

func (client *ClientMock) Close() {
  client.isConnected = false
}

func (client *ClientMock) Open() {
  client.isConnected = true
}

func (client *ClientMock) checkConnection() {
  if !client.isConnected {
    panic("Client is not connected")
  }
}

func (client *ClientMock) AddTable_(tablename string) {
  if client.tmpDB == nil {
    client.tmpDB = make(map[string][]SQLInterface)
  }
  client.tmpDB[tablename] = make([]SQLInterface, 0)
}

func (client *ClientMock) Fetch(obj SQLInterface, tableType string, id int64) {
  client.checkConnection()
  fetchQuery(obj, tableType, "id = ?")
  _, _, addresses := listFields(obj, true)
  values, _, _ := listFields(client.tmpDB[obj.TableName(tableType)][id], true)
  assignAll(values, addresses)
}

func assignAll(values []interface{}, addresses []interface{}) {
  for i := 0; i < len(addresses); i++ {
    switch address := addresses[i].(type) {
    case *int:
      *address = values[i].(int)
    case *int64:
      *address = values[i].(int64)
    case *string:
      *address = values[i].(string)
    case *float64:
      *address = values[i].(float64)
    case *sql.NullString:
      *address = values[i].(sql.NullString)
    case * sql.NullInt64:
      *address = values[i].(sql.NullInt64)
    default:
      panic("Unknown type")
    }
  }
}

func (client *ClientMock) FetchAll(asInterface interface{}, tableType string, whereClause string, whereValues ...interface{}) {
  client.checkConnection()
  asSlice := reflect.ValueOf(asInterface).Elem()
  asSlice.Set(reflect.MakeSlice(asSlice.Type(), 0, 0))

  generatedObj := reflect.New(asSlice.Type().Elem())


  var obj SQLInterface
  arrayOf := "?"
  if generatedObj.Elem().Kind() == reflect.Ptr {
    arrayOf = "pointers"
    generatedObj = reflect.New(asSlice.Type().Elem().Elem())
    obj = generatedObj.Interface().(SQLInterface)
  } else {
    arrayOf = "structs"
    obj = generatedObj.Interface().(SQLInterface)
  }
  tmpTable := client.tmpDB[obj.TableName(tableType)]
  fetchQuery(obj, tableType, whereClause)

  re := regexp.MustCompile(`\s*(\w+)\s*=\s*\?`)
  conditionMatches := re.FindAllStringSubmatch(whereClause, -1)
  for i := 0; i < len(tmpTable); i++ {
    objFromTable := tmpTable[i]
    values, names, _ := listFields(objFromTable, true)
    conditionMet := true
    for j := 0; j < len(names); j++ {
      for k := 0; k < len(conditionMatches); k++ {
        if names[j] == conditionMatches[k][1] {
          if values[j] != whereValues[k] {
            conditionMet = false
          }
          break
        }
      }
    }

    if conditionMet {
      generatedObj := reflect.New(asSlice.Type().Elem())
      if arrayOf == "pointers" {
        generatedObj = reflect.New(asSlice.Type().Elem().Elem())
      }
      _, _, addresses := listFields(generatedObj.Interface(), true)
      values, _, _ := listFields(objFromTable, true)
      assignAll(values, addresses)
      if arrayOf == "pointers" {
        asSlice.Set(reflect.Append(asSlice, generatedObj))
      } else {
        asSlice.Set(reflect.Append(asSlice, generatedObj.Elem()))
      }
    }
  }
}

func (client *ClientMock) FetchMany(asInterface interface{}, parentId int64, parentTableName string, childTableName string, valueName string, childType string, reverseAccess bool) {
  childTableNameWithoutAppName := strings.Replace(childTableName, "plan_", "", 1)
  parentTableNameWithoutAppName := strings.Replace(parentTableName, "plan_", "", 1)

  if reverseAccess {
    tmpVariable := childTableNameWithoutAppName
    childTableNameWithoutAppName = parentTableNameWithoutAppName
    parentTableNameWithoutAppName = tmpVariable
  }

  whereClause := fmt.Sprintf("id IN (SELECT %s_id FROM %s_%s WHERE %s_id = ?)",
                              childTableNameWithoutAppName,
                              parentTableName,
                              valueName,
                              parentTableNameWithoutAppName,
                            )

  client.FetchAll(asInterface, childType, whereClause, parentId)
}

func (client *ClientMock) Update(obj SQLInterface, tableType string) {
  client.checkConnection()
  updateQuery(obj, tableType)
  tableName := obj.TableName(tableType)
  client.tmpDB[tableName][*obj.GetId()] = obj
}

func (client *ClientMock) Save(obj SQLInterface, tableType string) {
  client.checkConnection()
  saveQuery(obj, tableType, 1)
  tableName := obj.TableName(tableType)
  insert_id := len(client.tmpDB[tableName])
  *(obj.GetId()) = int64(insert_id)
  client.tmpDB[tableName] = append(client.tmpDB[tableName], obj)
}

func (client *ClientMock) SaveAll(asInterface interface{}, tableType string) {
  client.checkConnection()
  objs := reflect.ValueOf(asInterface).Elem()
  if objs.Len() <= 0 {
    return
  }

  objValue := objs.Index(0)
  var obj SQLInterface
  if objValue.Kind() == reflect.Ptr {
    obj = objValue.Interface().(SQLInterface)
  } else {
    obj = objValue.Addr().Interface().(SQLInterface)
  }
  saveQuery(obj, tableType, objs.Len())
  tableName := obj.TableName(tableType)
  insert_id := len(client.tmpDB[tableName])
  for i := 0; i < objs.Len(); i++ {
    objValue := objs.Index(i)
    var obj SQLInterface
    if objValue.Kind() == reflect.Ptr {
      obj = objValue.Interface().(SQLInterface)
    } else {
      obj = objValue.Addr().Interface().(SQLInterface)
    }
    *(obj.GetId()) = int64(insert_id + i)
    client.tmpDB[tableName] = append(client.tmpDB[tableName], obj)
  }
}

func (client *ClientMock) Print(tablename string) {
  fmt.Printf("<=%s=>\n", tablename)
  for i := 0; i < len(client.tmpDB[tablename]); i++ {
    fmt.Printf("<%d>\n", i)
    fmt.Printf("%v+\n\n", client.tmpDB[tablename][i])
  }
}

func (client *ClientMock) Delete(obj SQLInterface, tableType string) {
  panic("Not yet implemented")
}
