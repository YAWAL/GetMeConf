// Package entity contains database entities
package entity

//ConfigInterface is an interface for all config structures
type ConfigInterface interface {
	TabName() string
}

//PersistedData stores the information about all config types in database and is used during searching for a config by name and type
type PersistedData struct {
	ConfigType ConfigInterface
	IDField    string
}

//Mongodb is an random config example
type Mongodb struct {
	Domain  string `json:"domain"gorm:"type:varchar(100);primary_key"validate:"nonzero"`
	Mongodb bool   `json:"mongodb"gorm:"type:boolean"validate:"regexp=^(true|false)$"`
	Host    string `json:"host"gorm:"type:varchar(100)"validate:"nonzero"`
	Port    string `json:"port"gorm:"type:varchar(100)"validate:"nonzero"`
}

func (*Mongodb) TabName() string {
	return "mongodbs"
}

//Tsconfig is an random config example
type Tsconfig struct {
	Module    string `json:"module"gorm:"type:varchar(100);primary_key"validate:"nonzero"`
	Target    string `json:"target"gorm:"type:varchar(100)"validate:"nonzero"`
	SourceMap bool   `json:"sourceMap"gorm:"type:boolean"validate:"regexp=^(true|false)$"`
	Excluding int    `json:"excluding"gorm:"type:integer"`
}

func (*Tsconfig) TabName() string {
	return "tsconfigs"
}

//Tempconfig is an random config example
type Tempconfig struct {
	RestApiRoot    string `json:"restApiRoot"gorm:"type:varchar(100);primary_key"validate:"nonzero"`
	Host           string `json:"host"gorm:"type:varchar(100)"validate:"nonzero"`
	Port           string `json:"port"gorm:"type:varchar(100)"validate:"nonzero"`
	Remoting       string `json:"remoting"gorm:"type:varchar(100)"validate:"nonzero"`
	LegasyExplorer bool   `json:"legasyExplorer"gorm:"type:boolean"validate:"regexp=^(true|false)$"`
}

func (*Tempconfig) TabName() string {
	return "tempconfigs"
}
