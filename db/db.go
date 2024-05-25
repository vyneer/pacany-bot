package db

import (
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	"github.com/glebarez/sqlite"
	slogGorm "github.com/orandin/slog-gorm"
	gocache "github.com/patrickmn/go-cache"
	"gorm.io/gorm"

	"github.com/vyneer/pacany-bot/config"
)

type DB struct {
	gormdb        *gorm.DB
	tagCache      *cache.Cache[[]Tag]
	timezoneCache *cache.Cache[[]Timezone]
}

func New(c *config.Config) (DB, error) {
	db, err := gorm.Open(sqlite.Open(c.DBPath), &gorm.Config{
		Logger: slogGorm.New(),
	})
	if err != nil {
		return DB{}, err
	}

	err = db.AutoMigrate(&Tag{}, &Timezone{})
	if err != nil {
		return DB{}, err
	}

	tagGocacheClient := gocache.New(10*time.Minute, 15*time.Minute)
	tagGocacheStore := gocache_store.NewGoCache(tagGocacheClient)
	tagCacheManager := cache.New[[]Tag](tagGocacheStore)

	timezoneGocacheClient := gocache.New(10*time.Minute, 15*time.Minute)
	timezoneGocacheStore := gocache_store.NewGoCache(timezoneGocacheClient)
	timezoneCacheManager := cache.New[[]Timezone](timezoneGocacheStore)

	return DB{
		gormdb:        db,
		tagCache:      tagCacheManager,
		timezoneCache: timezoneCacheManager,
	}, nil
}
