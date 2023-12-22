package db

import (
	"context"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	"github.com/glebarez/sqlite"
	slogGorm "github.com/orandin/slog-gorm"
	gocache "github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/vyneer/tg-tagbot/config"
)

var (
	ErrTagAlreadyExists  = errors.New("tag already exists")
	ErrNoTagsFound       = errors.New("no tags found")
	ErrTagDoesntExist    = errors.New("tag doesn't exist")
	ErrEmptyTag          = errors.New("tag is now empty, removing it")
	ErrUsersAlreadyInTag = errors.New("provided users are already mentioned by the tag")
	ErrUsersNotInTag     = errors.New("provided users are not mentioned by the tag")
)

type DB struct {
	cache  *cache.Cache[[]Tag]
	gormdb *gorm.DB
}

type Chat struct {
	ID int64 `gorm:"uniqueIndex"`
}

type Tag struct {
	gorm.Model
	ChatID   int64  `gorm:"uniqueIndex"`
	Name     string `gorm:"uniqueIndex"`
	Mentions string
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

	gocacheClient := gocache.New(10*time.Minute, 15*time.Minute)
	gocacheStore := gocache_store.NewGoCache(gocacheClient)
	cacheManager := cache.New[[]Tag](gocacheStore)

	return DB{
		gormdb: db,
		cache:  cacheManager,
	}, nil
}

func (db *DB) NewTag(ctx context.Context, chatID int64, name string, mentions ...string) error {
	t, err := db.getTag(chatID, name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if t.Name != "" {
		return ErrTagAlreadyExists
	}

	_, err = db.createTag(ctx, chatID, name, mentions...)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetTags(ctx context.Context, chatID int64) ([]Tag, error) {
	t, err := db.getTags(ctx, chatID)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (db *DB) RemoveTag(ctx context.Context, chatID int64, name string) error {
	t, err := db.getTag(chatID, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTagDoesntExist
		}
		return err
	}

	if err := db.removeTag(ctx, chatID, &t); err != nil {
		return err
	}

	return nil
}

func (db *DB) AddMentionsToTag(ctx context.Context, chatID int64, name string, mentions ...string) error {
	t, err := db.getTag(chatID, name)
	if err != nil {
		return err
	}

	oldMentions := strings.Fields(t.Mentions)
	l := len(oldMentions)
	for _, v := range mentions {
		if !slices.Contains[[]string](oldMentions, v) {
			oldMentions = append(oldMentions, v)
		}
	}

	if len(oldMentions) == l {
		return ErrUsersAlreadyInTag
	}

	t.Mentions = strings.Join(oldMentions, " ")

	_, err = db.updateTag(ctx, t)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) RemoveMentionsFromTag(ctx context.Context, chatID int64, name string, mentions ...string) error {
	t, err := db.getTag(chatID, name)
	if err != nil {
		return err
	}

	var newMentions []string
	oldMentions := strings.Fields(t.Mentions)
	for _, v := range oldMentions {
		if i := slices.Index[[]string](mentions, v); i == -1 {
			newMentions = append(newMentions, v)
		}
	}

	if len(newMentions) == 0 {
		return ErrEmptyTag
	}

	if len(newMentions) == len(oldMentions) {
		return ErrUsersNotInTag
	}

	t.Mentions = strings.Join(newMentions, " ")

	_, err = db.updateTag(ctx, t)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) getTags(ctx context.Context, chatID int64) ([]Tag, error) {
	t, err := db.cache.Get(ctx, chatID)
	if err == nil {
		return t, nil
	}

	res := db.gormdb.Where(&Tag{ChatID: chatID}).Find(&t)
	if res.Error != nil {
		return t, res.Error
	}

	if err := db.cache.Set(ctx, chatID, t); err != nil {
		return t, err
	}

	return t, nil
}

func (db *DB) getTag(chatID int64, name string) (Tag, error) {
	var t Tag

	res := db.gormdb.Where(&Tag{ChatID: chatID, Name: name}).First(&t)
	if res.Error != nil {
		return t, res.Error
	}

	return t, nil
}

func (db *DB) createTag(ctx context.Context, chatID int64, name string, mentions ...string) (Tag, error) {
	tags, err := db.getTags(ctx, chatID)
	if err != nil {
		return Tag{}, err
	}

	t := Tag{
		ChatID:   chatID,
		Name:     name,
		Mentions: strings.Join(mentions, " "),
	}

	res := db.gormdb.Create(&t)
	if res.Error != nil {
		return t, res.Error
	}

	if err := db.cache.Set(ctx, chatID, append(tags, t)); err != nil {
		return t, err
	}

	return t, nil
}

func (db *DB) removeTag(ctx context.Context, chatID int64, t *Tag) error {
	tags, err := db.getTags(ctx, chatID)
	if err != nil {
		return err
	}

	res := db.gormdb.Unscoped().Delete(t)
	if res.Error != nil {
		return res.Error
	}

	i := slices.IndexFunc[[]Tag](tags, func(innerT Tag) bool {
		return t.Name == innerT.Name
	})
	if i == -1 {
		return nil
	}

	tags[i] = tags[len(tags)-1]
	tags = tags[:len(tags)-1]

	if err := db.cache.Set(ctx, chatID, tags); err != nil {
		return err
	}

	return nil
}

func (db *DB) updateTag(ctx context.Context, t Tag) (Tag, error) {
	tags, err := db.getTags(ctx, t.ChatID)
	if err != nil {
		return t, err
	}

	res := db.gormdb.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&t)
	if res.Error != nil {
		return t, res.Error
	}

	i := slices.IndexFunc[[]Tag](tags, func(innerT Tag) bool {
		return t.Name == innerT.Name
	})
	if i == -1 {
		return t, res.Error
	}

	tags[i] = t

	if err := db.cache.Set(ctx, t.ChatID, tags); err != nil {
		return t, err
	}

	return t, nil
}
