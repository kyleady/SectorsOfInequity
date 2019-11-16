package utilities

import "reflect"
import "regexp"
import "fmt"

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

func (client *ClientMock) Fetch(obj SQLInterface, id int64) {
  client.checkConnection()
  fetchQuery(obj, "id = ?")
  _, _, addresses := listFields(obj, true)
  values, _, _ := listFields(client.tmpDB[obj.TableName()][id], true)
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
    default:
      panic("Unknown type")
    }
  }
}

func (client *ClientMock) FetchAll(asInterface interface{}, whereClause string, whereValues ...interface{}) {
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
  tmpTable := client.tmpDB[obj.TableName()]
  fetchQuery(obj, whereClause)

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

func (client *ClientMock) Update(obj SQLInterface) {
  panic("Not yet implemented")
}

func (client *ClientMock) Save(obj SQLInterface) {
  client.checkConnection()
  saveQuery(obj, 1)
  insert_id := len(client.tmpDB[obj.TableName()])
  *(obj.GetId()) = int64(insert_id)
  client.tmpDB[obj.TableName()] = append(client.tmpDB[obj.TableName()], obj)
}

func (client *ClientMock) SaveAll(asInterface interface{}) {
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
  saveQuery(obj, objs.Len())
  insert_id := len(client.tmpDB[obj.TableName()])
  for i := 0; i < objs.Len(); i++ {
    objValue := objs.Index(i)
    var obj SQLInterface
    if objValue.Kind() == reflect.Ptr {
      obj = objValue.Interface().(SQLInterface)
    } else {
      obj = objValue.Addr().Interface().(SQLInterface)
    }
    *(obj.GetId()) = int64(insert_id + i)
    client.tmpDB[obj.TableName()] = append(client.tmpDB[obj.TableName()], obj)
  }
}

func (client *ClientMock) Print(tablename string) {
  fmt.Printf("<=%s=>\n", tablename)
  for i := 0; i < len(client.tmpDB[tablename]); i++ {
    fmt.Printf("<%d>\n", i)
    fmt.Printf("%v+\n\n", client.tmpDB[tablename][i])
  }
}

func (client *ClientMock) Delete(obj SQLInterface) {
  panic("Not yet implemented")
}
