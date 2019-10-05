package utilities

import (
    "fmt"
    "strings"
    "reflect"
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
  client.populateSecrets()

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

func (client *Client) Fetch(obj SQLInterface, id int64) {
  _, names, addresses := listFields(obj)
  query := fmt.Sprintf("SELECT %s FROM %s WHERE id = ?;",
    strings.Join(names, ","), obj.TableName())
  rows, err := client.db.Query(query, id)
  if err != nil {
  	panic(err)
  }
  for rows.Next() {
  	err := rows.Scan(addresses...)
  	if err != nil {
  		panic(err)
  	}
  }
}

func (client *Client) Update(obj SQLInterface) {
  values, names, _ := listFields(obj)
  query := fmt.Sprintf("UPDATE %s SET %s%s WHERE id = ?;",
    obj.TableName(),
    strings.Join(names, " = ?, "),
    " = ?",
  )

  args := append(values, *obj.GetId())
  _, err := client.db.Exec(query, args...)
  if err != nil {
  	panic(err)
  }
}

func (client *Client) Save(obj SQLInterface) {
  values, names, _ := listFields(obj)
  questionMarks := make([]string, len(names))
  for i := 0; i < len(names); i++ {
    questionMarks[i] = "?"
  }
  query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id;",
    obj.TableName(),
    strings.Join(names, ","),
    strings.Join(questionMarks, ","),
  )
  rows, err := client.db.Query(query, values...)
  if err != nil {
  	panic(err)
  }
  for rows.Next() {
  	err := rows.Scan(obj.GetId())
  	if err != nil {
  		panic(err)
  	}
  }
}

func (client *Client) Delete(obj SQLInterface) {
  query := fmt.Sprintf("DELETE FROM %s WHERE id = ?;",
    obj.TableName(),
  )
  _, err := client.db.Exec(query, *obj.GetId())
  if err != nil {
  	panic(err)
  }
}

func (client *Client) populateSecrets() {
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

  client.dbUser = response.Username
  client.dbPassword = response.Password
  client.dbConnection = "tcp"
  client.dbHost = response.Host
  client.dbPort = response.Port
  client.dbName = client.Resource
}

func listFields (obj interface{}) ([]interface{}, []string, []interface{}) {
  rValue := reflect.ValueOf(obj).Elem()
  rType := reflect.TypeOf(obj).Elem()
  fieldCount := rValue.NumField()
  addresses := make([]interface{}, fieldCount)
  values := make([]interface{}, fieldCount)
  names := make([]string, fieldCount)
  for i := 0; i < fieldCount; i++ {
    field := rValue.Field(i)
    values[i] = field
    names[i] = rType.Field(i).Tag.Get("sql")
    addresses[i] = field.Addr().Interface()
  }

  return values, names, addresses
}

type svcResponse struct {
  Username string `json:"username"`
  Engine string `json:"engine"`
  Host string `json:"host"`
  Password string `json:"password"`
  Port float64 `json:"port"`
  DBInstanceIdentifier string `json:"dbInstanceIdentifier"`
}
