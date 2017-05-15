package memcache_test

import (
	"bytes"
	"compress/zlib"
	"io"
	"os"
	"testing"
	"time"

	gomemcache "github.com/bradfitz/gomemcache/memcache"
	"github.com/bukalapak/ottoman/memcache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MemcacheSuite struct {
	suite.Suite
	client *gomemcache.Client
	c      *memcache.Memcache
}

func (suite *MemcacheSuite) Addr() string {
	return os.Getenv("MEMCACHE_ADDR")
}

func (suite *MemcacheSuite) SetupTest() {
	suite.NewClient(suite.Addr(), memcache.Option{
		Compress: true,
	})
}

func (suite *MemcacheSuite) NewClient(addr string, option memcache.Option) {
	suite.client = gomemcache.New(addr)
	suite.c = memcache.New([]string{addr}, option)
}

func (suite *MemcacheSuite) TearDownTest() {
	suite.client.Delete("foo")
	suite.client.Delete("baz")
}

func (suite *MemcacheSuite) TestName() {
	assert.Equal(suite.T(), "Memcached", suite.c.Name())
}

func (suite *MemcacheSuite) TestRead() {
	suite.NewClient(suite.Addr(), memcache.Option{
		Compress: false,
	})

	suite.loadFixtures(false)

	b, err := suite.c.Read("foo")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), `{"foo":"bar"}`, string(b))
}

func (suite *MemcacheSuite) TestRead_miss() {
	suite.loadFixtures(false)

	b, err := suite.c.Read("boo")
	assert.Error(suite.T(), gomemcache.ErrCacheMiss, err.Error())
	assert.Nil(suite.T(), b)
}

func (suite *MemcacheSuite) TestRead_zlib() {
	suite.loadFixtures(true)

	b, err := suite.c.Read("foo")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), `{"foo":"bar"}`, string(b))
}

func (suite *MemcacheSuite) TestRead_zlib_uncompressedCache() {
	suite.loadFixtures(false)

	b, err := suite.c.Read("foo")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), `{"foo":"bar"}`, string(b))
}

func (suite *MemcacheSuite) TestReadMap() {
	suite.loadFixtures(true)

	m, err := suite.c.ReadMap("foo")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), map[string]interface{}{"foo": "bar"}, m)
}

func (suite *MemcacheSuite) TestReadMap_invalidCache() {
	suite.loadFixtures(true)

	m, err := suite.c.ReadMap("baz")
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), m)
}

func (suite *MemcacheSuite) TestReadMap_unknownCache() {
	suite.loadFixtures(true)

	m, err := suite.c.ReadMap("boo")
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), m)
}

func (suite *MemcacheSuite) TestReadMulti() {
	suite.loadFixtures(true)

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

func (suite *MemcacheSuite) TestReadMulti_failure() {
	suite.loadFixtures(true)
	suite.NewClient(suite.Addr(), memcache.Option{
		Timeout: 1 * time.Microsecond,
	})

	defer func() {
		suite.NewClient(suite.Addr(), memcache.Option{
			Compress: true,
		})
	}()

	m, err := suite.c.ReadMulti([]string{"foo", "boo", "fox"})
	assert.NotNil(suite.T(), err)
	assert.Len(suite.T(), m, 0)
}

func TestMemcacheSuite(t *testing.T) {
	suite.Run(t, new(MemcacheSuite))
}

func (suite *MemcacheSuite) loadFixtures(compress bool) {
	data := map[string]string{
		"foo": `{"foo":"bar"}`,
		"fox": `{"fox":"baz"}`,
		"baz": `x`,
	}

	for k, v := range data {
		var z []byte

		if compress {
			var b bytes.Buffer

			w := zlib.NewWriter(&b)
			io.WriteString(w, v)
			w.Close()

			z = b.Bytes()
		} else {
			z = []byte(v)
		}

		x := &gomemcache.Item{
			Key:        k,
			Value:      z,
			Expiration: int32(time.Minute.Seconds()),
		}

		err := suite.client.Set(x)
		if err != nil {
			panic(err)
		}
	}
}
