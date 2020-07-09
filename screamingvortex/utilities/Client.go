package utilities

import (
    "os"
    "fmt"
    "strings"
    "reflect"
    "io/ioutil"
    "encoding/json"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/secretsmanager"
)

type Client struct {
  Environment string
  Region string
  Resource string
  Secret string
  Local string

  dbUser string
  dbPassword string
  dbConnection string
  dbHost string
  dbPort float64
  dbName string
  db *sql.DB
}

func (client *Client) Close() {
  client.db.Close()
}

func (client *Client) Open() {
  response := new(svcResponse)
  if client.Local != "" {
    response = client.mockedResponse()
  } else {
    response = client.secretResponse()
  }

  client.dbUser = response.Username
  client.dbPassword = response.Password
  client.dbConnection = "tcp"
  client.dbHost = response.Host
  client.dbPort = response.Port
  client.dbName = client.Resource

  connectionString := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",
    client.dbUser,
    client.dbPassword,
    client.dbConnection,
    client.dbHost,
    int(client.dbPort),
    client.dbName,
  )
  db, err := sql.Open("mysql", connectionString)

  if err != nil {
      panic(err.Error())
  }

  client.db = db
}

func (client *Client) Fetch(obj SQLInterface, tableType string, id int64) {
  query := fetchQuery(obj, tableType, "id = ?")
  rows, err := client.db.Query(query, id)
  if err != nil {
  	panic(err)
  }
  _, _, addresses := listFields(obj, true)
  for rows.Next() {
  	err := rows.Scan(addresses...)
  	if err != nil {
  		panic(err)
  	}
  }
}

func (client *Client) FetchAll(asInterface interface{}, tableType string, whereClause string, whereValues ...interface{}) {
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
  query := fetchQuery(obj, tableType, whereClause)
  rows, err := client.db.Query(query, whereValues...)
  if err != nil {
  	panic(err)
  }
  for rows.Next() {
    if arrayOf == "pointers" {
      generatedObj = reflect.New(asSlice.Type().Elem().Elem())
    } else if arrayOf == "structs" {
      generatedObj = reflect.New(asSlice.Type().Elem())
    }
    _, _, addresses := listFields(generatedObj.Interface(), true)
    err := rows.Scan(addresses...)
    if err != nil {
  		panic(err)
  	}
    if arrayOf == "pointers" {
      asSlice.Set(reflect.Append(asSlice, generatedObj))
    } else {
      asSlice.Set(reflect.Append(asSlice, generatedObj.Elem()))
    }
  }
}

func (client *Client) FetchManyToManyChildIds(ids *[]int64, parendId int64, parentTableName string, childTableName string, valueName string, childType string, reverseAccess bool) {
  query := manyQuery(
      parentTableName,
      childTableName,
      valueName,
      reverseAccess,
    )
  rows, err := client.db.Query(query, parendId)
  if err != nil {
  	panic(err)
  }

  for rows.Next() {
    id := new(int64)
    err := rows.Scan(id)
    if err != nil {
  		panic(err)
  	}

    *ids = append(*ids, *id)
  }
}

func (client *Client) FetchMany(asInterface interface{}, parentId int64, parentTableName string, childTableName string, valueName string, childType string, reverseAccess bool) {
  whereClause := fmt.Sprintf("id IN (%s)", manyQuery(
      parentTableName,
      childTableName,
      valueName,
      reverseAccess,
    ))

  client.FetchAll(asInterface, childType, whereClause, parentId)
}

func manyTableAndColumnNames(parentTableName string, childTableName string, valueName string, reverseAccess bool) (string, string, string) {
  childTableNameWithoutAppName := strings.Replace(childTableName, "plan_", "", 1)
  parentTableNameWithoutAppName := strings.Replace(parentTableName, "plan_", "", 1)

  tableBaseName := parentTableName
  if reverseAccess {
    tableBaseName = childTableName
  }

  parentPrefix := ""
  childPrefix := ""
  if parentTableName == childTableName {
    parentPrefix = "from_"
    childPrefix = "to_"
  }

  return fmt.Sprintf("%s_%s", tableBaseName, valueName), fmt.Sprintf("%s%s_id", parentPrefix, parentTableNameWithoutAppName), fmt.Sprintf("%s%s_id", childPrefix, childTableNameWithoutAppName)
}

func manyQuery(parentTableName string, childTableName string, valueName string, reverseAccess bool) string {
  manyTableName, parentIdName, childIdName := manyTableAndColumnNames(parentTableName, childTableName, valueName, reverseAccess)
  return fmt.Sprintf("SELECT %s FROM %s WHERE %s = ?",
                              childIdName,
                              manyTableName,
                              parentIdName,
                            )
}

func (client *Client) Update(obj SQLInterface, tableType string) {
  values, _, _ := listFields(obj, false)
  query := updateQuery(obj, tableType)

  args := append(values, *obj.GetId())
  _, err := client.db.Exec(query, args...)
  if err != nil {
  	panic(err)
  }
}

func updateQuery(obj SQLInterface, tableType string) string {
  _, names, _ := listFields(obj, false)
  return fmt.Sprintf("UPDATE %s SET %s%s WHERE id = ?;",
    obj.TableName(tableType),
    strings.Join(names, " = ?, "),
    " = ?",
  )
}

func (client *Client) Save(obj SQLInterface, tableType string) {
  query := saveQuery(obj, tableType, 1)
  values, _, _ := listFields(obj, false)
  result, err := client.db.Exec(query, values...)
  if err != nil {
  	panic(err)
  }
  insert_id, err := result.LastInsertId()
  if err != nil {
  	panic(err)
  }
  *(obj.GetId()) = insert_id
}

func (client *Client) SaveAll(asInterface interface{}, tableType string) {
  objs := reflect.ValueOf(asInterface).Elem()
  if objs.Len() <= 0 {
    return
  }

  var obj SQLInterface
  objValue := objs.Index(0)
  arrayOf := "?"
  if objValue.Kind() == reflect.Ptr {
    arrayOf = "pointers"
    obj = objValue.Interface().(SQLInterface)
  } else {
    arrayOf = "structs"
    obj = objValue.Addr().Interface().(SQLInterface)
  }
  query := saveQuery(obj, tableType, objs.Len())

  all_values := make([]interface{}, 0, objs.Len())
  all_ids := make([]*int64, 0, objs.Len())
  for i := 0; i < objs.Len(); i++ {
    objValue := objs.Index(i)
    var obj SQLInterface
    if arrayOf == "pointers" {
      obj = objValue.Interface().(SQLInterface)
    } else if arrayOf == "structs" {
      obj = objValue.Addr().Interface().(SQLInterface)
    }
    values, _, _ := listFields(obj, false)
    all_values = append(all_values, values...)
    all_ids = append(all_ids, obj.GetId())
  }

  result, err := client.db.Exec(query, all_values...)
  if err != nil {
  	panic(err)
  }
  insert_id, err := result.LastInsertId()
  if err != nil {
  	panic(err)
  }

  for i := 0; i < objs.Len(); i++ {
    *all_ids[i] = insert_id + int64(i)
  }
}

func (client *Client) SaveMany2ManyLinks(parentObj SQLInterface, childObjsInterface interface{}, parentTableType string, childTableType string, valueName string, reverseAccess bool) {
  childObjs := reflect.ValueOf(childObjsInterface).Elem()
  if childObjs.Len() <= 0 {
    return
  }

  var exampleChildObj SQLInterface
  exampleChildObjValue := childObjs.Index(0)
  arrayOf := "?"
  if exampleChildObjValue.Kind() == reflect.Ptr {
    arrayOf = "pointers"
    exampleChildObj = exampleChildObjValue.Interface().(SQLInterface)
  } else {
    arrayOf = "structs"
    exampleChildObj = exampleChildObjValue.Addr().Interface().(SQLInterface)
  }

  query := saveMany2ManyLinksQuery(parentObj.TableName(parentTableType), exampleChildObj.TableName(childTableType), valueName, reverseAccess, childObjs.Len())
  allParentAndChildIds := make([]interface{}, 0, 2 * childObjs.Len())
  for i := 0; i < childObjs.Len(); i++ {
    childObjValue := childObjs.Index(i)
    var childObj SQLInterface
    if arrayOf == "pointers" {
      childObj = childObjValue.Interface().(SQLInterface)
    } else if arrayOf == "structs" {
      childObj = childObjValue.Addr().Interface().(SQLInterface)
    }

    allParentAndChildIds = append(allParentAndChildIds, parentObj.GetId())
    allParentAndChildIds = append(allParentAndChildIds, childObj.GetId())
  }

  _, err := client.db.Exec(query, allParentAndChildIds...)
  if err != nil {
  	panic(err)
  }
}

func fetchQuery(obj SQLInterface, tableType string, whereClause string) string {
  _, names, _ := listFields(obj, true)
  return fmt.Sprintf("SELECT %s FROM %s WHERE %s;",
    strings.Join(names, ","), obj.TableName(tableType), whereClause)
}

func saveQuery(obj SQLInterface, tableType string, objCount int) string {
  _, names, _ := listFields(obj, false)
  questionMarks := make([]string, len(names), len(names))
  for i := 0; i < len(names); i++ {
    questionMarks[i] = "?"
  }
  parens := make([]string, objCount, objCount)
  for i := 0; i < objCount; i++ {
    parens[i] = fmt.Sprintf("(%s)", strings.Join(questionMarks, ","))
  }
  return fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;",
    obj.TableName(tableType),
    strings.Join(names, ","),
    strings.Join(parens, ","),
  )
}

func saveMany2ManyLinksQuery(parentTableName string, childTableName string, valueName string, reverseAccess bool, childCount int) string {
  parens := make([]string, childCount, childCount)
  for i := 0; i < childCount; i++ {
    parens[i] = "(?,?)"
  }

  manyTableName, parentIdName, childIdName := manyTableAndColumnNames(parentTableName, childTableName, valueName, reverseAccess)
  return fmt.Sprintf("INSERT INTO %s (%s,%s) VALUES %s;",
    manyTableName,
    parentIdName,
    childIdName,
    strings.Join(parens, ","),
  )
}

func (client *Client) Delete(obj SQLInterface, tableType string) {
  query := fmt.Sprintf("DELETE FROM %s WHERE id = ?;",
    obj.TableName(tableType),
  )
  _, err := client.db.Exec(query, *obj.GetId())
  if err != nil {
  	panic(err)
  }
}

func (client *Client) secretResponse() *svcResponse {
  awsSession, sessionErr := session.NewSession(&aws.Config{
  	Region: aws.String(client.Region),
  })
  if sessionErr != nil {
    panic(sessionErr.Error())
  }

  secret_name := fmt.Sprintf("%s/%s/%s",
    client.Environment,
    client.Resource,
    client.Secret,
  )
  input := &secretsmanager.GetSecretValueInput{
    SecretId: aws.String(secret_name),
  }

  svc := secretsmanager.New(awsSession)
  result, svcError := svc.GetSecretValue(input)
  if svcError != nil {
    panic(svcError.Error())
  }

  response := new(svcResponse)
  jsonErr := json.Unmarshal([]byte(*result.SecretString), response)
  if jsonErr != nil {
    panic(jsonErr)
  }

  return response
}

func (client *Client) mockedResponse() *svcResponse {
  mockedResponse, err := os.Open(client.Local)
  defer mockedResponse.Close()
  if err != nil {
    panic(err)
  }

  byteResponse, _ := ioutil.ReadAll(mockedResponse)
  response := new(svcResponse)
  json.Unmarshal(byteResponse, response)
  return response
}

func listFields (obj interface{}, includeId bool) ([]interface{}, []string, []interface{}) {
  rValue := reflect.ValueOf(obj)
  rType := reflect.TypeOf(obj)
  if rType.Kind() == reflect.Ptr {
    rValue = rValue.Elem()
    rType = rType.Elem()
  }
  if rType.Kind() == reflect.Ptr {
    rValue = rValue.Elem()
    rType = rType.Elem()
  }

  fieldCount := 0
  fieldCount = rValue.NumField()
  addresses := make([]interface{}, fieldCount)
  values := make([]interface{}, fieldCount)
  names := make([]string, fieldCount)
  included := 0
  for i := 0; i < fieldCount; i++ {
    sqlname := rType.Field(i).Tag.Get("sql")
    if !includeId && sqlname == "id" {
      continue
    }

    if len(sqlname) > 0 {
      field := rValue.Field(i)
      values[included] = field.Interface()
      names[included] = sqlname
      addresses[included] = field.Addr().Interface()

      included++
    }
  }

  return values[:included], names[:included], addresses[:included]
}

type svcResponse struct {
  Username string `json:"username"`
  Engine string `json:"engine"`
  Host string `json:"host"`
  Password string `json:"password"`
  Port float64 `json:"port"`
  DBInstanceIdentifier string `json:"dbInstanceIdentifier"`
}
