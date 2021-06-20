package tests

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/conf"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/dao/mysql"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/dao/mysql/gen"
	"github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/dao/user_dao"
	"github.com/stretchr/testify/assert"
)

const (
	ENV_CONFIG_VAR = "CONFIG"
)

var (
	db *sql.DB
	dq user_dao.UserDao
)

func TestMain(m *testing.M) {
	defer cleanUpDB()
	setupTest()
	cleanUpDB()
	os.Exit(m.Run())
}

func setupTest() {
	conf, err := conf.LoadConfigFromEnv(ENV_CONFIG_VAR)
	if err != nil {
		panic(err.Error() + "check ENV_CONFIG_VAR")
	}
	mysqlConf := mysql.MakeConfig(conf)
	db = mysql.NewSqlClient(mysqlConf)
	dq = mysql.NewUserDao(db)
}

func cleanUpDB() {
	if db == nil {
		return
	}
	stmt, err := db.Prepare("truncate table users_db.users;")
	if err != nil {
		fmt.Println(err)
	}
	stmt.Exec()
}

func TestDaoSQLC(t *testing.T) {
	u := gen.InsertUserParams{
		FirstName:   nillableStr("fname"),
		LastName:    nillableStr("lname"),
		Email:       "email@domain.com",
		DateCreated: time.Now(),
		Status:      nillableStr("status"),
		Password:    nillableStr("hashed_psw"),
	}
	userId, err := dq.Save(u)
	if err != nil {
		t.Fatal(err)
	}
	result, err := dq.Get(userId)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, userId, result.ID)

	_, err = dq.Save(u)
	if err == nil {
		t.Fatal(err)
	}
	assert.NotNil(t, err, "same Email isn't allowed")

	findByEMailAndPswParams := gen.FindByEMailAndPswParams{
		Email:    u.Email,
		Password: u.Password,
		Status:   u.Status,
	}
	findResult, err := dq.FindByEmailAndPsw(findByEMailAndPswParams)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, u.FirstName, findResult.FirstName)

	users, err := dq.FindByStatus(u.Status.String)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(users) == 1)

	u.FirstName = nillableStr("changed")
	updateParams := gen.UpdateUserParams{
		FirstName: u.FirstName,
		Email:     u.Email,
		ID:        int32(userId),
	}
	err = dq.Update(updateParams)
	if err != nil {
		t.Fatal(err)
	}

	result, err = dq.Get(userId)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, "changed", result.FirstName.String)

	err = dq.Delete(userId)
	if err != nil {
		t.Fatal(err)
	}

	_, err = dq.Get(userId)
	if err == nil {
		t.Fatal(err)
	}
}

func nillableStr(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}
