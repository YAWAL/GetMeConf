package service

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/YAWAL/GetMeConf/repository"
	"github.com/YAWAL/GetMeConf/service/errortype"
	"github.com/YAWAL/GetMeConf/usecase"
	pb "github.com/YAWAL/GetMeConfAPI/api"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	pdbScheme = "PDB_SCHEME"
	pdbDSN    = "PDB_DSN"
	//pdbHost             	= "PDB_HOST"
	//pdbPort             	= "PDB_PORT"
	//pdbUser             	= "PDB_USER"
	//pdbPassword         	= "PDB_PASSWORD"
	//pdbName             	= "PDB_NAME"
	maxOpCon            = "MAX_OPENED_CONNECTIONS_TO_DB"
	maxIdleCon          = "MAX_IDLE_CONNECTIONS_TO_DB"
	vConnMaxLifetimeMin = "MB_CONN_MAX_LIFETIME_MINUTES"

	servicePort     = "SERVICE_PORT"
	cacheExpTime    = "CACHE_EXPIRATION_TIME"
	cacheCleanupInt = "CACHE_CLEANUP_INTERVAL"
)

var (
	defaultDbScheme = "postgres"
	defaultDbDSN    = "postgres://dlxifkbx:L7Cey-ucPY4L3T6VFlFdNykNE4jO0VjV@horton.elephantsql.com:5432/dlxifkbx?sslmode=disable"
	//defaultDbHost     = "horton.elephantsql.com"
	//defaultDbPort     = "5432"
	//defaultDbUser     = "dlxifkbx"
	//defaultDbPassword = "L7Cey-ucPY4L3T6VFlFdNykNE4jO0VjV"
	//defaultDbName     = "dlxifkbx"

	defaultMaxOpenedConnectionsToDb = 5
	defaultMaxIdleConnectionsToDb   = 0
	defaultmbConnMaxLifetimeMinutes = 30

	defaultPort                 = "3000"
	defaultCacheExpirationTime  = 5
	defaultCacheCleanupInterval = 10
)

// ServiceConfig structure contains the configuration information for the database.
//type postgresConfig struct {
//	dbSchema                 string
//	dbHost                   string `yaml:"dbhost"`
//	dbPort                   string `yaml:"dbport"`
//	dbUser                   string `yaml:"dbUser"`
//	dbPassword               string `yaml:"dbPassword"`
//	dbName                   string `yaml:"dbName"`
//	maxOpenedConnectionsToDb int    `yaml:"maxOpenedConnectionsToDb"`
//	maxIdleConnectionsToDb   int    `yaml:"maxIdleConnectionsToDb"`
//	mbConnMaxLifetimeMinutes int    `yaml:"mbConnMaxLifetimeMinutes"`
//}

type ConfigGRPCServer struct {
	log          *zap.Logger
	ctx          context.Context
	configServer usecase.ConfigServer
}

// NewConfigGRPCServer returns a new instance of ConfigGRPCServer
func NewConfigGRPCServer(ctx context.Context, log *zap.Logger, cs usecase.ConfigServer) *ConfigGRPCServer {
	return &ConfigGRPCServer{
		ctx:          ctx,
		configServer: cs,
		log:          log.With(zap.String("config service", "SharedConfigHandler")),
	}
}

// GetConfigByName returns one config in GetConfigResponce message.
func (s *ConfigGRPCServer) GetConfigByName(ctx context.Context, nameRequest *pb.GetConfigByNameRequest) (*pb.GetConfigResponce, error) {
	res, err := s.configServer.GetConfigByName(nameRequest.ConfigName, nameRequest.ConfigType)
	if err != nil {
		msg := "couldn't get config"
		s.log.Error(err.Error(), zap.String("error", msg))
		return nil, errortype.GrpcError(err, msg)
	}
	return res, nil
}

// GetConfigByName streams configs as GetConfigResponce messages.
func (s *ConfigGRPCServer) GetConfigsByType(typeRequest *pb.GetConfigsByTypeRequest, stream pb.ConfigService_GetConfigsByTypeServer) error {
	err := s.configServer.GetConfigsByType(typeRequest.ConfigType, stream)
	if err != nil {
		msg := "couldn't get configs"
		s.log.Error(err.Error(), zap.String("error", msg))
		return errortype.GrpcError(err, msg)
	}
	return nil
}

// CreateConfig calls the function from database package to add a new config record to the database,
// returns response structure containing a status message.
func (s *ConfigGRPCServer) CreateConfig(ctx context.Context, config *pb.Config) (*pb.Responce, error) {
	res, err := s.configServer.CreateConfig(config)
	if err != nil {
		msg := "couldn't create config"
		s.log.Error(err.Error(), zap.String("error", msg))
		return nil, errortype.GrpcError(err, msg)
	}
	return res, nil
}

// DeleteConfig removes config records from the database.
// If successful, returns the amount of deleted records in a status message of the response structure.
func (s *ConfigGRPCServer) DeleteConfig(ctx context.Context, delConfigRequest *pb.DeleteConfigRequest) (*pb.Responce, error) {
	res, err := s.configServer.DeleteConfig(delConfigRequest)
	if err != nil {
		msg := "couldn't delete config"
		s.log.Error(err.Error(), zap.String("error", msg))
		return nil, errortype.GrpcError(err, msg)
	}
	return res, nil
}

// UpdateConfig updates a config stored in database.
func (s *ConfigGRPCServer) UpdateConfig(ctx context.Context, config *pb.Config) (*pb.Responce, error) {
	res, err := s.configServer.UpdateConfig(config)
	if err != nil {
		msg := "couldn't update config"
		s.log.Error(err.Error(), zap.String("error", msg))
		return nil, errortype.GrpcError(err, msg)
	}
	return res, nil
}

func initServiceConfiguration(logger *zap.Logger) *usecase.ServiceConfiguration {
	port := os.Getenv(servicePort)
	if port == "" {
		logger.Info("error during reading env. variable", zap.String("default value is used", defaultPort))
		port = defaultPort
	}
	cacheConfiguration := initCacheConfiguration(logger)
	return &usecase.ServiceConfiguration{Port: port, CacheConf: cacheConfiguration}
}

// Run function starts the service.
func Run() {

	ctx, done := context.WithCancel(context.Background())
	defer done()

	var err error
	logger, err := zap.NewProduction()
	if err != nil {
		log.Printf("Error has occurred during logger initialization: %v", err)
	}
	defer logger.Sync()

	serviceConfiguration := initServiceConfiguration(logger)

	postgresConfig := initPostgresConfig(logger)

	postgresStorage, err := repository.NewPostgresStorage(postgresConfig)
	if err != nil {
		logger.Fatal("failed to init postgres db", zap.Error(err))
	}

	postgresStorage.Migrate()

	server := usecase.NewConfigServer(postgresStorage, serviceConfiguration)

	grpcServ := NewConfigGRPCServer(ctx, logger, server)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", serviceConfiguration.Port))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}
	logger.Info("Server started at", zap.String("port", serviceConfiguration.Port))

	grpcServer := grpc.NewServer()

	pb.RegisterConfigServiceServer(grpcServer, grpcServ)

	go func() {
		logger.Fatal("failed to serve", zap.Error(grpcServer.Serve(lis)))
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	logger.Info("shotdown signal received, exiting")
	grpcServer.GracefulStop()
}

func validatePostgresConfig(logger *zap.Logger, c *repository.PostgresConfig) {
	if c.Shema == "" {
		logger.Info("error during reading env. variable", zap.String("default value is used ", defaultDbScheme))
		c.Shema = defaultDbScheme
	}
	if c.DSN == "" {
		logger.Info("error during reading env. variable", zap.String("default value is used ", defaultDbDSN))
		c.DSN = defaultDbDSN
	}
	//if c.dbHost == "" {
	//	logger.Info("error during reading env. variable", zap.String("default value is used ", defaultDbHost))
	//	c.dbHost = defaultDbHost
	//}
	//if c.dbPort == "" {
	//	logger.Info("error during reading env. variable", zap.String("default value is used ", defaultDbPort))
	//	c.dbPort = defaultDbPort
	//}
	//if c.dbUser == "" {
	//	logger.Info("error during reading env. variable", zap.String("default value is used ", defaultDbUser))
	//	c.dbUser = defaultDbUser
	//}
	//if c.dbPassword == "" {
	//	logger.Info("error during reading env. variable", zap.String("default value is used ", defaultDbPassword))
	//	c.dbPassword = defaultDbPassword
	//}
	//if c.dbName == "" {
	//	logger.Info("error during reading env. variable", zap.String("default value is used ", defaultDbName))
	//	c.dbName = defaultDbName
	//}
	if c.MaxOpenedConnectionsToDb == 0 {
		logger.Info("maxOpenedConnectionsToDb = 0", zap.Int("default value is used ", defaultMaxOpenedConnectionsToDb))
		c.MaxOpenedConnectionsToDb = defaultMaxOpenedConnectionsToDb
	}
	if c.MaxIdleConnectionsToDb == 0 {
		logger.Info("maxIdleConnectionsToDb = 0", zap.Int("default value is used ", defaultMaxIdleConnectionsToDb))
		c.MaxIdleConnectionsToDb = defaultMaxIdleConnectionsToDb
	}
	if c.MbConnMaxLifetimeMinutes == 0 {
		logger.Info("mbConnMaxLifetimeMinutes = 0", zap.Int("default value is used ", defaultmbConnMaxLifetimeMinutes))
		c.MbConnMaxLifetimeMinutes = defaultmbConnMaxLifetimeMinutes
	}
}

func initPostgresConfig(logger *zap.Logger) *repository.PostgresConfig {
	c := new(repository.PostgresConfig)
	c.Shema = os.Getenv(pdbScheme)
	c.DSN = os.Getenv(pdbDSN)
	//c.dbHost = os.Getenv(pdbHost)
	//c.dbPort = os.Getenv(pdbPort)
	//c.dbUser = os.Getenv(pdbUser)
	//c.dbPassword = os.Getenv(pdbPassword)
	//c.dbName = os.Getenv(pdbName)
	var err error
	c.MaxOpenedConnectionsToDb, err = strconv.Atoi(os.Getenv(maxOpCon))
	if err != nil {
		logger.Info("error during reading env. variable. Could not convert from string to int", zap.Error(err))
	}
	c.MaxIdleConnectionsToDb, err = strconv.Atoi(os.Getenv(maxIdleCon))
	if err != nil {
		logger.Info("error during reading env. variable. Could not convert from string to int", zap.Error(err))
	}
	c.MbConnMaxLifetimeMinutes, err = strconv.Atoi(os.Getenv(vConnMaxLifetimeMin))
	if err != nil {
		logger.Info("error during reading env. variable. Could not convert from string to int", zap.Error(err))
	}
	validatePostgresConfig(logger, c)
	return c
}

func initCacheConfiguration(logger *zap.Logger) *usecase.CacheConfiguration {
	c := new(usecase.CacheConfiguration)
	cacheExpirationTime, err := strconv.Atoi(os.Getenv(cacheExpTime))
	if err != nil {
		logger.Info("error during reading env. variable", zap.Int("default value is used", defaultCacheExpirationTime))
		cacheExpirationTime = defaultCacheExpirationTime
	}
	c.CacheExpirationTime = cacheExpirationTime
	cacheCleanupInterval, err := strconv.Atoi(os.Getenv(cacheCleanupInt))
	if err != nil {
		logger.Info("error during reading env. variable", zap.Int("default value is used", defaultCacheCleanupInterval))
		cacheCleanupInterval = defaultCacheCleanupInterval
	}
	c.CacheCleanupInterval = cacheCleanupInterval
	return c
}
