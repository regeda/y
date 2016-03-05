package y

import (
	"database/sql"
	"fmt"
	"net"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const (
	FixturedID   int64 = 1
	FixturedName       = "foo"
)

type env func(string, string) string

type testrdb interface {
	Setup(env) (string, error)
	CreateTable() string
}

type mysqlrdb struct{}

func (m mysqlrdb) Setup(e env) (dsn string, err error) {
	user := e("MYSQL_TEST_USER", "root")
	pass := e("MYSQL_TEST_PASS", "")
	prot := e("MYSQL_TEST_PROT", "tcp")
	addr := e("MYSQL_TEST_ADDR", "localhost:3306")
	dbname := e("MYSQL_TEST_DBNAME", "y_test")
	netAddr := fmt.Sprintf("%s(%s)", prot, addr)
	dsn = fmt.Sprintf("%s:%s@%s/%s?timeout=30s&strict=true", user, pass, netAddr, dbname)
	c, err := net.Dial(prot, addr)
	if err == nil {
		c.Close()
	}
	return
}

func (m mysqlrdb) CreateTable() string {
	return `
CREATE TABLE y_test (
	id INTEGER PRIMARY KEY AUTO_INCREMENT,
	name TEXT,
	_version INTEGER
)`
}

type postgresrdb struct{}

func (p postgresrdb) Setup(e env) (dsn string, err error) {
	SetBuilderProvider(Postgres)
	user := e("PG_TEST_USER", "postgres")
	pass := e("PG_TEST_PASS", "")
	addr := e("PG_TEST_ADDR", "localhost:5432")
	dbname := e("PG_TEST_DBNAME", "y_test")
	dsn = fmt.Sprintf("postgres://%s:%s@%s/%s?connect_timeout=30", user, pass, addr, dbname)
	return
}

func (p postgresrdb) CreateTable() string {
	return `
CREATE TABLE y_test (
	id SERIAL PRIMARY KEY,
	name TEXT,
	_version INTEGER
)`
}

var (
	rdb     testrdb
	rdbtype string
	rdberr  error
	dsn     string
)

func rdbfactory(rdbtype string) (testrdb, error) {
	switch rdbtype {
	case "mysql":
		return mysqlrdb{}, nil
	case "postgres":
		return postgresrdb{}, nil
	}
	return nil, fmt.Errorf("Unknown RDBS: %s", rdbtype)
}

func init() {
	env := func(key, defaultValue string) string {
		if value := os.Getenv(key); value != "" {
			return value
		}
		return defaultValue
	}
	rdbtype = env("DB", "mysql")
	rdb, rdberr = rdbfactory(rdbtype)
	if rdberr == nil {
		dsn, rdberr = rdb.Setup(env)
	}
}

var _ = Describe("Proxy", func() {
	var db *sql.DB

	type Fixtured struct {
		ID   int64 `y:",pk,autoincr"`
		Name string
	}

	// common setup
	BeforeEach(func() {
		if rdberr != nil {
			Skip(fmt.Sprintf("%s server has got an error: %s", rdbtype, rdberr.Error()))
		}
		db, _ = sql.Open(rdbtype, dsn)
		_, err := db.Exec(rdb.CreateTable())
		if err != nil {
			Skip(err.Error())
		}
	})

	// common teardown
	AfterEach(func() {
		db.Exec("DROP TABLE y_test")
		db.Close()
	})

	Context("when one item created", func() {
		type YTest struct {
			Fixtured
		}

		var (
			err   error
			ytest YTest
		)

		BeforeEach(func() {
			ytest = YTest{Fixtured{Name: FixturedName}}
			_, err = New(&ytest).Put(db)
		})

		It("no error occurred", func() {
			Expect(err).To(BeNil())
		})
		It("the primary key is not empty", func() {
			Expect(ytest.ID).To(Equal(FixturedID))
		})
	})

	Context("when one item loaded", func() {
		type YTest struct {
			Fixtured
		}

		var (
			err   error
			ytest YTest
		)

		BeforeEach(func() {
			// insert a fixture
			New(&YTest{Fixtured{FixturedID, FixturedName}}).Put(db)
			// load the testable bean
			ytest = YTest{Fixtured{ID: FixturedID}}
			err = New(&ytest).Load(db)
		})

		It("no error occurred", func() {
			Expect(err).To(BeNil())
		})
		It("the name is not empty", func() {
			Expect(ytest.Name).To(Equal(FixturedName))
		})
	})

	Context("when an item updated", func() {
		type YTest struct {
			Fixtured
			Versionable
		}

		var (
			err   error
			ytest YTest
		)

		BeforeEach(func() {
			ytest = YTest{Fixtured{FixturedID, FixturedName}, MakeVersionable(1)}
			// insert a fixture
			New(&ytest).Put(db)
			// update the item
			err = New(&ytest).Update(db, Values{"name": "bar"})
		})

		It("no error occurred", func() {
			Expect(err).To(BeNil())
		})
		It("the name is changed", func() {
			Expect(ytest.Name).To(Equal("bar"))
		})
		It("the version is updated", func() {
			Expect(ytest.Version.Int64).To(Equal(int64(2)))
		})
	})
})
