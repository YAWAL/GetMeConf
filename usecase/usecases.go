package usecase

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/YAWAL/GetMeConf/entity"
	"github.com/YAWAL/GetMeConf/repository"
	pb "github.com/YAWAL/GetMeConfAPI/api"
	"github.com/patrickmn/go-cache"
	"gopkg.in/validator.v2"
)

const (
	mongodb    = "mongodb"
	tempconfig = "tempconfig"
	tsconfig   = "tsconfig"
)

// ConfigServer interface contains server functions.
type ConfigServer interface {
	GetConfigByName(name, confType string) (*pb.GetConfigResponce, error)
	GetConfigsByType(confType string, stream pb.ConfigService_GetConfigsByTypeServer) error
	CreateConfig(config *pb.Config) (*pb.Responce, error)
	DeleteConfig(delConfigRequest *pb.DeleteConfigRequest) (*pb.Responce, error)
	UpdateConfig(config *pb.Config) (*pb.Responce, error)
}

// ConfigServerImpl wraps the cache and a storage.
type ConfigServerImpl struct {
	configCache *cache.Cache
	repo        repository.Storage
}

// ServiceConfiguration contains service parameters.
type ServiceConfiguration struct {
	Port      string
	CacheConf *CacheConfiguration
}

// CacheConfiguration contains the cache parrameters.
type CacheConfiguration struct {
	CacheExpirationTime  int
	CacheCleanupInterval int
}

// NewConfigServer constructs new ConfigServer
func NewConfigServer(s repository.Storage, sc *ServiceConfiguration) *ConfigServerImpl {
	return &ConfigServerImpl{
		repo:        s,
		configCache: initCache(sc),
	}
}

func initCache(cc *ServiceConfiguration) *cache.Cache {
	return cache.New(
		time.Duration(cc.CacheConf.CacheExpirationTime)*time.Minute,
		time.Duration(cc.CacheConf.CacheCleanupInterval)*time.Minute,
	)
}

// GetConfigByName returns one config in GetConfigResponce message.
func (s *ConfigServerImpl) GetConfigByName(name, confType string) (*pb.GetConfigResponce, error) {
	configResponse, found := s.configCache.Get(name)
	if found {
		return configResponse.(*pb.GetConfigResponce), nil
	}
	var err error
	var res entity.ConfigInterface
	switch confType {
	case mongodb:
		res, err = s.repo.FindMongoDBConfig(name)
		if err != nil {
			return nil, err
		}
	case tempconfig:
		res, err = s.repo.FindTempConfig(name)
		if err != nil {
			return nil, err
		}
	case tsconfig:
		res, err = s.repo.FindTsConfig(name)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unexpected type " + confType)
	}
	byteRes, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	configResponse = &pb.GetConfigResponce{Config: byteRes}
	s.configCache.Set(name, configResponse, cache.DefaultExpiration)
	return configResponse.(*pb.GetConfigResponce), nil
}

// GetConfigsByType streams configs as GetConfigResponce messages.
func (s *ConfigServerImpl) GetConfigsByType(confType string, stream pb.ConfigService_GetConfigsByTypeServer) error {
	switch confType {
	case mongodb:
		res, err := s.repo.FindAllMongoDBConfig()
		if err != nil {
			return err
		}
		for k := range res {
			byteRes, err := json.Marshal(res[k])
			if err != nil {
				return err
			}
			if err = stream.Send(&pb.GetConfigResponce{Config: byteRes}); err != nil {
				return err
			}
		}
	case tempconfig:
		res, err := s.repo.FindAllTempConfig()
		if err != nil {
			return err
		}
		for k := range res {
			byteRes, err := json.Marshal(res[k])
			if err != nil {
				return err
			}
			if err = stream.Send(&pb.GetConfigResponce{Config: byteRes}); err != nil {
				return err
			}
		}
	case tsconfig:
		res, err := s.repo.FindAllTsConfig()
		if err != nil {
			return err
		}
		for k := range res {
			byteRes, err := json.Marshal(res[k])
			if err != nil {
				return err
			}
			if err = stream.Send(&pb.GetConfigResponce{Config: byteRes}); err != nil {
				return err
			}
		}
	default:
		return errors.New("unexpected type " + confType)
	}
	return nil
}

// CreateConfig calls the function from database package to add a new config record to the database, returns response
// structure containing a status message.
func (s *ConfigServerImpl) CreateConfig(config *pb.Config) (*pb.Responce, error) {
	switch config.ConfigType {
	case mongodb:
		configStr := entity.Mongodb{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			return nil, err
		}
		if err = validator.Validate(configStr); err != nil {
			return nil, err
		}
		response, err := s.repo.SaveMongoDBConfig(&configStr)
		if err != nil {
			return nil, err
		}
		s.configCache.Flush()
		return &pb.Responce{Status: response}, nil

	case tempconfig:
		configStr := entity.Tempconfig{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			return nil, err
		}
		if err = validator.Validate(configStr); err != nil {
			return nil, err
		}
		response, err := s.repo.SaveTempConfig(&configStr)
		if err != nil {
			return nil, err
		}
		s.configCache.Flush()
		return &pb.Responce{Status: response}, nil

	case tsconfig:
		configStr := entity.Tsconfig{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			return nil, err
		}
		if err = validator.Validate(configStr); err != nil {
			return nil, err
		}
		response, err := s.repo.SaveTsConfig(&configStr)
		if err != nil {
			return nil, err
		}
		s.configCache.Flush()
		return &pb.Responce{Status: response}, nil
	default:
		return nil, errors.New("unexpected type " + config.ConfigType)
	}
}

// DeleteConfig removes config records from the database. If successful, returns the amount of deleted records
// in a status message of the response structure.
func (s *ConfigServerImpl) DeleteConfig(delConfigRequest *pb.DeleteConfigRequest) (*pb.Responce, error) {
	switch delConfigRequest.ConfigType {
	case mongodb:
		response, err := s.repo.DeleteMongoDBConfig(delConfigRequest.ConfigName)
		if err != nil {
			return nil, err
		}
		s.configCache.Flush()
		return &pb.Responce{Status: response}, nil
	case tempconfig:
		response, err := s.repo.DeleteTempConfig(delConfigRequest.ConfigName)
		if err != nil {
			return nil, err
		}
		s.configCache.Flush()
		return &pb.Responce{Status: response}, nil
	case tsconfig:
		response, err := s.repo.DeleteTsConfig(delConfigRequest.ConfigName)
		if err != nil {
			return nil, err
		}
		s.configCache.Flush()
		return &pb.Responce{Status: response}, nil
	default:
		return nil, errors.New("unexpected type " + delConfigRequest.ConfigType)
	}
}

// UpdateConfig updates a config stored in database.
func (s *ConfigServerImpl) UpdateConfig(config *pb.Config) (*pb.Responce, error) {
	var status string
	switch config.ConfigType {
	case mongodb:
		configStr := entity.Mongodb{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			return nil, err
		}
		if err = validator.Validate(configStr); err != nil {
			return nil, err
		}
		status, err = s.repo.UpdateMongoDBConfig(&configStr)
		if err != nil {
			return nil, err
		}
	case tempconfig:
		configStr := entity.Tempconfig{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			return nil, err
		}
		if err = validator.Validate(configStr); err != nil {
			return nil, err
		}
		status, err = s.repo.UpdateTempConfig(&configStr)
		if err != nil {
			return nil, err
		}
	case tsconfig:
		configStr := entity.Tsconfig{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			return nil, err
		}
		if err = validator.Validate(configStr); err != nil {
			return nil, err
		}
		status, err = s.repo.UpdateTsConfig(&configStr)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unexpected type " + config.ConfigType)
	}
	s.configCache.Flush()
	return &pb.Responce{Status: status}, nil
}
