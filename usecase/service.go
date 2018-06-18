package usecase

import (
	"encoding/json"
	"time"

	"errors"

	pb "github.com/YAWAL/GetMeConfAPI/api"

	"github.com/YAWAL/GetMeConf/entity"
	"github.com/YAWAL/GetMeConf/repository"
	"github.com/patrickmn/go-cache"
	"gopkg.in/validator.v2"
)

const (
	mongodb    = "mongodb"
	tempconfig = "tempconfig"
	tsconfig   = "tsconfig"
)

type ConfigServer interface {
	GetConfigByName(name, confType string) (*pb.GetConfigResponce, error)
	GetConfigsByType(confType string, stream pb.ConfigService_GetConfigsByTypeServer) error
	CreateConfig(config *pb.Config) (*pb.Responce, error)
	DeleteConfig(delConfigRequest *pb.DeleteConfigRequest) (*pb.Responce, error)
	UpdateConfig(config *pb.Config) (*pb.Responce, error)
}

type configServerImpl struct {
	configCache *cache.Cache
	repo        *repository.StoragePostgresStorage
}

type ServiceConfiguration struct {
	Port      string
	CacheConf *CacheConfiguration
}

type CacheConfiguration struct {
	CacheExpirationTime  int
	CacheCleanupInterval int
}

// NewSharedConfigInteractor constructs new SharedConfigInteractor
func NewConfigServer(s *repository.PostgresStorage, sc *ServiceConfiguration) *configServerImpl {
	//	repo, _ := repository.CreatePostgresStorage()
	return &configServerImpl{
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
func (s *configServerImpl) GetConfigByName(name, confType string) (*pb.GetConfigResponce, error) {
	configResponse, found := s.configCache.Get(name)
	if found {
		return configResponse.(*pb.GetConfigResponce), nil
	}
	var err error
	var res entity.ConfigInterface

	switch confType {
	case mongodb:
		res, err = s.repo.MongoDBRepo.Find(name)
		if err != nil {
			return nil, err
		}
	case tempconfig:
		res, err = s.repo.TempRepo.Find(name)
		if err != nil {
			return nil, err
		}
	case tsconfig:
		res, err = s.repo.TsRepo.Find(name)
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

// GetConfigByName streams configs as GetConfigResponce messages.
func (s *configServerImpl) GetConfigsByType(confType string, stream pb.ConfigService_GetConfigsByTypeServer) error {
	switch confType {
	case mongodb:
		res, err := s.repo.MongoDBRepo.FindAll()
		if err != nil {
			return err
		}
		for _, v := range res {
			byteRes, err := json.Marshal(v)
			if err != nil {
				return err
			}
			if err = stream.Send(&pb.GetConfigResponce{Config: byteRes}); err != nil {
				return err
			}
		}
	case tempconfig:
		res, err := s.repo.TempRepo.FindAll()
		if err != nil {
			return err
		}
		for _, v := range res {
			byteRes, err := json.Marshal(v)
			if err != nil {
				return err
			}
			if err = stream.Send(&pb.GetConfigResponce{Config: byteRes}); err != nil {
				return err
			}
		}
	case tsconfig:
		res, err := s.repo.TsRepo.FindAll()
		if err != nil {
			return err
		}
		for _, v := range res {
			byteRes, err := json.Marshal(v)
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

// CreateConfig calls the function from database package to add a new config record to the database, returns response structure containing a status message.
func (s *configServerImpl) CreateConfig(config *pb.Config) (*pb.Responce, error) {
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
		response, err := s.repo.MongoDBRepo.Save(&configStr)
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
		response, err := s.repo.TempRepo.Save(&configStr)
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
		response, err := s.repo.TsRepo.Save(&configStr)
		if err != nil {
			return nil, err
		}
		s.configCache.Flush()
		return &pb.Responce{Status: response}, nil
	default:
		return nil, errors.New("unexpected type " + config.ConfigType)
	}
}

// DeleteConfig removes config records from the database. If successful, returns the amount of deleted records in a status message of the response structure.
func (s *configServerImpl) DeleteConfig(delConfigRequest *pb.DeleteConfigRequest) (*pb.Responce, error) {
	switch delConfigRequest.ConfigType {
	case mongodb:
		response, err := s.repo.MongoDBRepo.Delete(delConfigRequest.ConfigName)
		if err != nil {
			return nil, err
		}
		s.configCache.Flush()
		return &pb.Responce{Status: response}, nil
	case tempconfig:
		response, err := s.repo.TempRepo.Delete(delConfigRequest.ConfigName)
		if err != nil {
			return nil, err
		}
		s.configCache.Flush()
		return &pb.Responce{Status: response}, nil
	case tsconfig:
		response, err := s.repo.TsRepo.Delete(delConfigRequest.ConfigName)
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
func (s *configServerImpl) UpdateConfig(config *pb.Config) (*pb.Responce, error) {
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
		status, err = s.repo.MongoDBRepo.Update(&configStr)
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
		status, err = s.repo.TempRepo.Update(&configStr)
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
		status, err = s.repo.TsRepo.Update(&configStr)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unexpected type " + config.ConfigType)
	}
	s.configCache.Flush()
	return &pb.Responce{Status: status}, nil
}
