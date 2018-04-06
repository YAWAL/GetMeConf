// Package entitie contains database entities
package entitie

//Mongodb is an random config example
type Mongodb struct {
	Domain  string `json:"domain" validate:"nonzero"`
	Mongodb bool   `json:"mongodb" validate:"regexp=^(true|false)$`
	Host    string `json:"host" validate:"nonzero"`
	Port    string `json:"port" validate:"nonzero"`
}

//Tsconfig is an random config example
type Tsconfig struct {
	Module    string `json:"module" validate:"nonzero"`
	Target    string `json:"target" validate:"nonzero"`
	SourceMap bool   `json:"sourceMap" validate:"regexp=^(true|false)$`
	Excluding int    `json:"excluding" validate:"nonzero"`
}

//Tempconfig is an random config example
type Tempconfig struct {
	RestApiRoot    string `json:"restApiRoot" validate:"nonzero"`
	Host           string `json:"host" validate:"nonzero"`
	Port           string `json:"port" validate:"nonzero"`
	Remoting       string `json:"remoting" validate:"nonzero"`
	LegasyExplorer bool   `json:"legasyExplorer" validate:"regexp=^(true|false)$`
}

//ConfigInterface is an interface for all config structures
type ConfigInterface interface {
}

//PersistedData stores the information about all config types in database and is used during searching for a config by name and type
type PersistedData struct {
	ConfigType ConfigInterface
	IDField    string
}
