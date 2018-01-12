package redis_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bukalapak/ottoman/redis"
	envx "github.com/bukalapak/ottoman/x/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/subosito/gotenv"
	redisc "gopkg.in/redis.v3"
)

type Provider interface {
	redis.Provider
	Del(keys ...string) *redisc.IntCmd
	Set(key string, value interface{}, expiration time.Duration) *redisc.StatusCmd
	TTL(key string) *redisc.DurationCmd
}

type CommonSuite struct {
	suite.Suite
	client Provider
	c      *redis.Redis
}

func (suite *CommonSuite) SetupSuite() {
	gotenv.Load("../env.sample")
}

func (suite *CommonSuite) TearDownSuite() {
	os.Clearenv()
}

func (suite *CommonSuite) TearDownTest() {
	suite.client.Del("foo")
}

func (suite *CommonSuite) TestRead() {
	suite.loadFixtures()

	b, err := suite.c.Read("foo")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), []byte(`{"foo":"bar"}`), b)
}

func (suite *CommonSuite) TestReadMap() {
	suite.loadFixtures()

	m, err := suite.c.ReadMap("foo")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), map[string]interface{}{"foo": "bar"}, m)
}

func (suite *CommonSuite) TestReadMap_invalidCache() {
	suite.loadFixtures()

	m, err := suite.c.ReadMap("baz")
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), m)
}

func (suite *CommonSuite) TestReadMap_unknownCache() {
	suite.loadFixtures()

	m, err := suite.c.ReadMap("boo")
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), m)
}

func (suite *CommonSuite) TestReadMulti() {
	suite.loadFixtures()

	keys := []string{
		"{x}.foo",
		"{x}.boo",
		"{x}.fox",
	}

	m, err := suite.c.ReadMulti(keys)
	assert.Nil(suite.T(), err)
	assert.Len(suite.T(), m, 2)

	for _, key := range keys {
		switch key {
		case "{x}.boo":
			assert.Nil(suite.T(), m[key])
		case "{x}.foo":
			assert.Equal(suite.T(), []byte(`{"foo":"bar"}`), m[key])
		case "{x}.fox":
			assert.Equal(suite.T(), []byte(`{"fox":"baz"}`), m[key])
		}
	}
}

func (suite *CommonSuite) TestIncr() {
	n, err := suite.c.Incr("foo")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), int64(1), n)
}

func (suite *CommonSuite) TestExpire() {
	err := suite.client.Set("foo", "bar", time.Minute).Err()
	assert.Nil(suite.T(), err)

	b, err := suite.c.Expire("foo", time.Hour)
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), b)

	cmd := suite.client.TTL("foo")
	assert.Nil(suite.T(), cmd.Err())
	assert.Equal(suite.T(), time.Hour, cmd.Val())
}

func (suite *CommonSuite) loadFixtures() {
	data := map[string]string{
		"foo":     `{"foo":"bar"}`,
		"fox":     `{"fox":"baz"}`,
		"{x}.foo": `{"foo":"bar"}`,
		"{x}.fox": `{"fox":"baz"}`,
		"baz":     `x`,
	}

	for k, v := range data {
		err := suite.client.Set(k, []byte(v), time.Minute).Err()
		if err != nil {
			panic(err)
		}
	}
}

type RedisSuite struct {
	*CommonSuite
}

func (suite *RedisSuite) SetupTest() {
	suite.client = redisc.NewClient(&redisc.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		DB:   int64(envx.Int("REDIS_DB")),
	})

	suite.c = redis.New(&redis.Option{
		Addrs: []string{os.Getenv("REDIS_ADDR")},
		DB:    int64(envx.Int("REDIS_DB")),
	})
}

func (suite *RedisSuite) TestName() {
	assert.Equal(suite.T(), "Redis", suite.c.Name())
}

func TestRedisSuite(t *testing.T) {
	suite.Run(t, &RedisSuite{new(CommonSuite)})
}

type RedisClusterSuite struct {
	*CommonSuite
}

func (suite *RedisClusterSuite) SetupTest() {
	addrs := strings.Split(os.Getenv("REDIS_CLUSTER_ADDR"), ",")

	suite.client = redisc.NewClusterClient(&redisc.ClusterOptions{
		Addrs: addrs,
	})

	suite.c = redis.New(&redis.Option{
		Addrs: addrs,
	})
}

func (suite *RedisClusterSuite) TestName() {
	assert.Equal(suite.T(), "Redis Cluster", suite.c.Name())
}

func (suite *RedisClusterSuite) TestReadMulti_CROSSSLOT() {
	suite.loadFixtures()

	keys := []string{
		"foo",
		"{x}.fox",
	}

	m, err := suite.c.ReadMulti(keys)
	assert.Contains(suite.T(), err.Error(), "CROSSSLOT")
	assert.Nil(suite.T(), m)
}

func TestRedisClusterSuite(t *testing.T) {
	suite.Run(t, &RedisClusterSuite{new(CommonSuite)})
}
