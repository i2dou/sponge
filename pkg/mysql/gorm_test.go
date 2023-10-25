package mysql

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var dsn = "root:123456@(192.168.3.37:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

func TestInit(t *testing.T) {
	db, err := Init(dsn, WithEnableTrace())
	if err != nil {
		// ignore test error about not being able to connect to real mysql
		t.Logf(fmt.Sprintf("connect to mysql failed, err=%v, dsn=%s", err, dsn))
		return
	}

	t.Logf("%+v", db.Name())
}

func Test_gormConfig(t *testing.T) {
	o := defaultOptions()
	o.apply(
		WithLog(),
		WithLogging(nil),
		WithLogging(nil, 4),
		WithSlowThreshold(time.Millisecond*100),
		WithEnableTrace(),
		WithMaxIdleConns(5),
		WithMaxOpenConns(50),
		WithConnMaxLifetime(time.Minute*3),
		WithEnableForeignKey(),
		WithLogRequestIDKey("request_id"),
		WithRWSeparation([]string{
			"root:123456@(192.168.3.37:3306)/slave1",
			"root:123456@(192.168.3.37:3306)/slave2"},
			"root:123456@(192.168.3.37:3306)/master"),
		WithGormPlugin(nil),
	)

	c := gormConfig(o)
	assert.NotNil(t, c)

	err := rwSeparationPlugin(o)
	assert.NotNil(t, err)
}

func TestGetTableName(t *testing.T) {
	name := GetTableName(&userExample{})
	assert.NotEmpty(t, name)

	name = GetTableName(userExample{})
	assert.NotEmpty(t, name)

	name = GetTableName("table")
	assert.Empty(t, name)
}
