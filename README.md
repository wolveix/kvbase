# kvbase
[![build](https://img.shields.io/github/workflow/status/Wolveix/kvbase/Build?label=build)](https://github.com/Wolveix/kvbase/workflows/Go) [![report](https://goreportcard.com/badge/github.com/Wolveix/kvbase)](https://goreportcard.com/report/github.com/Wolveix/kvbase) [![coverage](https://img.shields.io/codecov/c/github/Wolveix/kvbase)](https://codecov.io/gh/Wolveix/kvbase) [![documentation](https://godoc.org/github.com/Wolveix/kvbase?status.svg)](https://pkg.go.dev/github.com/Wolveix/kvbase) [![license](https://img.shields.io/github/license/Wolveix/kvbase)](https://github.com/Wolveix/kvbase/blob/master/LICENSE) [![version](https://img.shields.io/github/v/tag/Wolveix/kvbase?label=version)](https://github.com/Wolveix/kvbase/releases/latest)

A simple abstraction library for key value stores.

Currently supported stores:
- [BadgerDB](https://github.com/dgraph-io/badger)
- [Bitcask](https://github.com/prologic/bitcask)
- [BoltDB](https://github.com/boltdb/bolt)
- [BboltDB](https://github.com/etcd-io/bbolt)
- [Diskv](https://github.com/peterbourgon/diskv)
- [Go-Cache](https://github.com/patrickmn/go-cache)
- [LevelDB](https://github.com/syndtr/goleveldb)

## Getting Started

### Installing

To start using kvbase, install Go and run `go get`:

```sh
$ go get github.com/Wolveix/kvbase/..
```

This will retrieve the library. You can interact with each supported store in the exact same way, simply swap out the specified backend when calling the `New` command.

<hr>

### Opening a database

All stores utilize the same `Backend` interface. The following functions are available for every backend:

- `Count(bucket string) (int, error)`
- `Create(bucket string, key string, model interface{}) error`
- `Delete(bucket string, key string) error`
- `Drop(bucket string) error`
- `Get(bucket string, model interface{}) (*map[string]interface{}, error)`
- `Initialize(backend string, source string, memory bool) error`
- `Read(bucket string, key string, model interface{}) error`
- `Update(bucket string, key string, model interface{}) error`

Stores can be opened similarly to how `database/sql` handles databases. Import `Wolveix/kvbase` as well as the backend you want to use `Wolveix/kvbase/backend/badgerdb`, then call `kvbase.New("badgerdb", "data", false)`:

```go
package main

import (
	"log"

	"github.com/Wolveix/kvbase"
	_ "github.com/Wolveix/kvbase/backend/badgerdb"
)

func main() {
    kv, err := kvbase.New("badgerdb", "data", false)
    if err != nil {
        log.Fatal(err)
    }
}
```

These functions expect a source to be specified. Some drivers utilize a file, others utilize a folder. Not all backends require the boolean value after the source (this value enables in-memory mode, disabling persistent database storage).

<hr>

### Counting entries within a bucket

The `Count()` function expects a bucket (as a `string`),:

```go
counter, err := kv.Count("users")
if err != nil {
    log.Fatal(err)
}

fmt.Print(counter) //This will output 1.
```

### Creating an entry

The `Create()` function expects a bucket (as a `string`), a key (as a `string`) and a struct containing your data (as an `interface{}`):

```go
type User struct {
	Password string
	Username string
}

user := User{
    "Password123",
    "JohnSmith123",
}

if err := kv.Create("users", user.Username, &user); err != nil {
    log.Fatal(err)
}
```
If the key already exists, this will **fail**.

### Deleting an entry

The `Delete()` function expects a bucket (as a `string`), a key (as a `string`):

```go
if err := kv.Delete("users", "JohnSmith01"); err != nil {
    log.Fatal(err)
}
```

### Dropping a bucket

The `Drop()` function expects a bucket (as a `string`):

```go
if err := kv.Drop("users"); err != nil {
    log.Fatal(err)
}
```

### Getting all entries

The `Get()` function expects a bucket (as a `string`), and a struct to unmarshal your data into (as an `interface{}`):

```go
results, err := kv.Get("users", User{})
if err != nil {
    log.Fatal(err)
}

s, _ := json.MarshalIndent(results, "", "\t")
fmt.Print(string(s))
```

`results` will now contain a `*map[string]interface{}` object. Note that the object doesn't support indexing, so `results["JohnSmith01"]` won't work; however, you can loop through the map to find specific keys.

### Reading an entry

The `Read()` function expects a bucket (as a `string`), a key (as a `string`) and a struct to unmarshal your data into (as an `interface{}`):

```go
user := User{}

if err := kv.Read("users", "JohnSmith01", &user); err != nil {
    log.Fatal(err)
}

fmt.Print(user.Password) //This will output Password123
```

### Updating an entry

The `Update()` function expects a bucket (as a `string`), a key (as a `string`) and a struct containing your data (as an `interface{}`):

```go
user := User{
    "Password456",
    "JohnSmith123",
}

if err := kv.Update("users", user.Username, &user); err != nil {
    log.Fatal(err)
}
```
If the key doesn't already exist, this will **fail**.

<hr>

## Credits
- Creator: [Robert Thomas](https://github.com/Wolveix)
- License: [GNU General Public License v3.0](https://github.com/Wolveix/kvbase/blob/master/LICENSE)
