package main

import (
	"encoding/json"
	"log"
	"time"

	"errors"

	"os"

	pb "github.com/YAWAL/GetMeConf/api"
	micro "github.com/micro/go-micro"

	"strconv"

	"github.com/YAWAL/GetMeConf/entitie"
	"github.com/YAWAL/GetMeConf/repository"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/server"
	"github.com/patrickmn/go-cache"
	"golang.org/x/net/context"
	"gopkg.in/validator.v2"
)

var (
	defaultPort                 = "3000"
	defaultCacheExpirationTime  = 5
	defaultCacheCleanupInterval = 10
)

const (
	mongodb    = "mongodb"
	tempconfig = "tempconfig"
	tsconfig   = "tsconfig"
)

type configServer struct {
	configCache       *cache.Cache
	mongoDBConfigRepo repository.MongoDBConfigRepo
	tempConfigRepo    repository.TempConfigRepo
	tsConfigRepo      repository.TsConfigRepo
	PubSub            broker.Broker
}

type serviceConfiguration struct {
	port                 string
	cacheExpirationTime  int
	cacheCleanupInterval int
}

//GetConfigByName returns one config in GetConfigResponce message
func (s *configServer) GetConfigByName(ctx context.Context, nameRequest *pb.GetConfigByNameRequest, configResponce *pb.GetConfigResponce) error {

	configFromCache, found := s.configCache.Get(nameRequest.ConfigName)
	if found {
		configResponce.Config = configFromCache.(*pb.GetConfigResponce).Config
		return nil
	}
	var err error
	var res entitie.ConfigInterface

	switch nameRequest.ConfigType {
	case mongodb:
		res, err = s.mongoDBConfigRepo.Find(nameRequest.ConfigName)
		if err != nil {
			return err
		}
	case tempconfig:
		res, err = s.tempConfigRepo.Find(nameRequest.ConfigName)
		if err != nil {
			return err
		}
	case tsconfig:
		res, err = s.tsConfigRepo.Find(nameRequest.ConfigName)
		if err != nil {
			return err
		}
	default:
		log.Print("unexpected type")
		return errors.New("unexpected type")
	}
	byteRes, err := json.Marshal(res)
	if err != nil {
		return err
	}
	configResponce.Config = byteRes
	s.configCache.Set(nameRequest.ConfigName, configResponce, cache.DefaultExpiration)
	return nil
}

//GetConfigByName streams configs as GetConfigResponce messages
func (s *configServer) GetConfigsByType(ctx context.Context, typeRequest *pb.GetConfigsByTypeRequest, stream pb.ConfigService_GetConfigsByTypeStream) error {
	switch typeRequest.ConfigType {
	case mongodb:
		res, err := s.mongoDBConfigRepo.FindAll()
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
		res, err := s.tempConfigRepo.FindAll()
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
		res, err := s.tsConfigRepo.FindAll()
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
		log.Print("unexpected type")
		return errors.New("unexpected type")
	}
	return nil
}

//CreateConfig calls the function from database package to add a new config record to the database, returns response structure containing a status message
func (s *configServer) CreateConfig(ctx context.Context, config *pb.Config, responce *pb.Responce) error {
	var responceStatus string
	switch config.ConfigType {
	case mongodb:
		configStr := entitie.Mongodb{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			log.Printf("unmarshal config err: %v", err)
			return err
		}
		if err = validator.Validate(configStr); err != nil {
			return err
		}
		responceStatus, err = s.mongoDBConfigRepo.Save(&configStr)
		if err != nil {
			return err
		}
		s.configCache.Flush()
		responce.Status = responceStatus
		return nil

	case tempconfig:
		configStr := entitie.Tempconfig{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			log.Printf("unmarshal config err: %v", err)
			return err
		}
		if err = validator.Validate(configStr); err != nil {
			return err
		}
		responceStatus, err = s.tempConfigRepo.Save(&configStr)
		if err != nil {
			return err
		}
		s.configCache.Flush()
		responce.Status = responceStatus
		return nil

	case tsconfig:
		configStr := entitie.Tsconfig{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			log.Printf("unmarshal config err: %v", err)
			return err
		}
		if err = validator.Validate(configStr); err != nil {
			return err
		}
		responceStatus, err = s.tsConfigRepo.Save(&configStr)
		if err != nil {
			return err
		}
		s.configCache.Flush()
		responce.Status = responceStatus
		return nil
	default:
		log.Print("unexpected type")
		return errors.New("unexpected type")
	}
}

//DeleteConfig removes config records from the database. If successful, returns the amount of deleted records in a status message of the response structure
func (s *configServer) DeleteConfig(ctx context.Context, delConfigRequest *pb.DeleteConfigRequest, response *pb.Responce) error {
	switch delConfigRequest.ConfigType {
	case mongodb:
		responseStatus, err := s.mongoDBConfigRepo.Delete(delConfigRequest.ConfigName)
		if err != nil {
			return err
		}
		s.configCache.Flush()
		response.Status = responseStatus
		return nil
	case tempconfig:
		responseStatus, err := s.tempConfigRepo.Delete(delConfigRequest.ConfigName)
		if err != nil {
			return err
		}
		s.configCache.Flush()
		response.Status = responseStatus
		return nil
	case tsconfig:
		responseStatus, err := s.tsConfigRepo.Delete(delConfigRequest.ConfigName)
		if err != nil {
			return err
		}
		s.configCache.Flush()
		response.Status = responseStatus
		return nil
	default:
		log.Print("unexpected type")
		return errors.New("unexpected type")
	}
}

//UpdateConfig
func (s *configServer) UpdateConfig(ctx context.Context, config *pb.Config, responce *pb.Responce) error {
	var status string
	switch config.ConfigType {
	case mongodb:
		configStr := entitie.Mongodb{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			log.Printf("unmarshal config err: %v", err)
			return err
		}
		if err = validator.Validate(configStr); err != nil {
			return err
		}
		status, err = s.mongoDBConfigRepo.Update(&configStr)
		if err != nil {
			return err
		}
	case tempconfig:
		configStr := entitie.Tempconfig{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			log.Printf("unmarshal config err: %v", err)
			return err
		}
		if err = validator.Validate(configStr); err != nil {
			return err
		}
		status, err = s.tempConfigRepo.Update(&configStr)
		if err != nil {
			return err
		}
	case tsconfig:
		configStr := entitie.Tsconfig{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			log.Printf("unmarshal config err: %v", err)
			return err
		}
		if err = validator.Validate(configStr); err != nil {
			return err
		}
		status, err = s.tsConfigRepo.Update(&configStr)
		if err != nil {
			return err
		}
	default:
		log.Print("unexpected type")
		return errors.New("unexpected type")
	}
	s.configCache.Flush()
	responce.Status = status
	return nil
}

func initServiceConfiguration() *serviceConfiguration {
	cacheExpirationTime, err := strconv.Atoi(os.Getenv("CACHE_EXPIRATION_TIME"))
	if err != nil {
		log.Printf("error during reading env. variable: %v, default value is used", err)
		cacheExpirationTime = defaultCacheExpirationTime
	}
	cacheCleanupInterval, err := strconv.Atoi(os.Getenv("CACHE_CLEANUP_INTERVAL"))
	if err != nil {
		log.Printf("error during reading env. variable: %v, default value is used", err)
		cacheCleanupInterval = defaultCacheCleanupInterval
	}
	return &serviceConfiguration{cacheCleanupInterval: cacheCleanupInterval, cacheExpirationTime: cacheExpirationTime}
}

func main() {

	serviceConfiguration := initServiceConfiguration()
	dbConn, err := repository.InitPostgresDB()
	if err != nil {
		log.Fatalf("failed to init postgres db: %v", err)
	}
	mongoDBRepo := repository.MongoDBConfigRepoImpl{DB: dbConn}
	tsConfigRepo := repository.TsConfigRepoImpl{DB: dbConn}
	tempConfigRepo := repository.TempConfigRepoImpl{DB: dbConn}

	srv := micro.NewService(
		micro.Name("api"),
		micro.Version("latest"),
	)

	srv.Server().Init(server.Wait(true))

	configCache := cache.New(time.Duration(serviceConfiguration.cacheExpirationTime)*time.Minute, time.Duration(serviceConfiguration.cacheCleanupInterval)*time.Minute)

	pb.RegisterConfigServiceHandler(srv.Server(), &configServer{configCache: configCache, mongoDBConfigRepo: &mongoDBRepo, tsConfigRepo: &tsConfigRepo, tempConfigRepo: &tempConfigRepo})

	if err := srv.Run(); err != nil {
		log.Println(err)
	}
}
