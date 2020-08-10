package config

type ConfigType struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
}

func (configType *ConfigType) TableName(configTypeType string) string {
  return "plan_config_name"
}

func (configType *ConfigType) GetId() *int64 {
  return &configType.Id
}

func FetchConfigType(manager *ConfigManager, configTypeId int64) *ConfigType {
  configType := new(ConfigType)
  manager.Client.Fetch(configType, "", configTypeId)
  return configType
}
