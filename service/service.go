package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"errors"

	"os"

	pb "github.com/YAWAL/GetMeConfAPI/api"

	"strconv"

	"os/signal"
	"syscall"

	"github.com/YAWAL/GetMeConf/entity"
	"github.com/YAWAL/GetMeConf/repository"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/validator.v2"
)

const (
	mongodb    = "mongodb"
	tempconfig = "tempconfig"
	tsconfig   = "tsconfig"

	servicePort     = "SERVICE_PORT"
	cacheExpTime    = "CACHE_EXPIRATION_TIME"
	cacheCleanupInt = "CACHE_CLEANUP_INTERVAL"
)

var (
	defaultPort                 = "3000"
	defaultCacheExpirationTime  = 5
	defaultCacheCleanupInterval = 10
)

var logger *zap.Logger

type configServer struct {
	configCache       *cache.Cache
	mongoDBConfigRepo repository.MongoDBConfigRepo
	tempConfigRepo    repository.TempConfigRepo
	tsConfigRepo      repository.TsConfigRepo
}

type serviceConfiguration struct {
	port                 string
	cacheExpirationTime  int
	cacheCleanupInterval int
}

// GetConfigByName returns one config in GetConfigResponce message.
func (s *configServer) GetConfigByName(ctx context.Context, nameRequest *pb.GetConfigByNameRequest) (*pb.GetConfigResponce, error) {

	configResponse, found := s.configCache.Get(nameRequest.ConfigName)
	if found {
		return configResponse.(*pb.GetConfigResponce), nil
	}
	var err error
	var res entity.ConfigInterface

	switch nameRequest.ConfigType {
	case mongodb:
		res, err = s.mongoDBConfigRepo.Find(nameRequest.ConfigName)
		if err != nil {
			return nil, err
		}
	case tempconfig:
		res, err = s.tempConfigRepo.Find(nameRequest.ConfigName)
		if err != nil {
			return nil, err
		}
	case tsconfig:
		res, err = s.tsConfigRepo.Find(nameRequest.ConfigName)
		if err != nil {
			return nil, err
		}
	default:
		logger.Info("unexpected type", zap.String("Such config does not exist: ", nameRequest.ConfigType))
		return nil, errors.New("unexpected type")
	}
	byteRes, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	configResponse = &pb.GetConfigResponce{Config: byteRes}
	s.configCache.Set(nameRequest.ConfigName, configResponse, cache.DefaultExpiration)
	return configResponse.(*pb.GetConfigResponce), nil
}

// GetConfigByName streams configs as GetConfigResponce messages.
func (s *configServer) GetConfigsByType(typeRequest *pb.GetConfigsByTypeRequest, stream pb.ConfigService_GetConfigsByTypeServer) error {
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
		logger.Info("unexpected type", zap.String("Such config does not exist: ", typeRequest.ConfigType))
		return errors.New("unexpected type")
	}
	return nil
}

// CreateConfig calls the function from database package to add a new config record to the database, returns response structure containing a status message.
func (s *configServer) CreateConfig(ctx context.Context, config *pb.Config) (*pb.Responce, error) {
	switch config.ConfigType {
	case mongodb:
		configStr := entity.Mongodb{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			logger.Info("unmarshal config error", zap.Error(err))
			return nil, err
		}
		if err = validator.Validate(configStr); err != nil {
			return nil, err
		}
		response, err := s.mongoDBConfigRepo.Save(&configStr)
		if err != nil {
			return nil, err
		}
		s.configCache.Flush()
		return &pb.Responce{Status: response}, nil

	case tempconfig:
		configStr := entity.Tempconfig{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			logger.Info("unmarshal config error", zap.Error(err))
			return nil, err
		}
		if err = validator.Validate(configStr); err != nil {
			return nil, err
		}
		response, err := s.tempConfigRepo.Save(&configStr)
		if err != nil {
			return nil, err
		}
		s.configCache.Flush()
		return &pb.Responce{Status: response}, nil

	case tsconfig:
		configStr := entity.Tsconfig{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			logger.Info("unmarshal config error", zap.Error(err))
			return nil, err
		}
		if err = validator.Validate(configStr); err != nil {
			return nil, err
		}
		response, err := s.tsConfigRepo.Save(&configStr)
		if err != nil {
			return nil, err
		}
		s.configCache.Flush()
		return &pb.Responce{Status: response}, nil
	default:
		logger.Info("unexpected type", zap.String("Such config does not exist: ", config.ConfigType))
		return nil, errors.New("unexpected type")
	}
}

// DeleteConfig removes config records from the database. If successful, returns the amount of deleted records in a status message of the response structure.
func (s *configServer) DeleteConfig(ctx context.Context, delConfigRequest *pb.DeleteConfigRequest) (*pb.Responce, error) {
	switch delConfigRequest.ConfigType {
	case mongodb:
		response, err := s.mongoDBConfigRepo.Delete(delConfigRequest.ConfigName)
		if err != nil {
			return nil, err
		}
		s.configCache.Flush()
		return &pb.Responce{Status: response}, nil
	case tempconfig:
		response, err := s.tempConfigRepo.Delete(delConfigRequest.ConfigName)
		if err != nil {
			return nil, err
		}
		s.configCache.Flush()
		return &pb.Responce{Status: response}, nil
	case tsconfig:
		response, err := s.tsConfigRepo.Delete(delConfigRequest.ConfigName)
		if err != nil {
			return nil, err
		}
		s.configCache.Flush()
		return &pb.Responce{Status: response}, nil
	default:
		logger.Info("unexpected type", zap.String("Such config does not exist: ", delConfigRequest.ConfigType))
		return nil, errors.New("unexpected type")
	}
}

// UpdateConfig updates a config stored in database.
func (s *configServer) UpdateConfig(ctx context.Context, config *pb.Config) (*pb.Responce, error) {
	var status string
	switch config.ConfigType {
	case mongodb:
		configStr := entity.Mongodb{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			logger.Info("unmarshal config error", zap.Error(err))
			return nil, err
		}
		if err = validator.Validate(configStr); err != nil {
			return nil, err
		}
		status, err = s.mongoDBConfigRepo.Update(&configStr)
		if err != nil {
			return nil, err
		}
	case tempconfig:
		configStr := entity.Tempconfig{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			logger.Info("unmarshal config error", zap.Error(err))
			return nil, err
		}
		if err = validator.Validate(configStr); err != nil {
			return nil, err
		}
		status, err = s.tempConfigRepo.Update(&configStr)
		if err != nil {
			return nil, err
		}
	case tsconfig:
		configStr := entity.Tsconfig{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			logger.Info("unmarshal config error", zap.Error(err))
			return nil, err
		}
		if err = validator.Validate(configStr); err != nil {
			return nil, err
		}
		status, err = s.tsConfigRepo.Update(&configStr)
		if err != nil {
			return nil, err
		}
	default:
		logger.Info("unexpected type", zap.String("Such config does not exist: ", config.ConfigType))
		return nil, errors.New("unexpected type")
	}
	s.configCache.Flush()
	return &pb.Responce{Status: status}, nil
}

func initServiceConfiguration() *serviceConfiguration {
	port := os.Getenv(servicePort)
	if port == "" {
		logger.Info("error during reading env. variable", zap.String("default value is used", defaultPort))
		port = defaultPort
	}
	cacheExpirationTime, err := strconv.Atoi(os.Getenv(cacheExpTime))
	if err != nil {
		logger.Info("error during reading env. variable", zap.Int("default value is used", defaultCacheExpirationTime))
		cacheExpirationTime = defaultCacheExpirationTime
	}
	cacheCleanupInterval, err := strconv.Atoi(os.Getenv(cacheCleanupInt))
	if err != nil {
		logger.Info("error during reading env. variable", zap.Int("default value is used", defaultCacheCleanupInterval))
		cacheCleanupInterval = defaultCacheCleanupInterval
	}
	return &serviceConfiguration{port: port, cacheCleanupInterval: cacheCleanupInterval, cacheExpirationTime: cacheExpirationTime}
}

// Run function starts the service.
func Run() {

	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Printf("Error has occurred during logger initialization: %v", err)
	}
	defer logger.Sync()
	serviceConfiguration := initServiceConfiguration()

	repository.InitZapLogger(logger)
	dbConn, err := repository.InitPostgresDB()
	if err != nil {
		logger.Fatal("failed to init postgres db", zap.Error(err))
	}
	mongoDBRepo := repository.MongoDBConfigRepoImpl{DB: dbConn}
	tsConfigRepo := repository.TsConfigRepoImpl{DB: dbConn}
	tempConfigRepo := repository.TempConfigRepoImpl{DB: dbConn}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", serviceConfiguration.port))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}
	logger.Info("Server started at", zap.String("port", serviceConfiguration.port))

	grpcServer := grpc.NewServer()

	configCache := cache.New(time.Duration(serviceConfiguration.cacheExpirationTime)*time.Minute, time.Duration(serviceConfiguration.cacheCleanupInterval)*time.Minute)

	pb.RegisterConfigServiceServer(grpcServer, &configServer{configCache: configCache, mongoDBConfigRepo: &mongoDBRepo, tsConfigRepo: &tsConfigRepo, tempConfigRepo: &tempConfigRepo})

	go func() {
		logger.Fatal("failed to serve", zap.Error(grpcServer.Serve(lis)))
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	logger.Info("shotdown signal received, exiting")
	grpcServer.GracefulStop()
}
