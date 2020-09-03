# y

[![Build Status](https://travis-ci.org/regeda/y.svg?branch=master)](https://travis-ci.org/regeda/y)
[![Go Report Card](https://goreportcard.com/badge/github.com/regeda/y)](https://goreportcard.com/report/github.com/regeda/y)

Be faster with Y. The simplest ORM-like framework for Golang.
The purpose of the library is a pure data modeling from a scratch.

## Install

```bash
go get github.com/regeda/y
```

## RDBS
You can use both databases: MySQL (default) and Postgres.
If you wanted to use Postgres you should call ```SetBuilderProvider``` method once before any Y routine:
```go
y.SetBuilderProvider(y.Postgres)
```
Y has been tested with following drivers:
- ```github.com/go-sql-driver/mysql```
- ```github.com/lib/pq```

## Actions

* [Fetch](#fetch)
* [FindBy](#findby)
* [Load](#load)
* [Put](#put)
* [Update](#update)
* [Delete](#delete)
* [Join](#join)

### Fetch
**Fetch** executes ```SELECT``` statement and returns a collection of objects.
```go
type Account struct {
  ID   int64  `y:"id,pk"`
  Name string `y:"name"`
}
c, err := y.New(Account{}).Fetch(db)
if err != nil {
  log.Panicln(err)
}
log.Printf("%#v\n", c.List())
```

### FindBy
**FindBy** modifies a query for custom selection.
```go
type Order struct {
  ID    int64 `y:"id,pk"`
  Price int   `y:"price"`
}
c, err := y.New(Order{}).
  FindBy(
  func(b squirrel.SelectBuilder) squirrel.SelectBuilder {
    return b.Where("price > ?", 10)
  }).
  Fetch(db)
if err != nil {
  log.Panicln(err)
}
log.Printf("%#v\n", c.List())
```

### Load
**Load** executes ```SELECT``` statement for one row and loads the object in self.
```go
type Todo struct {
  ID    int64  `y:"id,pk"`
  Title string `y:"title"`
}
todo := Todo{ID: 1}
err := y.New(&todo).Load(db)
if err != nil {
  log.Panicln(err)
}
log.Printf("%#v\n", todo)
```

### Put
**Put** executes ```INSERT``` statement and assigns auto-increment value.
```go
type User struct {
  ID   int64  `y:"id,pk,autoincr"`
  Name string `y:"name"`
}
user := User{Name: "Harry"}
_, err := y.New(&user).Put(db)
if err != nil {
  log.Panicln(err)
}
log.Printf("%#v\n", user)
```
You can use **Put** for batch statement also.
```go
type Log struct {
  Msg string
}
logs := []Log{
  {"It"}, {"Works"},
}
affected, err := y.New(logs).Put(db)
if err != nil {
  log.Panicln(err)
}
log.Printf("%#v\n", affected)
```

### Update
**Update** executes ```UPDATE``` statement. The action compares origin and modified objects by their version in the database.
```go
type Car struct {
  ID    int64 `y:"id,pk"`
  Power int   `y:"power"`
  y.Versionable
}
var err error
car := Car{ID: 1}
err = y.New(&car).MustLoad(db).Update(db, y.Values{"power": 50})
if err != nil {
  log.Panicln(err)
}
log.Printf("%#v\n", car)
```

### Delete
**Delete** executes ```DELETE``` statement. The action removes an object by primary keys.
```go
type Account struct {
	ID    int64 `y:",pk"`
	Email string
}
acc := Account{ID: 1}
affected, err := y.New(acc).Delete(db)
if err != nil {
  log.Panicln(err)
}
log.Printf("Affected rows: %d\n", affected)
```

### Join
**Join** builds relations by foreign keys
```go
type Device struct {
  ID     int64 `y:",pk"`
  UserID int64 `y:",fk"`
  Name   string
}
type User struct {
  ID          int64     `y:"id,pk"`
  DeviceArray []*Device `y:"-"`
}
users, err := y.New(User{}).Fetch(db)
if err != nil {
  log.Panicln(err)
}
if !users.Empty() {
  devices, _ := y.New(Device{}).Join(db, users)
  log.Printf("All users with their devices: %#v\n", users)
  log.Printf("All devices: %#v\n", devices)
}
```

## Pull requests
Please, be free to make pull requests or issues posting.

## License
MIT

Original source: https://github.com/regeda/y
