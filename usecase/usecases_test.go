package usecase

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/YAWAL/GetMeConf/entity"
	pb "github.com/YAWAL/GetMeConfAPI/api"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"gopkg.in/validator.v2"
)

const (
	testDefaultExpirationTimeOfCacheMin = 5
	testCleanupInternalOfCacheMin       = 10
)

var errDB = errors.New("database error")

type mockPostgresStorage struct{}

func (m *mockPostgresStorage) Migrate() error {
	return nil
}

type configServerImplWrap struct {
	i *ConfigServerImpl
	ConfigServer
	grpc.ServerStream
	Results []*pb.GetConfigResponce
}

func (mcs *configServerImplWrap) Send(response *pb.GetConfigResponce) error {
	mcs.Results = append(mcs.Results, response)
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
	return nil, errDB
}
func (m *mockPostgresStorageErr) FindAllMongoDBConfig() ([]entity.Mongodb, error) {
	return nil, errDB
}
func (m *mockPostgresStorageErr) UpdateMongoDBConfig(config *entity.Mongodb) (string, error) {
	return "", errDB
}
func (m *mockPostgresStorageErr) SaveMongoDBConfig(config *entity.Mongodb) (string, error) {
	return "", errDB
}
func (m *mockPostgresStorageErr) DeleteMongoDBConfig(configName string) (string, error) {
	return "", errDB
}

func (m *mockPostgresStorageErr) FindTempConfig(configName string) (*entity.Tempconfig, error) {
	return nil, errDB
}
func (m *mockPostgresStorageErr) FindAllTempConfig() ([]entity.Tempconfig, error) {
	return nil, errDB
}
func (m *mockPostgresStorageErr) UpdateTempConfig(config *entity.Tempconfig) (string, error) {
	return "", errDB
}
func (m *mockPostgresStorageErr) SaveTempConfig(config *entity.Tempconfig) (string, error) {
	return "", errDB
}
func (m *mockPostgresStorageErr) DeleteTempConfig(configName string) (string, error) {
	return "", errDB
}

func (m *mockPostgresStorageErr) FindTsConfig(configName string) (*entity.Tsconfig, error) {
	return nil, errDB
}
func (m *mockPostgresStorageErr) FindAllTsConfig() ([]entity.Tsconfig, error) {
	return nil, errDB
}
func (m *mockPostgresStorageErr) UpdateTsConfig(config *entity.Tsconfig) (string, error) {
	return "", errDB
}
func (m *mockPostgresStorageErr) SaveTsConfig(config *entity.Tsconfig) (string, error) {
	return "", errDB
}
func (m *mockPostgresStorageErr) DeleteTsConfig(configName string) (string, error) {
	return "", errDB
}

func TestGetConfigByName(t *testing.T) {
	configCache := cache.New(testDefaultExpirationTimeOfCacheMin*time.Minute, testCleanupInternalOfCacheMin*time.Minute)
	mock := &ConfigServerImpl{}
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
	expectedError := errDB
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
	mock := &ConfigServerImpl{}
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

func TestGetConfigsByType(t *testing.T) {

	configCache := cache.New(testDefaultExpirationTimeOfCacheMin*time.Minute, testCleanupInternalOfCacheMin*time.Minute)

	mock := &ConfigServerImpl{}
	mock.configCache = configCache
	mock.repo = &mockPostgresStorage{}

	m := configServerImplWrap{i: mock}

	err := mock.GetConfigsByType(mongodb, &m)
	if err != nil {
		t.Error("error during unit testing of GetConfigsByType function: ", err)
	}
	assert.Equal(t, 1, len(m.Results), "expected to contain 1 item")
	if err != nil {
		t.Error("error during unit testing of GetConfigsByType function: ", err)
	}
	err = mock.GetConfigsByType(tsconfig, &m)
	if err != nil {
		t.Error("error during unit testing of GetConfigsByType function: ", err)
	}
	assert.Equal(t, 2, len(m.Results), "expected to contain 1 item")
	err = mock.GetConfigsByType(tempconfig, &m)
	assert.Equal(t, 3, len(m.Results), "expected to contain 1 item")
	if err != nil {
		t.Error("error during unit testing of GetConfigsByType function: ", err)
	}

	// testing error cases
	err = mock.GetConfigsByType("unexpectedConfigType", &m)
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("unexpected type unexpectedConfigType"), err)
	}

	expectedError := errDB
	err = nil
	mock.repo = &mockPostgresStorageErr{}
	err = mock.GetConfigsByType(mongodb, &m)
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
	err = nil
	err = mock.GetConfigsByType(tsconfig, &m)
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
	err = nil
	err = mock.GetConfigsByType(tempconfig, &m)
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
}

func TestCreateConfig(t *testing.T) {
	configCache := cache.New(testDefaultExpirationTimeOfCacheMin*time.Minute, testCleanupInternalOfCacheMin*time.Minute)
	mock := &ConfigServerImpl{}
	mock.configCache = configCache
	mock.repo = &mockPostgresStorage{}
	testConfMongo := entity.Mongodb{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"}
	byteResMongo, err := json.Marshal(testConfMongo)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}

	res, err := mock.CreateConfig(&pb.Config{ConfigType: mongodb, Config: byteResMongo})
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

	res, err = mock.CreateConfig(&pb.Config{ConfigType: tsconfig, Config: byteResTs})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expectedResponse, res)

	testConfTemp := entity.Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	byteResTemp, err := json.Marshal(testConfTemp)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}

	res, err = mock.CreateConfig(&pb.Config{ConfigType: tempconfig, Config: byteResTemp})
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
	_, resultingErr := mock.CreateConfig(&pb.Config{ConfigType: mongodb, Config: byteResMongoEmpty})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expectedError, resultingErr)
	}

	// ts validation error
	resultingErr = nil
	testConfTsEmpty := entity.Tsconfig{Excluding: 0, Target: "", Module: "", SourceMap: false}
	byteResTsEmpty, err := json.Marshal(testConfTsEmpty)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}

	expectedError = validator.ErrorMap{"Module": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Target": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}}
	_, resultingErr = mock.CreateConfig(&pb.Config{ConfigType: tsconfig, Config: byteResTsEmpty})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expectedError, resultingErr)
	}

	//	temp validation error
	resultingErr = nil
	testConfTempEmpty := entity.Tempconfig{Host: "", Port: "", Remoting: "", LegasyExplorer: false, RestApiRoot: ""}
	byteResTempEmpty, err := json.Marshal(testConfTempEmpty)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}

	expectedError = validator.ErrorMap{"RestApiRoot": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Host": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Port": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Remoting": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}}
	_, resultingErr = mock.CreateConfig(&pb.Config{ConfigType: tempconfig, Config: byteResTempEmpty})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expectedError, resultingErr)
	}

	// mongodb saving error
	expError := errDB
	resultingErr = nil
	_, resultingErr = mock.CreateConfig(&pb.Config{ConfigType: mongodb, Config: byteResMongo})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expError, resultingErr)
	}

	// temp saving error
	resultingErr = nil
	_, resultingErr = mock.CreateConfig(&pb.Config{ConfigType: "tempconfig", Config: byteResTemp})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expError, resultingErr)
	}

	// ts saving error
	_, resultingErr = mock.CreateConfig(&pb.Config{ConfigType: tsconfig, Config: byteResTs})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expError, resultingErr)
	}

	// unexpectedType error
	resultingErr = nil
	_, resultingErr = mock.CreateConfig(&pb.Config{ConfigType: "anotherType", Config: byteResMongo})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, errors.New("unexpected type anotherType"), resultingErr)
	}
}

func TestDeleteConfig(t *testing.T) {

	configCache := cache.New(testDefaultExpirationTimeOfCacheMin*time.Minute, testCleanupInternalOfCacheMin*time.Minute)
	mock := &ConfigServerImpl{}
	mock.configCache = configCache
	mock.repo = &mockPostgresStorage{}

	res, err := mock.DeleteConfig(&pb.DeleteConfigRequest{ConfigType: mongodb, ConfigName: "testName"})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	expectedResponse := &pb.Responce{Status: "OK"}
	assert.Equal(t, expectedResponse, res)

	res, err = mock.DeleteConfig(&pb.DeleteConfigRequest{ConfigType: tsconfig, ConfigName: "testName"})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expectedResponse, res)

	res, err = mock.DeleteConfig(&pb.DeleteConfigRequest{ConfigType: tempconfig, ConfigName: "testName"})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expectedResponse, res)

	// testing error cases
	mock.repo = &mockPostgresStorageErr{}
	expectedError := errDB
	_, resultingErr := mock.DeleteConfig(&pb.DeleteConfigRequest{ConfigType: mongodb, ConfigName: "errorTestName"})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expectedError, resultingErr)
	}
	resultingErr = nil
	_, resultingErr = mock.DeleteConfig(&pb.DeleteConfigRequest{ConfigType: tsconfig, ConfigName: "errorTestName"})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expectedError, resultingErr)
	}
	resultingErr = nil
	_, resultingErr = mock.DeleteConfig(&pb.DeleteConfigRequest{ConfigType: tempconfig, ConfigName: "errorTestName"})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expectedError, resultingErr)
	}
	resultingErr = nil
	_, resultingErr = mock.DeleteConfig(&pb.DeleteConfigRequest{ConfigType: "errorTestType", ConfigName: "errorTestName"})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, errors.New("unexpected type errorTestType"), resultingErr)
	}
}

func TestUpdateConfig(t *testing.T) {

	configCache := cache.New(testDefaultExpirationTimeOfCacheMin*time.Minute, testCleanupInternalOfCacheMin*time.Minute)
	mock := &ConfigServerImpl{}
	mock.configCache = configCache
	mock.repo = &mockPostgresStorage{}

	testConfMongo := entity.Mongodb{Domain: "testName", Mongodb: true, Host: "testHost", Port: "testPort"}
	byteResMongo, err := json.Marshal(testConfMongo)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	testConfTs := entity.Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}
	byteResTs, err := json.Marshal(testConfTs)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	testConfTemp := entity.Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	byteResTemp, err := json.Marshal(testConfTemp)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}

	resp, err := mock.UpdateConfig(&pb.Config{ConfigType: mongodb, Config: byteResMongo})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, &pb.Responce{Status: "OK"}, resp)
	resp, err = mock.UpdateConfig(&pb.Config{ConfigType: tsconfig, Config: byteResTs})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, &pb.Responce{Status: "OK"}, resp)
	resp, err = mock.UpdateConfig(&pb.Config{ConfigType: tempconfig, Config: byteResTemp})
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, &pb.Responce{Status: "OK"}, resp)

	// testing error cases
	// testing unexpected type
	_, err = mock.UpdateConfig(&pb.Config{ConfigType: "unexpectedConfigType"})
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("unexpected type unexpectedConfigType"), err)
	}

	// mongodb validation error
	expectedError := validator.ErrorMap{"Domain": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Host": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Port": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}}
	testConfMongoEmpty := entity.Mongodb{Domain: "", Mongodb: false, Host: "", Port: ""}
	byteResMongoEmpty, err := json.Marshal(testConfMongoEmpty)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	_, resultingErr := mock.UpdateConfig(&pb.Config{ConfigType: mongodb, Config: byteResMongoEmpty})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expectedError, resultingErr)
	}

	// ts validation error
	resultingErr = nil
	testConfTsEmpty := entity.Tsconfig{Excluding: 0, Target: "", Module: "", SourceMap: false}
	byteResTsEmpty, err := json.Marshal(testConfTsEmpty)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	expectedError = validator.ErrorMap{"Module": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Target": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}}
	_, resultingErr = mock.UpdateConfig(&pb.Config{ConfigType: tsconfig, Config: byteResTsEmpty})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expectedError, resultingErr)
	}

	// temp validation error
	resultingErr = nil
	testConfTempEmpty := entity.Tempconfig{Host: "", Port: "", Remoting: "", LegasyExplorer: false, RestApiRoot: ""}
	byteResTempEmpty, err := json.Marshal(testConfTempEmpty)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}

	expectedError = validator.ErrorMap{"RestApiRoot": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Host": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Port": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}, "Remoting": validator.ErrorArray{validator.TextErr{Err: errors.New("zero value")}}}
	_, resultingErr = mock.UpdateConfig(&pb.Config{ConfigType: tempconfig, Config: byteResTempEmpty})
	if assert.Error(t, resultingErr) {
		assert.Equal(t, expectedError, resultingErr)
	}

	// mongodb saving error
	expError := errDB
	resultingErr = nil
	mock.repo = &mockPostgresStorageErr{}
	err = nil
	_, err = mock.UpdateConfig(&pb.Config{ConfigType: mongodb, Config: byteResMongo})
	if assert.Error(t, err) {
		assert.Equal(t, expError, err)
	}

	// ts saving error
	err = nil
	_, err = mock.UpdateConfig(&pb.Config{ConfigType: tsconfig, Config: byteResTs})
	if assert.Error(t, err) {
		assert.Equal(t, expError, err)
	}

	// temp saving error
	err = nil
	_, err = mock.UpdateConfig(&pb.Config{ConfigType: tempconfig, Config: byteResTemp})
	if assert.Error(t, err) {
		assert.Equal(t, expError, err)
	}
}

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
