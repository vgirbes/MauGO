package main

import (
	"os"
    "time"
    "database/sql"
	"log"
	"net/http"
    "github.com/allegro/bigcache"
)

var db *sql.DB
var cache *bigcache.BigCache

func main() {
    var err error
    var initErr error

    router := NewRouter()

    port := os.Getenv("PORT")

    if port == "" {
        log.Fatal("PORT environment variable was not set")
	}

    db, err = sql.Open("mysql", "mau_db:U28vX9uXJDnvmY6W@tcp(v1db.mobincube.com:3306)/mau_db")
    db.SetMaxOpenConns(100)

    if err != nil {
        log.Panic(err)
    }

    config := bigcache.Config {
        // number of shards (must be a power of 2)
        Shards: 1024,
        // time after which entry can be evicted
        LifeWindow: 10 * time.Minute,
        // rps * lifeWindow, used only in initial memory allocation
        MaxEntriesInWindow: 1000 * 10 * 60,
        // max entry size in bytes, used only in initial memory allocation
        MaxEntrySize: 500,
        // prints information about additional memory allocation
        Verbose: true,
        // cache will not allocate more memory than this limit, value in MB
        // if value is reached then the oldest entries can be overridden for the new ones
        // 0 value means no size limit
        HardMaxCacheSize: 8192,
        // callback fired when the oldest entry is removed because of its
        // expiration time or no space left for the new entry. Default value is nil which
        // means no callback and it prevents from unwrapping the oldest entry.
        OnRemove: nil,
    }

    cache, initErr = bigcache.NewBigCache(config)
    cache.Set("init", []byte("init"))

    if initErr != nil {
        log.Fatal(initErr)
    }

    log.Fatal(http.ListenAndServe(":" + port, router))
    defer db.Close()
}
