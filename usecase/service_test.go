package usecase

import (
	//"context"
	"encoding/json"
	//"os"

	"testing"
	//	"time"

	pb "github.com/YAWAL/GetMeConfAPI/api"
	//	validator "gopkg.in/validator.v2"

	"errors"

	"time"

	"github.com/YAWAL/GetMeConf/entity"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"gopkg.in/validator.v2"
)

const (
	testDefaultExpirationTimeOfCacheMin = 5
	testCleanupInternalOfCacheMin       = 10
)

type mockPostgresStorage struct{}

func (m *mockPostgresStorage) Migrate() error {
	return nil
}

func (m *mockPostgresStorage) FindMongoDBConfig(configName string) (*entity.Mongodb, error) {
	return &entity.Mongodb{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"}, nil
}
func (m *mockPostgresStorage) FindAllMongoDBConfig() ([]entity.Mongodb, error) {
	return []entity.Mongodb{{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"}}, nil
}
func (m *mockPostgresStorage) UpdateMongoDBConfig(config *entity.Mongodb) (string, error) {
	return "OK", nil
}
func (m *mockPostgresStorage) SaveMongoDBConfig(config *entity.Mongodb) (string, error) {
	return "OK", nil
}
func (m *mockPostgresStorage) DeleteMongoDBConfig(configName string) (string, error) { return "OK", nil }

func (m *mockPostgresStorage) FindTempConfig(configName string) (*entity.Tempconfig, error) {
	return &entity.Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}, nil
}
func (m *mockPostgresStorage) FindAllTempConfig() ([]entity.Tempconfig, error) {
	return []entity.Tempconfig{{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}}, nil
}
func (m *mockPostgresStorage) UpdateTempConfig(config *entity.Tempconfig) (string, error) {
	return "OK", nil
}
func (m *mockPostgresStorage) SaveTempConfig(config *entity.Tempconfig) (string, error) {
	return "OK", nil
}
func (m *mockPostgresStorage) DeleteTempConfig(configName string) (string, error) { return "OK", nil }

func (m *mockPostgresStorage) FindTsConfig(configName string) (*entity.Tsconfig, error) {
	return &entity.Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}, nil
}
func (m *mockPostgresStorage) FindAllTsConfig() ([]entity.Tsconfig, error) {
	return []entity.Tsconfig{{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}}, nil
}
func (m *mockPostgresStorage) UpdateTsConfig(config *entity.Tsconfig) (string, error) {
	return "OK", nil
}
func (m *mockPostgresStorage) SaveTsConfig(config *entity.Tsconfig) (string, error) { return "OK", nil }
func (m *mockPostgresStorage) DeleteTsConfig(configName string) (string, error)     { return "OK", nil }

type mockPostgresStorageErr struct{}

func (m *mockPostgresStorageErr) Migrate() error {
	return nil
}

func (m *mockPostgresStorageErr) FindMongoDBConfig(configName string) (*entity.Mongodb, error) {
	return nil, errors.New("error from database querying")
}
func (m *mockPostgresStorageErr) FindAllMongoDBConfig() ([]entity.Mongodb, error) {
	return nil, errors.New("error from database querying")
}
func (m *mockPostgresStorageErr) UpdateMongoDBConfig(config *entity.Mongodb) (string, error) {
	return "", errors.New("error from database querying")
}
func (m *mockPostgresStorageErr) SaveMongoDBConfig(config *entity.Mongodb) (string, error) {
	return "", errors.New("error from database querying")
}
func (m *mockPostgresStorageErr) DeleteMongoDBConfig(configName string) (string, error) {
	return "", errors.New("error from database querying")
}

func (m *mockPostgresStorageErr) FindTempConfig(configName string) (*entity.Tempconfig, error) {
	return nil, errors.New("error from database querying")
}
func (m *mockPostgresStorageErr) FindAllTempConfig() ([]entity.Tempconfig, error) {
	return nil, errors.New("error from database querying")
}
func (m *mockPostgresStorageErr) UpdateTempConfig(config *entity.Tempconfig) (string, error) {
	return "", errors.New("error from database querying")
}
func (m *mockPostgresStorageErr) SaveTempConfig(config *entity.Tempconfig) (string, error) {
	return "", errors.New("error from database querying")
}
func (m *mockPostgresStorageErr) DeleteTempConfig(configName string) (string, error) {
	return "", errors.New("error from database querying")
}

func (m *mockPostgresStorageErr) FindTsConfig(configName string) (*entity.Tsconfig, error) {
	return nil, errors.New("error from database querying")
}
func (m *mockPostgresStorageErr) FindAllTsConfig() ([]entity.Tsconfig, error) {
	return nil, errors.New("error from database querying")
}
func (m *mockPostgresStorageErr) UpdateTsConfig(config *entity.Tsconfig) (string, error) {
	return "", errors.New("error from database querying")
}
func (m *mockPostgresStorageErr) SaveTsConfig(config *entity.Tsconfig) (string, error) {
	return "", errors.New("error from database querying")
}
func (m *mockPostgresStorageErr) DeleteTsConfig(configName string) (string, error) {
	return "", errors.New("error from database querying")
}

//
//
//
//func (m *mockMongoDBConfigRepo) Find(configName string) (*entity.Mongodb, error) {
//	return &entity.Mongodb{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"}, nil
//}
//
//func (m *mockMongoDBConfigRepo) FindAll() ([]entity.Mongodb, error) {
//	return []entity.Mongodb{{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"}}, nil
//}
//
//func (m *mockMongoDBConfigRepo) Update(config *entity.Mongodb) (string, error) {
//	return "OK", nil
//}
//
//func (m *mockMongoDBConfigRepo) Save(config *entity.Mongodb) (string, error) {
//	return "OK", nil
//}
//
//func (m *mockMongoDBConfigRepo) Delete(configName string) (string, error) {
//	return "OK", nil
//}
//
//type mockErrorMongoDBConfigRepo struct {
//}
//
//func (m *mockErrorMongoDBConfigRepo) Find(configName string) (*entity.Mongodb, error) {
//	return nil, errors.New("error from database querying")
//}
//
//func (m *mockErrorMongoDBConfigRepo) FindAll() ([]entity.Mongodb, error) {
//	return nil, errors.New("error from database querying")
//}
//
//func (m *mockErrorMongoDBConfigRepo) Update(config *entity.Mongodb) (string, error) {
//	return "", errors.New("error from database querying")
//}
//
//func (m *mockErrorMongoDBConfigRepo) Save(config *entity.Mongodb) (string, error) {
//	return "", errors.New("error from database querying")
//}
//func (m *mockErrorMongoDBConfigRepo) Delete(configName string) (string, error) {
//	return "", errors.New("error from database querying")
//}
//
//type mockTsConfigRepo struct {
//}
//
//func (m *mockTsConfigRepo) Find(configName string) (*entity.Tsconfig, error) {
//	return &entity.Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}, nil
//}
//
//func (m *mockTsConfigRepo) FindAll() ([]entity.Tsconfig, error) {
//	return []entity.Tsconfig{{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}}, nil
//}
//
//func (m *mockTsConfigRepo) Update(config *entity.Tsconfig) (string, error) {
//	return "OK", nil
//}
//
//func (m *mockTsConfigRepo) Save(config *entity.Tsconfig) (string, error) {
//	return "OK", nil
//}
//
//func (m *mockTsConfigRepo) Delete(configName string) (string, error) {
//	return "OK", nil
//}
//
//type mockErrorTsConfigRepo struct {
//}
//
//func (m *mockErrorTsConfigRepo) Find(configName string) (*entity.Tsconfig, error) {
//	return nil, errors.New("error from database querying")
//}
//
//func (m *mockErrorTsConfigRepo) FindAll() ([]entity.Tsconfig, error) {
//	return nil, errors.New("error from database querying")
//}
//
//func (m *mockErrorTsConfigRepo) Update(config *entity.Tsconfig) (string, error) {
//	return "", errors.New("error from database querying")
//}
//
//func (m *mockErrorTsConfigRepo) Save(config *entity.Tsconfig) (string, error) {
//	return "", errors.New("error from database querying")
//}
//func (m *mockErrorTsConfigRepo) Delete(configName string) (string, error) {
//	return "", errors.New("error from database querying")
//}
//
//type mockTempConfigRepo struct {
//}
//
//func (m *mockTempConfigRepo) Find(configName string) (*entity.Tempconfig, error) {
//	return &entity.Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}, nil
//}
//
//func (m *mockTempConfigRepo) FindAll() ([]entity.Tempconfig, error) {
//	return []entity.Tempconfig{{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}}, nil
//}
//
//func (m *mockTempConfigRepo) Update(config *entity.Tempconfig) (string, error) {
//	return "OK", nil
//}
//
//func (m *mockTempConfigRepo) Save(config *entity.Tempconfig) (string, error) {
//	return "OK", nil
//}
//
//func (m *mockTempConfigRepo) Delete(configName string) (string, error) {
//	return "OK", nil
//}
//
//type mockErrorTempConfigRepo struct {
//}
//
//func (m *mockErrorTempConfigRepo) Find(configName string) (*entity.Tempconfig, error) {
//	return nil, errors.New("error from database querying")
//}
//
//func (m *mockErrorTempConfigRepo) FindAll() ([]entity.Tempconfig, error) {
//	return nil, errors.New("error from database querying")
//}
//
//func (m *mockErrorTempConfigRepo) Update(config *entity.Tempconfig) (string, error) {
//	return "", errors.New("error from database querying")
//}
//
//func (m *mockErrorTempConfigRepo) Save(config *entity.Tempconfig) (string, error) {
//	return "", errors.New("error from database querying")
//}
//func (m *mockErrorTempConfigRepo) Delete(configName string) (string, error) {
//	return "", errors.New("error from database querying")
//}

func TestGetConfigByName(t *testing.T) {
	configCache := cache.New(testDefaultExpirationTimeOfCacheMin*time.Minute, testCleanupInternalOfCacheMin*time.Minute)
	mock := &configServerImpl{}
	mock.configCache = configCache
	mock.repo = &mockPostgresStorage{}

	res, err := mock.GetConfigByName("testNameMongo", "mongodb")
	if err != nil {
		t.Error("error during unit testing: ", err)
	}

	var expectedConfig []byte
	expectedConfig, err = json.Marshal(entity.Mongodb{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expectedConfig, res.Config)

	res, err = mock.GetConfigByName("testNameTs", "tsconfig")
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	expectedConfig, err = json.Marshal(entity.Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expectedConfig, res.Config)

	res, err = mock.GetConfigByName("testNameTemp", "tempconfig")
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	expectedConfig, err = json.Marshal(entity.Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expectedConfig, res.Config)

	mock.configCache.Flush()

	// testing error cases
	mock.repo = &mockPostgresStorageErr{}
	expectedError := errors.New("error from database querying")
	_, err = mock.GetConfigByName("testNameMongo", "mongodb")
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
	_, err = mock.GetConfigByName("testNameTs", "tsconfig")
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
	_, err = mock.GetConfigByName("testNameTemp", "tempconfig")
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
	_, err = mock.GetConfigByName("testNameTemp", "unexpectedConfigType")
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("unexpected type unexpectedConfigType"), err)
	}
}

func TestGetConfigByName_FromCache(t *testing.T) {
	testName := "testName"
	testConf := entity.Mongodb{Domain: testName, Mongodb: true, Host: "testHost", Port: "testPort"}
	configCache := cache.New(testDefaultExpirationTimeOfCacheMin*time.Minute, testCleanupInternalOfCacheMin*time.Minute)
	mock := &configServerImpl{}
	mock.configCache = configCache

	byteRes, err := json.Marshal(testConf)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	configResponse := &pb.GetConfigResponce{Config: byteRes}
	mock.configCache.Set(testName, configResponse, testDefaultExpirationTimeOfCacheMin*time.Minute)
	res, err := mock.GetConfigByName("testName", "mongodb")
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	var expectedConfig []byte
	expectedConfig, err = json.Marshal(entity.Mongodb{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expectedConfig, res.Config)
}

//func TestGetConfigsByType(t *testing.T) {
//
//	mock := &mockConfigServerImpl{}
//	err := mock.GetConfigsByType("mongodb", mock)
//	assert.Equal(t, 1, len(mock.Results), "expected to contain 1 item")
//	mock.tsConfigRepo = &mockTsConfigRepo{}
//	err = mock.GetConfigsByType("tsconfig", mock)
//	assert.Equal(t, 2, len(mock.Results), "expected to contain 1 item")
//	mock.tempConfigRepo = &mockTempConfigRepo{}
//	err = mock.GetConfigsByType("tempconfig", mock)
//	assert.Equal(t, 3, len(mock.Results), "expected to contain 1 item")
//	if err != nil {
//		t.Error("error during unit testing of GetConfigsByType function: ", err)
//	}
//
//	// testing error cases
//	err = mock.GetConfigsByType("unexpectedConfigType", mock)
//	if assert.Error(t, err) {
//		assert.Equal(t, errors.New("unexpected type"), err)
//	}
//
//	expectedError := errors.New("error from database querying")
//	err = nil
//	mock.mongoDBConfigRepo = &mockErrorMongoDBConfigRepo{}
//	err = mock.GetConfigsByType("mongodb", mock)
//	if assert.Error(t, err) {
//		assert.Equal(t, expectedError, err)
//	}
//
//	err = nil
//	mock.tsConfigRepo = &mockErrorTsConfigRepo{}
//	err = mock.GetConfigsByType("tsconfig", mock)
//	if assert.Error(t, err) {
//		assert.Equal(t, errors.New("error from database querying"), err)
//	}
//	err = nil
//	mock.tempConfigRepo = &mockErrorTempConfigRepo{}
//	err = mock.GetConfigsByType("tempconfig", mock)
//	if assert.Error(t, err) {
//		assert.Equal(t, expectedError, err)
//	}
//	err = nil
//	mock.tempConfigRepo = &mockErrorTempConfigRepo{}
//	err = mock.GetConfigsByType("unexpectedType", mock)
//	if assert.Error(t, err) {
//		assert.Equal(t, errors.New("unexpected type"), err)
//	}
//
//}

type mockConfigServerImpl struct {
	ConfigServer
	grpc.ServerStream
	Results []*pb.GetConfigResponce
}

//type mockConfigServer struct {
//	mongoDBConfigRepo repository.MongoDBConfigRepo
//	tsConfigRepo      repository.TsConfigRepo
//	tempConfigRepo    repository.TempConfigRepo
//	ConfigServer
//	grpc.ServerStream
//	Results     []*pb.GetConfigResponce
//	configCache *cache.Cache
//}
//
//func (mcs *mockConfigServer) Send(response *pb.GetConfigResponce) error {
//	mcs.Results = append(mcs.Results, response)
//	return nil
//}

//func TestInitServiceConfiguration(t *testing.T) {
//	os.Setenv("SERVICE_PORT", "")
//	os.Setenv("CACHE_EXPIRATION_TIME", "test")
//	os.Setenv("CACHE_CLEANUP_INTERVAL", "test")
//	expectedOut := serviceConfiguration{port: defaultPort, cacheExpirationTime: defaultCacheExpirationTime, cacheCleanupInterval: defaultCacheCleanupInterval}
//	realOutput := initServiceConfiguration()
//	assert.Equal(t, &expectedOut, realOutput)
//}

func TestCreateConfig(t *testing.T) {

	configCache := cache.New(testDefaultExpirationTimeOfCacheMin*time.Minute, testCleanupInternalOfCacheMin*time.Minute)

	mock := &configServerImpl{}
	mock.configCache = configCache
	mock.repo = &mockPostgresStorage{}

	testConfMongo := entity.Mongodb{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"}
	byteResMongo, err := json.Marshal(testConfMongo)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}

	res, err := mock.CreateConfig(&pb.Config{ConfigType: "mongodb", Config: byteResMongo})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	expectedResponse := &pb.Responce{Status: "OK"}
	assert.Equal(t, expectedResponse, res)

	testConfTs := entity.Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}
	byteResTs, err := json.Marshal(testConfTs)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}

	res, err = mock.CreateConfig(&pb.Config{ConfigType: "tsconfig", Config: byteResTs})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expectedResponse, res)

	testConfTemp := entity.Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	byteResTemp, err := json.Marshal(testConfTemp)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}

	res, err = mock.CreateConfig(&pb.Config{ConfigType: "tempconfig", Config: byteResTemp})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expectedResponse, res)

	// testing error cases
	mock.repo = &mockPostgresStorageErr{}

	// mongodb validation error
	expectedError := validator.ErrorMap{"Domain": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Host": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Port": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}}
	testConfMongoEmpty := entity.Mongodb{Domain: "", Mongodb: false, Host: "", Port: ""}
	byteResMongoEmpty, err := json.Marshal(testConfMongoEmpty)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	_, resultingErr := mock.CreateConfig(&pb.Config{ConfigType: "mongodb", Config: byteResMongoEmpty})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expectedError, resultingErr)
	}

	// mongodb saving error
	expError := errors.New("error from database querying")
	resultingErr = nil
	_, resultingErr = mock.CreateConfig(&pb.Config{ConfigType: "mongodb", Config: byteResMongo})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expError, resultingErr)
	}

	// ts validation error
	resultingErr = nil
	testConfTsEmpty := entity.Tsconfig{Excluding: 0, Target: "", Module: "", SourceMap: false}
	byteResTsEmpty, err := json.Marshal(testConfTsEmpty)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}

	expectedError = validator.ErrorMap{"Excluding": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Module": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Target": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}}
	_, resultingErr = mock.CreateConfig(&pb.Config{ConfigType: "tsconfig", Config: byteResTsEmpty})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expectedError, resultingErr)
	}

	// ts saving error
	expError = errors.New("error from database querying")
	_, resultingErr = mock.CreateConfig(&pb.Config{ConfigType: "tsconfig", Config: byteResTs})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expError, resultingErr)
	}

	// temp validation error
	resultingErr = nil
	testConfTempEmpty := entity.Tempconfig{Host: "", Port: "", Remoting: "", LegasyExplorer: false, RestApiRoot: ""}
	byteResTempEmpty, err := json.Marshal(testConfTempEmpty)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}

	expectedError = validator.ErrorMap{"RestApiRoot": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Host": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Port": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Remoting": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}}
	_, resultingErr = mock.CreateConfig(&pb.Config{ConfigType: "tempconfig", Config: byteResTempEmpty})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expectedError, resultingErr)
	}

	// temp saving error
	resultingErr = nil
	expError = errors.New("error from database querying")
	_, resultingErr = mock.CreateConfig(&pb.Config{ConfigType: "tempconfig", Config: byteResTemp})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expError, resultingErr)
	}

	// unexpectedType error
	resultingErr = nil
	_, resultingErr = mock.CreateConfig(&pb.Config{ConfigType: "unexpectedType", Config: byteResMongo})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, errors.New("unexpected type"), resultingErr)
	}
}

//func TestDeleteConfig(t *testing.T) {
//
//	configCache := cache.New(testDefaultExpirationTimeOfCacheMin*time.Minute, testCleanupInternalOfCacheMin*time.Minute)
//	mock := &mockConfigServer{}
//	mock.configCache = configCache
//	mock.mongoDBConfigRepo = &mockMongoDBConfigRepo{}
//	mock.tsConfigRepo = &mockTsConfigRepo{}
//	mock.tempConfigRepo = &mockTempConfigRepo{}
//
//	res, err := mock.DeleteConfig(context.Background(), &pb.DeleteConfigRequest{ConfigType: "mongodb", ConfigName: "testName"})
//	if err != nil {
//		t.Error("error during unit testing: ", err)
//	}
//	expectedResponse := &pb.Responce{Status: "OK"}
//	assert.Equal(t, expectedResponse, res)
//
//	res, err = mock.DeleteConfig(context.Background(), &pb.DeleteConfigRequest{ConfigType: "tsconfig", ConfigName: "testName"})
//	if err != nil {
//		t.Error("error during unit testing: ", err)
//	}
//	assert.Equal(t, expectedResponse, res)
//
//	res, err = mock.DeleteConfig(context.Background(), &pb.DeleteConfigRequest{ConfigType: "tempconfig", ConfigName: "testName"})
//	if err != nil {
//		t.Error("error during unit testing: ", err)
//	}
//	assert.Equal(t, expectedResponse, res)
//
//	// testing error cases
//	mock.mongoDBConfigRepo = &mockErrorMongoDBConfigRepo{}
//	mock.tsConfigRepo = &mockErrorTsConfigRepo{}
//	mock.tempConfigRepo = &mockErrorTempConfigRepo{}
//	expectedError := errors.New("error from database querying")
//	_, resultingErr := mock.DeleteConfig(context.Background(), &pb.DeleteConfigRequest{ConfigType: "mongodb", ConfigName: "errorTestName"})
//	if assert.Error(t, resultingErr) {
//		assert.Equal(t, expectedError, resultingErr)
//	}
//	resultingErr = nil
//	_, resultingErr = mock.DeleteConfig(context.Background(), &pb.DeleteConfigRequest{ConfigType: "tsconfig", ConfigName: "errorTestName"})
//	if assert.Error(t, resultingErr) {
//		assert.Equal(t, expectedError, resultingErr)
//	}
//	resultingErr = nil
//	_, resultingErr = mock.DeleteConfig(context.Background(), &pb.DeleteConfigRequest{ConfigType: "tempconfig", ConfigName: "errorTestName"})
//	if assert.Error(t, resultingErr) {
//		assert.Equal(t, expectedError, resultingErr)
//	}
//	resultingErr = nil
//	_, resultingErr = mock.DeleteConfig(context.Background(), &pb.DeleteConfigRequest{ConfigType: "unexpectedType", ConfigName: "errorTestName"})
//	if assert.Error(t, resultingErr) {
//		assert.Equal(t, errors.New("unexpected type"), resultingErr)
//	}
//}
//
//func TestUpdateConfig(t *testing.T) {
//
//	configCache := cache.New(testDefaultExpirationTimeOfCacheMin*time.Minute, testCleanupInternalOfCacheMin*time.Minute)
//	mock := &mockConfigServer{}
//	mock.configCache = configCache
//	mock.mongoDBConfigRepo = &mockMongoDBConfigRepo{}
//	mock.tsConfigRepo = &mockTsConfigRepo{}
//	mock.tempConfigRepo = &mockTempConfigRepo{}
//
//	testConfMongo := entity.Mongodb{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"}
//	byteResMongo, err := json.Marshal(testConfMongo)
//	if err != nil {
//		t.Error("error during unit testing: ", err)
//	}
//	testConfTs := entity.Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}
//	byteResTs, err := json.Marshal(testConfTs)
//	if err != nil {
//		t.Error("error during unit testing: ", err)
//	}
//	testConfTemp := entity.Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
//	byteResTemp, err := json.Marshal(testConfTemp)
//	if err != nil {
//		t.Error("error during unit testing: ", err)
//	}
//
//	resp, err := mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "mongodb", Config: byteResMongo})
//	assert.Equal(t, &pb.Responce{Status: "OK"}, resp)
//	resp, err = mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "tsconfig", Config: byteResTs})
//	assert.Equal(t, &pb.Responce{Status: "OK"}, resp)
//	resp, err = mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "tempconfig", Config: byteResTemp})
//	assert.Equal(t, &pb.Responce{Status: "OK"}, resp)
//
//	// testing error cases
//	// testing unexpected type
//	_, err = mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "unexpectedConfigType"})
//	if assert.Error(t, err) {
//		assert.Equal(t, errors.New("unexpected type"), err)
//	}
//
//	// mongodb validation error
//	expectedError := validator.ErrorMap{"Domain": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Host": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Port": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}}
//	testConfMongoEmpty := entity.Mongodb{Domain: "", Mongodb: false, Host: "", Port: ""}
//	byteResMongoEmpty, err := json.Marshal(testConfMongoEmpty)
//	if err != nil {
//		t.Error("error during unit testing: ", err)
//	}
//	_, resultingErr := mock.CreateConfig(context.Background(), &pb.Config{ConfigType: "mongodb", Config: byteResMongoEmpty})
//	if assert.Error(t, resultingErr) {
//		assert.Equal(t, expectedError, resultingErr)
//	}
//
//	// mongodb saving error
//	expError := errors.New("error from database querying")
//	resultingErr = nil
//	mock.mongoDBConfigRepo = &mockErrorMongoDBConfigRepo{}
//	err = nil
//	_, err = mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "mongodb", Config: byteResMongo})
//	if assert.Error(t, err) {
//		assert.Equal(t, expError, err)
//	}
//
//	// ts validation error
//	resultingErr = nil
//	testConfTsEmpty := entity.Tsconfig{Excluding: 0, Target: "", Module: "", SourceMap: false}
//	byteResTsEmpty, err := json.Marshal(testConfTsEmpty)
//	if err != nil {
//		t.Error("error during unit testing: ", err)
//	}
//	expectedError = validator.ErrorMap{"Excluding": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Module": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Target": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}}
//	_, resultingErr = mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "tsconfig", Config: byteResTsEmpty})
//	if assert.Error(t, resultingErr) {
//		assert.Equal(t, expectedError, resultingErr)
//	}
//
//	// ts saving error
//	expError = errors.New("error from database querying")
//	err = nil
//	mock.tsConfigRepo = &mockErrorTsConfigRepo{}
//	_, err = mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "tsconfig", Config: byteResTs})
//	if assert.Error(t, err) {
//		assert.Equal(t, expError, err)
//	}
//
//	// temp validation error
//	resultingErr = nil
//	testConfTempEmpty := entity.Tempconfig{Host: "", Port: "", Remoting: "", LegasyExplorer: false, RestApiRoot: ""}
//	byteResTempEmpty, err := json.Marshal(testConfTempEmpty)
//	if err != nil {
//		t.Error("error during unit testing: ", err)
//	}
//
//	expectedError = validator.ErrorMap{"RestApiRoot": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Host": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Port": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Remoting": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}}
//	_, resultingErr = mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "tempconfig", Config: byteResTempEmpty})
//	if assert.Error(t, resultingErr) {
//		assert.Equal(t, expectedError, resultingErr)
//	}
//
//	// temp saving error
//	resultingErr = nil
//	expError = errors.New("error from database querying")
//	err = nil
//	mock.tempConfigRepo = &mockErrorTempConfigRepo{}
//	_, err = mock.UpdateConfig(context.Background(), &pb.Config{ConfigType: "tempconfig", Config: byteResTemp})
//	if assert.Error(t, err) {
//		assert.Equal(t, expError, err)
//	}
//}
//
//func BenchmarkCreateConfig(b *testing.B) {
//
//	configCache := cache.New(testDefaultExpirationTimeOfCacheMin*time.Minute, testCleanupInternalOfCacheMin*time.Minute)
//	mock := &mockConfigServer{}
//	mock.configCache = configCache
//	mock.mongoDBConfigRepo = &mockMongoDBConfigRepo{}
//	mock.tsConfigRepo = &mockTsConfigRepo{}
//	mock.tempConfigRepo = &mockTempConfigRepo{}
//
//	testConfMongo := entity.Mongodb{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"}
//	byteResMongo, err := json.Marshal(testConfMongo)
//	if err != nil {
//		b.Error("error during unit testing: ", err)
//	}
//	//b.ReportAllocs()
//	for i := 0; i < b.N; i++ {
//		_, err := mock.CreateConfig(context.Background(), &pb.Config{ConfigType: "mongodb", Config: byteResMongo})
//		if err != nil {
//			b.Error("error during unit testing: ", err)
//		}
//	}
//
//}
