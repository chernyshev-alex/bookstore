package users

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/config"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/datasources/mysql/users_db"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/domain/users"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/stretchr/testify/assert"
)

const (
	ENV_CONFIG_VAR = "CONFIG"
)

var (
	db *sql.DB
	dq IUserDao
)

func TestMain(m *testing.M) {
	defer cleanUpDB()
	setupTest()
	cleanUpDB()
	os.Exit(m.Run())
}

func setupTest() {
	var conf config.Config

	configFile, _ := filepath.Abs(os.Getenv(ENV_CONFIG_VAR))
	if err := cleanenv.ReadConfig(configFile, &conf); err != nil {
		panic(err)
	}

	mysqlConf := users_db.MakeMySQLConfig(conf)
	db = users_db.ProvideSqlClient(mysqlConf)
	dq = NewUserDao(db)
}

func cleanUpDB() {
	stmt, err := db.Prepare("truncate table users_db.users;")
	if err != nil {
		fmt.Println(err)
	}
	stmt.Exec()
}

func TestDaoSQLC(t *testing.T) {
	u := users.User{
		FirstName: "fname",
		LastName:  "lname",
		Email:     "email@domain.com",
		Status:    "status",
		Password:  "hashed_psw",
	}
	err := dq.Save(&u)
	if err != nil {
		t.Fatal(err)
	}
	result, err := dq.Get(int32(u.Id))
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, u.Id, result.Id)

	err = dq.Save(&u)
	if err == nil {
		t.Fatal(err)
	}
	assert.NotNil(t, err, "same Email isn't allowed")

	err = dq.FindByEmailAndPsw(&u)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, "fname", u.FirstName)

	users, err := dq.FindByStatus("status")
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(users) == 1)

	u.FirstName = "changed"
	err = dq.Update(&u)
	if err != nil {
		t.Fatal(err)
	}

	result, err = dq.Get(int32(u.Id))
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, "changed", result.FirstName)

	err = dq.Delete(&u)
	if err != nil {
		t.Fatal(err)
	}

	_, err = dq.Get(int32(u.Id))
	if err == nil {
		t.Fatal(err)
	}
}
