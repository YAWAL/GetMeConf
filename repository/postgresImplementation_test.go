// Package repository contains repository interfaces as well as their implementations for given databases
package repository

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/YAWAL/GetMeConf/entity"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var logger, _ = zap.NewDevelopment()

var errDB = errors.New("database error")

var errDBDelete = errors.New("could not delete from database")

func newDB() (sqlmock.Sqlmock, *gorm.DB, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal("can not create sql mock " + err.Error())
		return nil, nil, err
	}
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		logger.Fatal("can not open gorm connection " + err.Error())
		return nil, nil, err
	}
	gormDB.LogMode(true)
	return mock, gormDB, nil
}

func formatRequest(s string) string {
	return fmt.Sprintf("^%s$", regexp.QuoteMeta(s))
}

func TestPostgresStorage_FindMongoDBConfig(t *testing.T) {
	m, db, _ := newDB()
	repo := PostgresStorage{DB: db}
	mongodbConfig := entity.Mongodb{Domain: "testDomain", Mongodb: true, Host: "testHost", Port: "testPort"}
	mongoRows := getMongoDBRows(mongodbConfig.Domain)
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\" WHERE (domain = $1)")).WithArgs("testDomain").
		WillReturnRows(mongoRows)
	returnedMongoConfigs, err := repo.FindMongoDBConfig("testDomain")
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, &mongodbConfig, returnedMongoConfigs)

	configName := "notExistingConfig"
	expectedError := errDB
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\" WHERE (domain = $1)")).WithArgs("notExistingConfig").
		WillReturnError(expectedError)
	_, returnedErr := repo.FindMongoDBConfig(configName)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}
func TestPostgresStorage_FindTsConfig(t *testing.T) {
	m, db, _ := newDB()
	repo := PostgresStorage{DB: db}
	tsConfig := entity.Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}
	tsRows := getTsConfigRows(tsConfig.Module)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\" WHERE (module = $1)")).WithArgs("testModule").
		WillReturnRows(tsRows)
	returnedTsConfigs, err := repo.FindTsConfig("testModule")
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	fmt.Println(returnedTsConfigs)
	assert.Equal(t, &tsConfig, returnedTsConfigs)

	configName := "notExistingConfig"
	expectedError := errDB
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\" WHERE (module = $1)")).
		WithArgs("notExistingConfig").WillReturnError(expectedError)
	_, returnedErr := repo.FindTsConfig(configName)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}

func TestPostgresStorage_FindTempConfig(t *testing.T) {
	m, db, _ := newDB()
	repo := PostgresStorage{DB: db}
	tempConfig := entity.Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort",
		Remoting: "testRemoting", LegasyExplorer: true}
	tempRows := getTempConfigRows(tempConfig.RestApiRoot)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\" WHERE (rest_api_root = $1)")).
		WithArgs("testRestApiRoot").WillReturnRows(tempRows)
	returnedTempConfigs, err := repo.FindTempConfig("testRestApiRoot")
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	fmt.Println(returnedTempConfigs)
	assert.Equal(t, &tempConfig, returnedTempConfigs)

	configName := "notExistingConfig"
	expectedError := errDB
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\" WHERE (rest_api_root = $1)")).
		WithArgs("notExistingConfig").WillReturnError(expectedError)
	_, returnedErr := repo.FindTempConfig(configName)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}

}

func TestPostgresStorage_FindAllMongoDBConfigl(t *testing.T) {
	m, db, _ := newDB()
	repo := PostgresStorage{DB: db}
	mongodbConfig := entity.Mongodb{Domain: "testDomain", Mongodb: true, Host: "testHost", Port: "testPort"}
	mongoRows := getMongoDBRows(mongodbConfig.Domain)
	expConfigs := []entity.Mongodb{mongodbConfig}
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\"")).WillReturnRows(mongoRows)
	returnedMongoConfigs, err := repo.FindAllMongoDBConfig()
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expConfigs, returnedMongoConfigs)

	expectedError := errDB
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\"")).WillReturnError(expectedError)
	_, returnedErr := repo.FindAllMongoDBConfig()
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}

func TestPostgresStorage_FindAllTsConfig(t *testing.T) {
	m, db, _ := newDB()
	repo := PostgresStorage{DB: db}
	tsConfig := entity.Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}
	tsRows := getTsConfigRows(tsConfig.Module)
	expTsConfigs := []entity.Tsconfig{tsConfig}
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\"")).WillReturnRows(tsRows)
	returnedTsConfigs, err := repo.FindAllTsConfig()
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, expTsConfigs, returnedTsConfigs)

	expectedError := errDB
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\"")).WillReturnError(expectedError)
	_, returnedErr := repo.FindAllTsConfig()
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}
func TestPostgresStorage_FindAllTempConfig(t *testing.T) {
	m, db, _ := newDB()
	repo := PostgresStorage{DB: db}
	tempConfig := entity.Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort",
		Remoting: "testRemoting", LegasyExplorer: true}
	tempRows := getTempConfigRows(tempConfig.RestApiRoot)
	expTempConfigs := []entity.Tempconfig{tempConfig}
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\"")).WillReturnRows(tempRows)
	returnedTempConfigs, err := repo.FindAllTempConfig()
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	fmt.Println(returnedTempConfigs)
	assert.Equal(t, expTempConfigs, returnedTempConfigs)

	expectedError := errDB
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\"")).WillReturnError(expectedError)
	_, returnedErr := repo.FindAllTempConfig()
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}

func TestPostgresStorage_SaveMongoDBConfig(t *testing.T) {
	m, db, _ := newDB()
	mockRepo := PostgresStorage{DB: db}
	mongodbConfig := entity.Mongodb{Domain: "testDomain", Mongodb: true, Host: "testHost", Port: "testPort"}

	//rows := sqlmock.NewRows([]string{"testDomain", "true", "testHost", "testPort"})

	m.ExpectExec(formatRequest("INSERT INTO \"mongodbs\" (\"domain\",\"mongodb\",\"host\",\"port\") "+
		"VALUES ($1,$2,$3,$4) RETURNING \"mongodbs\".\"domain\"")).
		WithArgs("testDomain", true, "testHost", "testPort").WillReturnResult(sqlmock.NewResult(0, 1))
	//m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\" WHERE \"mongodbs\".\"domain\" = $1")).
	//	WithArgs("testDomain").WillReturnRows(rows)
	result, err := mockRepo.SaveMongoDBConfig(&mongodbConfig)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, "OK", result)

}

func TestPostgresStorage_SaveMongoDBConfigWithError(t *testing.T) {
	m, db, _ := newDB()
	mockRepo := PostgresStorage{DB: db}

	mongodbConfigErr := entity.Mongodb{Domain: "testDomainError", Mongodb: true, Host: "testHost", Port: "testPort"}
	expectedError := errors.New("db error")
	m.ExpectQuery(formatRequest("INSERT INTO \"mongodbs\" (\"domain\",\"mongodb\",\"host\",\"port\") "+
		"VALUES ($1,$2,$3,$4) RETURNING \"mongodbs\".\"domain\"")).
		WithArgs("testDomainError", true, "testHost", "testPort").
		WillReturnError(expectedError)
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\" WHERE \"mongodbs\".\"domain\" = $1")).
		WithArgs("testDomainError").WillReturnError(expectedError)
	_, returnedErr := mockRepo.SaveMongoDBConfig(&mongodbConfigErr)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}

func TestPostgresStorage_SaveTsConfig(t *testing.T) {
	m, db, _ := newDB()
	mockRepo := PostgresStorage{DB: db}

	tsConfig := entity.Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}

	rows := sqlmock.NewRows([]string{"testModule", "testTarget", "true", "1"})

	m.ExpectQuery(formatRequest("INSERT INTO \"tsconfigs\" (\"module\",\"target\",\"source_map\",\"excluding\") "+
		"VALUES ($1,$2,$3,$4) RETURNING \"tsconfigs\".\"module\"")).
		WithArgs("testModule", "testTarget", true, 1).
		WillReturnRows(rows)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\" WHERE \"tsconfigs\".\"module\" = $1")).
		WithArgs("testModule").WillReturnRows(rows)
	result, err := mockRepo.SaveTsConfig(&tsConfig)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, "OK", result)

}

func TestPostgresStorage_SaveTsConfigWithError(t *testing.T) {
	m, db, _ := newDB()
	mockRepo := PostgresStorage{DB: db}

	tsConfigErr := entity.Tsconfig{Module: "testModuleError", Target: "testTarget", SourceMap: true, Excluding: 1}
	expectedError := errors.New("db error")
	m.ExpectQuery(formatRequest("INSERT INTO \"tsconfigs\" (\"module\",\"target\",\"source_map\",\"excluding\") "+
		"VALUES ($1,$2,$3,$4) RETURNING \"tsconfigs\".\"module\"")).
		WithArgs("testModuleError", "testTarget", true, 1).
		WillReturnError(expectedError)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\" WHERE \"tsconfigs\".\"module\" = $1")).
		WithArgs("testModuleError").WillReturnError(expectedError)
	_, returnedErr := mockRepo.SaveTsConfig(&tsConfigErr)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}

func TestPostgresStorage_SaveTempConfig(t *testing.T) {
	m, db, _ := newDB()
	mockRepo := PostgresStorage{DB: db}

	tempConfig := entity.Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort",
		Remoting: "testRemoting", LegasyExplorer: true}
	rows := sqlmock.NewRows([]string{"testApiRoot", "testHost", "testPort", "testRemoting", "true"})

	m.ExpectQuery(formatRequest("INSERT INTO \"tempconfigs\" (\"rest_api_root\",\"host\",\"port\",\"remoting\",\"legasy_explorer\") "+
		"VALUES ($1,$2,$3,$4,$5) RETURNING \"tempconfigs\".\"rest_api_root\"")).
		WithArgs("testApiRoot", "testHost", "testPort", "testRemoting", true).
		WillReturnRows(rows)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\" WHERE \"tempconfigs\".\"rest_api_root\" = $1")).
		WithArgs("testApiRoot").WillReturnRows(rows)
	result, err := mockRepo.SaveTempConfig(&tempConfig)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, "OK", result)
}

func TestPostgresStorage_SaveTempConfigWithError(t *testing.T) {
	m, db, _ := newDB()
	mockRepo := PostgresStorage{DB: db}

	tempConfigErr := entity.Tempconfig{RestApiRoot: "testApiRootError", Host: "testHost", Port: "testPort",
		Remoting: "testRemoting", LegasyExplorer: true}
	expectedError := errors.New("db error")
	m.ExpectQuery(formatRequest("INSERT INTO \"tempconfigs\" (\"rest_api_root\",\"host\",\"port\",\"remoting\",\"legasy_explorer\") "+
		"VALUES ($1,$2,$3,$4,$5) RETURNING \"tempconfigs\".\"rest_api_root\"")).
		WithArgs("testApiRootError", "testHost", "testPort", "testRemoting", true).
		WillReturnError(expectedError)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\" WHERE \"tempconfigs\".\"rest_api_root\" = $1")).
		WithArgs("testApiRootError").WillReturnError(expectedError)
	_, returnedErr := mockRepo.SaveTempConfig(&tempConfigErr)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}

func TestPostgresStorage_DeleteMongoDBConfig(t *testing.T) {
	m, db, _ := newDB()
	repo := PostgresStorage{DB: db}
	testType := "mongodb"
	testID := "testID"
	m.ExpectExec(formatRequest("DELETE FROM \"" + testType + "s\" WHERE (domain = $1)")).
		WithArgs("testID").WillReturnResult(sqlmock.NewResult(0, 1))
	res, err := repo.DeleteMongoDBConfig(testID)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, "deleted 1 row(s)", res)

	testID = "notExistingTestID"
	m.ExpectExec(formatRequest("DELETE FROM \"" + testType + "s\" WHERE (domain = $1)")).
		WithArgs("notExistingTestID").WillReturnError(errDBDelete)
	_, returnedErr := repo.DeleteMongoDBConfig(testID)
	expectedError := errDBDelete
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}
func TestPostgresStorage_DeleteTsConfig(t *testing.T) {
	m, db, _ := newDB()
	repo := PostgresStorage{DB: db}
	testType := "tsconfig"
	testID := "testID"
	m.ExpectExec(formatRequest("DELETE FROM \"" + testType + "s\" WHERE (module = $1)")).
		WithArgs("testID").WillReturnResult(sqlmock.NewResult(0, 1))
	res, err := repo.DeleteTsConfig(testID)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, "deleted 1 row(s)", res)

	testID = "notExistingTestID"
	m.ExpectExec(formatRequest("DELETE FROM \"" + testType + "s\" WHERE (module = $1)")).
		WithArgs("notExistingTestID").WillReturnError(errDBDelete)
	_, returnedErr := repo.DeleteTsConfig(testID)
	expectedError := errDBDelete
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}

func TestPostgresStorage_DeleteTempConfig(t *testing.T) {
	m, db, _ := newDB()
	repo := PostgresStorage{DB: db}
	testType := "tempconfig"
	testID := "testID"
	m.ExpectExec(formatRequest("DELETE FROM \"" + testType + "s\" WHERE (rest_api_root = $1)")).
		WithArgs("testID").WillReturnResult(sqlmock.NewResult(0, 1))
	res, err := repo.DeleteTempConfig(testID)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, "deleted 1 row(s)", res)

	testID = "notExistingTestID"
	m.ExpectExec(formatRequest("DELETE FROM \"" + testType + "s\" WHERE (rest_api_root = $1)")).
		WithArgs("notExistingTestID").WillReturnError(errDBDelete)
	_, returnedErr := repo.DeleteTempConfig(testID)
	expectedError := errDBDelete
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedError, returnedErr)
	}
}

func TestPostgresStorage_UpdateMongoDBConfig(t *testing.T) {
	m, db, _ := newDB()
	repo := PostgresStorage{DB: db}
	config := entity.Mongodb{Domain: "testDomain", Mongodb: true, Host: "testHost", Port: "8080"}
	rows := getMongoDBRows(config.Domain)
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\" WHERE (domain = $1)")).
		WithArgs("testDomain").
		WillReturnRows(rows)
	m.ExpectExec(formatRequest("UPDATE mongodbs SET mongodb = $1, port = $2, host = $3 WHERE domain = $4")).
		WithArgs(strconv.FormatBool(config.Mongodb), config.Port, config.Host, config.Domain).
		WillReturnResult(sqlmock.NewResult(0, 1))
	result, err := repo.UpdateMongoDBConfig(&config)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, "OK", result)

	configErrOne := entity.Mongodb{Domain: "errOneConfig", Mongodb: true, Host: "testHost", Port: "8080"}
	_ = getMongoDBRows(configErrOne.Domain)
	expectedErrorOne := errors.New("record not found")
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\" WHERE (domain = $1)")).
		WithArgs("errOneConfig").
		WillReturnError(expectedErrorOne)
	_, returnedErr := repo.UpdateMongoDBConfig(&configErrOne)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedErrorOne, returnedErr)
	}

	expectedErrorTwo := errors.New("db error")
	configErrTwo := entity.Mongodb{Domain: "errTwoConfig", Mongodb: true, Host: "testHost", Port: "8080"}
	rows = getMongoDBRows(configErrTwo.Domain)
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\" WHERE (domain = $1)")).
		WithArgs("errTwoConfig").
		WillReturnRows(rows)
	m.ExpectExec(formatRequest("UPDATE mongodbs SET mongodb = $1, port = $2, host = $3 WHERE domain = $4")).
		WithArgs(strconv.FormatBool(configErrTwo.Mongodb), configErrTwo.Port, configErrTwo.Host, configErrTwo.Domain).
		WillReturnError(expectedErrorTwo)
	_, returnedErr = repo.UpdateMongoDBConfig(&configErrTwo)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedErrorTwo, returnedErr)
	}

	expectedErrorThree := errors.New("fields are empty")
	configErrThree := entity.Mongodb{Domain: "errThreeConfig", Mongodb: true, Host: "", Port: ""}
	rows = getMongoDBRows(configErrThree.Domain)
	m.ExpectQuery(formatRequest("SELECT * FROM \"mongodbs\" WHERE (domain = $1)")).
		WithArgs("errThreeConfig").
		WillReturnRows(rows)
	m.ExpectExec(formatRequest("UPDATE mongodbs SET mongodb = $1, port = $2, host = $3 WHERE domain = $4")).
		WithArgs(strconv.FormatBool(configErrThree.Mongodb), configErrThree.Port, configErrThree.Host, configErrThree.Domain).
		WillReturnError(expectedErrorThree)
	_, returnedErr = repo.UpdateMongoDBConfig(&configErrThree)
	if assert.Error(t, returnedErr) {
		assert.Equal(t, expectedErrorThree, returnedErr)
	}
}

func TestPostgresStorage_UpdateTsConfig(t *testing.T) {
	m, db, _ := newDB()

	repo := PostgresStorage{DB: db}
	tsConfig := entity.Tsconfig{Module: "testModule", Target: "testTarget", SourceMap: true, Excluding: 1}
	tsRows := getTsConfigRows(tsConfig.Module)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\" WHERE (module = $1)")).
		WithArgs("testModule").
		WillReturnRows(tsRows)
	m.ExpectExec(formatRequest("UPDATE tsconfigs SET target = $1, source_map = $2, excluding = $3 WHERE module = $4")).
		WithArgs(tsConfig.Target, strconv.FormatBool(tsConfig.SourceMap), strconv.Itoa(tsConfig.Excluding), tsConfig.Module).
		WillReturnResult(sqlmock.NewResult(0, 1))
	tsResult, err := repo.UpdateTsConfig(&tsConfig)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, "OK", tsResult)

	tsConfigErrOne := entity.Tsconfig{Module: "errOneConfig", Target: "testTarget", SourceMap: true, Excluding: 1}
	expectedTsErrorOne := errDB
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\" WHERE (module = $1)")).
		WithArgs("errOneConfig").
		WillReturnError(expectedTsErrorOne)
	_, tsReturnedErr := repo.UpdateTsConfig(&tsConfigErrOne)
	if assert.Error(t, tsReturnedErr) {
		assert.Equal(t, expectedTsErrorOne, tsReturnedErr)
	}

	expectedTsErrorTwo := errDB
	tsConfigErrTwo := entity.Tsconfig{Module: "errTwoConfig", Target: "testTarget", SourceMap: true, Excluding: 1}
	tsRows = getTsConfigRows(tsConfigErrTwo.Module)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\" WHERE (module = $1)")).
		WithArgs("errTwoConfig").
		WillReturnRows(tsRows)
	m.ExpectExec(formatRequest("UPDATE tsconfigs SET target = $1, source_map = $2, excluding = $3 WHERE module = $4")).
		WithArgs(tsConfigErrTwo.Target, strconv.FormatBool(tsConfigErrTwo.SourceMap),
			strconv.Itoa(tsConfigErrTwo.Excluding), tsConfigErrTwo.Module).
		WillReturnError(expectedTsErrorTwo)
	_, tsReturnedErrTwo := repo.UpdateTsConfig(&tsConfigErrTwo)
	if assert.Error(t, tsReturnedErrTwo) {
		assert.Equal(t, expectedTsErrorTwo, tsReturnedErrTwo)
	}

	expectedTsErrorThree := errors.New("fields are empty")
	tsConfigErrThree := entity.Tsconfig{Module: "errThreeConfig", Target: "", SourceMap: true, Excluding: 1}
	tsRows = getTsConfigRows(tsConfigErrThree.Module)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tsconfigs\" WHERE (module = $1)")).
		WithArgs("errThreeConfig").
		WillReturnRows(tsRows)
	m.ExpectExec(formatRequest("UPDATE tsconfigs SET target = $1, source_map = $2, excluding = $3 WHERE module = $4")).
		WithArgs(tsConfigErrThree.Target, strconv.FormatBool(tsConfigErrThree.SourceMap),
			strconv.Itoa(tsConfigErrThree.Excluding), tsConfigErrThree.Module).
		WillReturnError(expectedTsErrorThree)
	_, tsReturnedErrThree := repo.UpdateTsConfig(&tsConfigErrThree)
	if assert.Error(t, tsReturnedErrThree) {
		assert.Equal(t, expectedTsErrorThree, tsReturnedErrThree)
	}
}

func TestPostgresStorage_UpdateTempConfig(t *testing.T) {
	m, db, _ := newDB()

	repo := PostgresStorage{DB: db}
	tempConfig := entity.Tempconfig{RestApiRoot: "testApiRoot", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	tempRows := getTempConfigRows(tempConfig.RestApiRoot)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\" WHERE (rest_api_root = $1)")).
		WithArgs("testApiRoot").
		WillReturnRows(tempRows)
	m.ExpectExec(formatRequest("UPDATE tempconfigs SET remoting = $1, port = $2, host = $3, legasy_explorer = $4 WHERE rest_api_root = $5")).
		WithArgs(tempConfig.Remoting, tempConfig.Port, tempConfig.Host,
			strconv.FormatBool(tempConfig.LegasyExplorer), tempConfig.RestApiRoot).
		WillReturnResult(sqlmock.NewResult(0, 1))
	tempResult, err := repo.UpdateTempConfig(&tempConfig)
	if err != nil {
		t.Error("error during unit testing: ", err)
	}
	assert.Equal(t, "OK", tempResult)

	tempConfigErrOne := entity.Tempconfig{RestApiRoot: "errOneConfig", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	expectedTempErrorOne := errors.New("record not found")
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\" WHERE (rest_api_root = $1)")).
		WithArgs("errOneConfig").
		WillReturnError(expectedTempErrorOne)
	_, tempReturnedErr := repo.UpdateTempConfig(&tempConfigErrOne)
	if assert.Error(t, tempReturnedErr) {
		assert.Equal(t, expectedTempErrorOne, tempReturnedErr)
	}

	expectedTempErrorTwo := errors.New("db error")
	tempConfigErrTwo := entity.Tempconfig{RestApiRoot: "errTwoConfig", Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	tempRows = getTempConfigRows(tempConfigErrTwo.RestApiRoot)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\" WHERE (rest_api_root = $1)")).
		WithArgs("errTwoConfig").
		WillReturnRows(tempRows)
	m.ExpectExec(formatRequest("UPDATE tempconfigs SET remoting = $1, port = $2, host = $3, legasy_explorer = $4 WHERE rest_api_root = $5")).
		WithArgs(tempConfigErrTwo.Remoting, tempConfigErrTwo.Port, tempConfigErrTwo.Host,
			strconv.FormatBool(tempConfigErrTwo.LegasyExplorer), tempConfigErrTwo.RestApiRoot).
		WillReturnError(expectedTempErrorTwo)
	_, tempReturnedErrTwo := repo.UpdateTempConfig(&tempConfigErrTwo)
	if assert.Error(t, tempReturnedErrTwo) {
		assert.Equal(t, expectedTempErrorTwo, tempReturnedErrTwo)
	}

	expectedTempErrorThree := errors.New("fields are empty")
	tempConfigErrThree := entity.Tempconfig{RestApiRoot: "errThreeConfig", Host: "", Port: "", Remoting: "", LegasyExplorer: true}
	tempRows = getTempConfigRows(tempConfigErrThree.RestApiRoot)
	m.ExpectQuery(formatRequest("SELECT * FROM \"tempconfigs\" WHERE (rest_api_root = $1)")).
		WithArgs("errThreeConfig").
		WillReturnRows(tempRows)
	m.ExpectExec(formatRequest("UPDATE tempconfigs SET remoting = $1, port = $2, host = $3, legasy_explorer = $4 WHERE rest_api_root = $5")).
		WithArgs(tempConfigErrThree.Remoting, tempConfigErrThree.Port, tempConfigErrThree.Host,
			strconv.FormatBool(tempConfigErrThree.LegasyExplorer), tempConfigErrThree.RestApiRoot).
		WillReturnError(expectedTempErrorThree)
	_, tempReturnedErrThree := repo.UpdateTempConfig(&tempConfigErrThree)
	if assert.Error(t, tempReturnedErrThree) {
		assert.Equal(t, expectedTempErrorThree, tempReturnedErrThree)
	}
}

func getMongoDBRows(configID string) *sqlmock.Rows {
	var fieldNames = []string{"domain", "mongodb", "host", "port"}
	rows := sqlmock.NewRows(fieldNames)
	mongodbConfig := entity.Mongodb{Domain: configID, Mongodb: true, Host: "testHost", Port: "testPort"}
	rows = rows.AddRow(mongodbConfig.Domain, mongodbConfig.Mongodb, mongodbConfig.Host, mongodbConfig.Port)
	return rows
}

func getTsConfigRows(configID string) *sqlmock.Rows {
	var fieldNames = []string{"module", "target", "source_map", "excluding"}
	rows := sqlmock.NewRows(fieldNames)
	tsConfig := entity.Tsconfig{Module: configID, Target: "testTarget", SourceMap: true, Excluding: 1}
	rows = rows.AddRow(tsConfig.Module, tsConfig.Target, tsConfig.SourceMap, tsConfig.Excluding)
	return rows
}

func getTempConfigRows(configID string) *sqlmock.Rows {
	var fieldNames = []string{"rest_api_root", "host", "port", "remoting", "legasy_explorer"}
	rows := sqlmock.NewRows(fieldNames)
	tempConfig := entity.Tempconfig{RestApiRoot: configID, Host: "testHost", Port: "testPort", Remoting: "testRemoting", LegasyExplorer: true}
	rows = rows.AddRow(tempConfig.RestApiRoot, tempConfig.Host, tempConfig.Port, tempConfig.Remoting, tempConfig.LegasyExplorer)
	return rows
}
