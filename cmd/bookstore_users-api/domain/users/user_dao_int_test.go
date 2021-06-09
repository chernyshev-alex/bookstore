package users

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/config"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/datasources/mysql/users_db"
	"github.com/chernyshev-alex/bookstore/pkg/bookstore_utils_go/date_utils"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/stretchr/testify/assert"
)

const (
	ENV_CONFIG_VAR = "CONFIG"
)

var (
	dao *UserDAO
)

func setupTest() {

	configFile, _ := filepath.Abs(os.Getenv(ENV_CONFIG_VAR))

	var conf config.Config
	if err := cleanenv.ReadConfig(configFile, &conf); err != nil {
		panic(err)
	}

	mysqlConf := users_db.MakeMySQLConfig(conf)
	dao = ProvideUserDao(users_db.ProvideSqlClient(mysqlConf))
}

func cleanUpDB() {
	stmt, _ := dao.SqlClient.Prepare("truncate table users_db.users;")
	if _, err := stmt.Exec(); err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	setupTest()
	defer cleanUpDB()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestGetUser(t *testing.T) {
	defer cleanUpDB()
	_, err := dao.Get(0)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.Status())
}

func TestSave(t *testing.T) {
	defer cleanUpDB()
	u := User{Email: "user@email", DateCreated: date_utils.GetNowDbFormat()}
	err := dao.Save(&u)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, u.Id, int64(0))
	err = dao.Save(&u)
	assert.Equal(t, http.StatusBadRequest, err.Status())
}

func TestUpdate(t *testing.T) {
	defer cleanUpDB()
	u := User{FirstName: "change_me", Email: "user@email", DateCreated: date_utils.GetNowDbFormat()}
	err := dao.Save(&u)
	assert.Nil(t, err)
	u.FirstName = "updated"
	err = dao.Update(&u)
	assert.Nil(t, err)
	updated, update_error := dao.Get(u.Id)
	assert.Nil(t, update_error)
	assert.Equal(t, updated.FirstName, u.FirstName)
}

func TestDelete(t *testing.T) {
	defer cleanUpDB()
	u := User{FirstName: "delete_me", Email: "user@email", DateCreated: date_utils.GetNowDbFormat()}
	err := dao.Save(&u)
	assert.Nil(t, err)
	u.FirstName = "updated"
	err = dao.Delete(&u)
	assert.Nil(t, err)
	_, err = dao.Get(u.Id)
	assert.Equal(t, http.StatusNotFound, err.Status())
}

func TestFindByStatus(t *testing.T) {
	defer cleanUpDB()
	u := User{Status: "active", Email: "user@email", DateCreated: date_utils.GetNowDbFormat()}
	err := dao.Save(&u)
	assert.Nil(t, err)
	ls, findErr := dao.FindByStatus(u.Status)
	assert.Nil(t, findErr)
	assert.GreaterOrEqual(t, len(ls), 0)
	assert.Equal(t, u, ls[0])
}

func TestFindByEmailAndPsw(t *testing.T) {
	defer cleanUpDB()
	u := User{Status: "active", Email: "user@email", Password: "md5psw", DateCreated: date_utils.GetNowDbFormat()}
	err := dao.Save(&u)
	assert.Nil(t, err)
	err = dao.FindByEmailAndPsw(&u)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, u.Id, int64(0))
}
