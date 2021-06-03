package foundationdb

import (
	"log"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/directory"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
)

var UrlMatch subspace.Subspace
var HashCodeMatch subspace.Subspace

type FdbWrapper struct {
	UrlMatch      subspace.Subspace
	HashCodeMatch subspace.Subspace
	db            fdb.Database
}

// func (fdbObject *FdbWrapper) SetDb(object fdb.Database) {

// }

func (fdbObject *FdbWrapper) CreateDb(firstSpace string, secondSpace string) {
	fdb.MustAPIVersion(620)
	fdbObject.db = fdb.MustOpenDefault()
	fdbObject.db.Options().SetTransactionTimeout(60000)
	fdbObject.db.Options().SetTransactionRetryLimit(100)

	schedulingDir, err := directory.CreateOrOpen(fdbObject.db, []string{"urlShortner"}, nil)
	if err != nil {
		log.Fatal(err)
	}
	fdbObject.UrlMatch = schedulingDir.Sub(firstSpace)
	fdbObject.HashCodeMatch = schedulingDir.Sub(secondSpace)
}

func (fdbObject *FdbWrapper) InsertIntoUrlMatch(key string, value string) {
	fdbObject.db.Transact(func(tr fdb.Transaction) (ret interface{}, err error) {
		tr.Set(fdbObject.UrlMatch.Pack((tuple.Tuple{key})), []byte(value))
		return
	})
}

func (fdbObject *FdbWrapper) InsertIntoHashCodeMatch(key string, value string) {
	fdbObject.db.Transact(func(tr fdb.Transaction) (ret interface{}, err error) {
		tr.Set(fdbObject.HashCodeMatch.Pack((tuple.Tuple{key})), []byte(value))
		return
	})
}

func (fdbObject *FdbWrapper) LookIntoHasCodeMatch(key string) (bool, error) {
	result := true
	ret, err := fdbObject.db.Transact(func(tr fdb.Transaction) (ret interface{}, err error) {
		ret = tr.Get(fdbObject.UrlMatch.Pack((tuple.Tuple{key}))).MustGet()
		tr.GetReadVersion()
		return
	})
	if err != nil {
		return false, err
	}
	if string(ret.([]byte)) == "" {
		result = false
	}
	return result, nil
}

func (fdbObject *FdbWrapper) GetFromUrlMatch(key string) (string, error, bool) {
	resultantKey := ""
	exists := false
	ret, err := fdbObject.db.Transact(func(tr fdb.Transaction) (ret interface{}, err error) {
		ret = tr.Get(fdbObject.UrlMatch.Pack((tuple.Tuple{key}))).MustGet()
		return
	})
	if err != nil {
		log.Println("Here is the Error---------------------------------------->", err)
		return "", err, false
	}
	if ret != nil {
		resultantKey = string(ret.([]byte))
		exists = true
	}
	return resultantKey, nil, exists
}

// ret, err := db.Transact(func (tr fdb.Transaction) (ret interface{}, e error) {
// 	ret = tr.Get(fdb.Key("hello")).MustGet()
// 	return
// })
// if err != nil {
// 	log.Fatalf("Unable to read FDB database value (%v)", err)
// }

// v := ret.([]byte)
// fmt.Printf("hello, %s\n", string(v))
