package db

import (
	"context"
	"errors"
	"slices"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrTimezoneAlreadyExists = errors.New("timezone already exists")
	ErrNoTimezonesFound      = errors.New("no timezones found")
	ErrTimezoneDoesntExist   = errors.New("timezone doesn't exist")
)

type Timezone struct {
	gorm.Model
	ChatID      int64 `gorm:"index:idx_chatid_userid,unique"`
	UserID      int64 `gorm:"index:idx_chatid_userid,unique"`
	Name        string
	Timezone    string
	Description string
}

func (db *DB) NewTimezone(ctx context.Context, chatID, userID int64, name, tz, description string) error {
	_, err := db.createTimezone(ctx, chatID, userID, name, tz, description)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetTimezones(ctx context.Context, chatID int64) ([]Timezone, error) {
	t, err := db.getTimezones(ctx, chatID)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (db *DB) RemoveTimezone(ctx context.Context, chatID, userID int64) error {
	t, err := db.getTimezone(chatID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTimezoneDoesntExist
		}
		return err
	}

	if err := db.removeTimezone(ctx, chatID, &t); err != nil {
		return err
	}

	return nil
}

func (db *DB) getTimezones(ctx context.Context, chatID int64) ([]Timezone, error) {
	t, err := db.timezoneCache.Get(ctx, chatID)
	if err == nil {
		return t, nil
	}

	res := db.gormdb.Where(&Timezone{ChatID: chatID}).Find(&t)
	if res.Error != nil {
		return t, res.Error
	}

	if err := db.timezoneCache.Set(ctx, chatID, t); err != nil {
		return t, err
	}

	return t, nil
}

func (db *DB) getTimezone(chatID, userID int64) (Timezone, error) {
	var t Timezone

	res := db.gormdb.Where(&Timezone{ChatID: chatID, UserID: userID}).First(&t)
	if res.Error != nil {
		return t, res.Error
	}

	return t, nil
}

func (db *DB) createTimezone(ctx context.Context, chatID int64, userID int64, name, tz, description string) (Timezone, error) {
	tzs, err := db.getTimezones(ctx, chatID)
	if err != nil {
		return Timezone{}, err
	}

	t := Timezone{
		ChatID:      chatID,
		UserID:      userID,
		Name:        name,
		Timezone:    tz,
		Description: description,
	}

	res := db.gormdb.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"name", "timezone", "description"}),
	}).Create(&t)
	if res.Error != nil {
		return t, res.Error
	}

	i := slices.IndexFunc[[]Timezone](tzs, func(innerT Timezone) bool {
		return userID == innerT.UserID
	})
	if i == -1 {
		if err := db.timezoneCache.Set(ctx, chatID, append(tzs, t)); err != nil {
			return t, err
		}
	} else {
		tzs[i] = t

		if err := db.timezoneCache.Set(ctx, chatID, tzs); err != nil {
			return t, err
		}
	}

	return t, nil
}

func (db *DB) removeTimezone(ctx context.Context, chatID int64, t *Timezone) error {
	tzs, err := db.getTimezones(ctx, chatID)
	if err != nil {
		return err
	}

	res := db.gormdb.Unscoped().Delete(t)
	if res.Error != nil {
		return res.Error
	}

	i := slices.IndexFunc[[]Timezone](tzs, func(innerT Timezone) bool {
		return t.UserID == innerT.UserID
	})
	if i == -1 {
		return nil
	}

	tzs[i] = tzs[len(tzs)-1]
	tzs = tzs[:len(tzs)-1]

	if err := db.timezoneCache.Set(ctx, chatID, tzs); err != nil {
		return err
	}

	return nil
}
