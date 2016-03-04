package y

import (
	"database/sql"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	// models work with mysql only
	_ "github.com/go-sql-driver/mysql"
)

const (
	FixturedID   int64 = 1
	FixturedName       = "foo"
)

var (
	dsn string
)

func init() {
	env := func(key, defaultValue string) string {
		if value := os.Getenv(key); value != "" {
			return value
		}
		return defaultValue
	}
	user := env("MYSQL_TEST_USER", "root")
	pass := env("MYSQL_TEST_PASS", "")
	prot := env("MYSQL_TEST_PROT", "tcp")
	addr := env("MYSQL_TEST_ADDR", "localhost:3306")
	dbname := env("MYSQL_TEST_DBNAME", "y_test")
	netAddr := fmt.Sprintf("%s(%s)", prot, addr)
	dsn = fmt.Sprintf("%s:%s@%s/%s?timeout=30s&strict=true", user, pass, netAddr, dbname)
}

var _ = Describe("Proxy", func() {
	var db *sql.DB

	type Fixtured struct {
		ID   int64 `y:",pk,autoincr"`
		Name string
	}

	BeforeEach(func() {
		db, _ = sql.Open("mysql", dsn)
		_, err := db.Exec(`
      CREATE TABLE y_test (
        id INTEGER PRIMARY KEY AUTO_INCREMENT,
        name TEXT,
        _version INTEGER
      )`)
		if err != nil {
			panic(err.Error())
		}
	})

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

		It("no error occuried", func() {
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

		It("no error occuried", func() {
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

		It("no error occuried", func() {
			Expect(err).To(BeNil())
		})
		It("the name is changed", func() {
			Expect(ytest.Name).To(Equal("bar"))
		})
		It("the versions is updated", func() {
			Expect(ytest.Version.Int64).To(Equal(int64(2)))
		})
	})
})
