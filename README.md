# y

[![Build Status](https://travis-ci.org/Repo2/y.svg?branch=master)](https://travis-ci.org/Repo2/y)

Be faster with Y. The simplest ORM-like framework for Golang.

## Install

```bash
go get github.com/lann/squirrel
go get github.com/Repo2/y
```

## Actions

### Fetch
**Fetch** executes ```SELECT``` statement and returns a collection of objects.
```go
type Account struct {
  ID   int64  `db:"id,pk"`
  Name string `db:"name"`
}
c, err := y.New(Account{}).Fetch(db)
if err != nil {
  log.Panicln(err)
}
log.Printf("%#v\n", c.List())
```

### Find
**Find** modifies a query for custom selection.
```go
type Order struct {
  ID    int64 `db:"id,pk"`
  Price int   `db:"price"`
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
  ID    int64  `db:"id,pk"`
  Title string `db:"title"`
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
  ID   int64  `db:"id,pk,autoincr"`
  Name string `db:"name"`
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
type Versionable struct {
  Version int `db:"_version"`
}
type Car struct {
  ID    int64 `db:"id,pk"`
  Power int   `db:"power"`
  Versionable
}
car := Car{ID: 1}
err := y.New(&car).Update(db, y.Values{"power": 50})
if err != nil {
  log.Panicln(err)
}
log.Printf("%#v\n", car)
```  

### Delete
**Delete** executes ```DELETE``` statement. The action removes an object by primary keys.
```go
type Account struct {
	ID    int64 `db:",pk"`
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
  ID     int64 `db:",pk"`
  UserID int64 `db:",fk"`
  Name   string
}
type User struct {
  ID          int64     `db:"id,pk"`
  DeviceArray []*Device `db:"-"`
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

## TODO
- [ ] More tests!
