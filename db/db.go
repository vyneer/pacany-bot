package db

import (
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	"github.com/glebarez/sqlite"
	slogGorm "github.com/orandin/slog-gorm"
	gocache "github.com/patrickmn/go-cache"
	"gorm.io/gorm"

	"github.com/vyneer/tg-tagbot/config"
)

type DB struct {
	gormdb   *gorm.DB
	tagCache *cache.Cache[[]Tag]
}

func New(c *config.Config) (DB, error) {
	db, err := gorm.Open(sqlite.Open(c.DBPath), &gorm.Config{
		Logger: slogGorm.New(),
	})
	if err != nil {
		return DB{}, err
	}

	err = db.AutoMigrate(&Tag{})
	if err != nil {
		return DB{}, err
	}

	tagGocacheClient := gocache.New(10*time.Minute, 15*time.Minute)
	tagGocacheStore := gocache_store.NewGoCache(tagGocacheClient)
	tagCacheManager := cache.New[[]Tag](tagGocacheStore)

	return DB{
		gormdb:   db,
		tagCache: tagCacheManager,
	}, nil
}
