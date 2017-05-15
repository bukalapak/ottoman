package redis_test

import (
	"os"
	"testing"
	"time"

	"github.com/bukalapak/ottoman/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/subosito/envx"
	"github.com/subosito/gotenv"
	redisc "gopkg.in/redis.v3"
)

type RedisSuite struct {
	suite.Suite
	client *redisc.Client
	c      *redis.Redis
}

func (suite *RedisSuite) SetupSuite() {
	gotenv.Load("../env.sample")
}

func (suite *RedisSuite) TearDownSuite() {
	os.Clearenv()
}

func (suite *RedisSuite) SetupTest() {
	options := &redisc.Options{
		Addr: os.Getenv("REDIS_CACHE_ADDR"),
		DB:   int64(envx.Int("REDIS_CACHE_DB")),
	}

	suite.client = redisc.NewClient(options)
	suite.c = redis.New(options)
}

func (suite *RedisSuite) TearDownTest() {
	suite.client.Del("foo")
}

func (suite *RedisSuite) TestName() {
	assert.Equal(suite.T(), "Redis", suite.c.Name())
}

func (suite *RedisSuite) TestRead() {
	suite.loadFixtures()

	b, err := suite.c.Read("foo")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), []byte(`{"foo":"bar"}`), b)
}

func (suite *RedisSuite) TestReadMap() {
	suite.loadFixtures()

	m, err := suite.c.ReadMap("foo")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), map[string]interface{}{"foo": "bar"}, m)
}

func (suite *RedisSuite) TestReadMap_invalidCache() {
	suite.loadFixtures()

	m, err := suite.c.ReadMap("baz")
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), m)
}

func (suite *RedisSuite) TestReadMap_unknownCache() {
	suite.loadFixtures()

	m, err := suite.c.ReadMap("boo")
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), m)
}

func (suite *RedisSuite) TestReadMulti() {
	suite.loadFixtures()

	keys := []string{
		"foo",
		"boo",
		"fox",
	}

	m, err := suite.c.ReadMulti(keys)
	assert.Nil(suite.T(), err)
	assert.Len(suite.T(), m, 2)

	for _, key := range keys {
		switch key {
		case "boo":
			assert.Nil(suite.T(), m[key])
		case "foo":
			assert.Equal(suite.T(), []byte(`{"foo":"bar"}`), m[key])
		case "fox":
			assert.Equal(suite.T(), []byte(`{"fox":"baz"}`), m[key])
		}
	}
}

func TestRedisSuite(t *testing.T) {
	suite.Run(t, new(RedisSuite))
}

func (suite *RedisSuite) loadFixtures() {
	data := map[string]string{
		"foo": `{"foo":"bar"}`,
		"fox": `{"fox":"baz"}`,
		"baz": `x`,
	}

	for k, v := range data {
		err := suite.client.Set(k, []byte(v), time.Minute).Err()
		if err != nil {
			panic(err)
		}
	}
}
