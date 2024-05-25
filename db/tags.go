package db

import (
	"context"
	"errors"
	"slices"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrTagAlreadyExists  = errors.New("tag already exists")
	ErrNoTagsFound       = errors.New("no tags found")
	ErrTagDoesntExist    = errors.New("tag doesn't exist")
	ErrEmptyTag          = errors.New("tag is now empty, removing it")
	ErrUsersAlreadyInTag = errors.New("provided users are already mentioned by the tag")
	ErrUsersNotInTag     = errors.New("provided users are not mentioned by the tag")
)

type Tag struct {
	gorm.Model
	ChatID      int64  `gorm:"index:idx_tags_chatid_name,unique"`
	Name        string `gorm:"index:idx_tags_chatid_name,unique"`
	Description string
	Mentions    string
}

func (db *DB) NewTag(ctx context.Context, chatID int64, name, description string, mentions ...string) error {
	t, err := db.getTag(chatID, name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if t.Name != "" {
		return ErrTagAlreadyExists
	}

	_, err = db.createTag(ctx, chatID, name, description, mentions...)
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
	oldTag, err := db.getTag(chatID, name)
	if err != nil {
		return err
	}

	oldMentions := strings.Fields(oldTag.Mentions)
	oldMentionsLower := strings.Fields(strings.ToLower(oldTag.Mentions))
	l := len(oldMentions)
	for _, v := range mentions {
		if !slices.Contains[[]string](oldMentionsLower, strings.ToLower(v)) {
			oldMentions = append(oldMentions, v)
		}
	}

	if len(oldMentions) == l {
		return ErrUsersAlreadyInTag
	}

	newTag := oldTag
	newTag.Mentions = strings.Join(oldMentions, " ")

	_, err = db.updateTag(ctx, oldTag, newTag)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) RemoveMentionsFromTag(ctx context.Context, chatID int64, name string, mentions ...string) error {
	mentionsLower := strings.Fields(strings.ToLower(strings.Join(mentions, " ")))

	oldTag, err := db.getTag(chatID, name)
	if err != nil {
		return err
	}

	var newMentions []string
	oldMentions := strings.Fields(oldTag.Mentions)
	for _, v := range oldMentions {
		if i := slices.Index[[]string](mentionsLower, strings.ToLower(v)); i == -1 {
			newMentions = append(newMentions, v)
		}
	}

	if len(newMentions) == 0 {
		return ErrEmptyTag
	}

	if len(newMentions) == len(oldMentions) {
		return ErrUsersNotInTag
	}

	newTag := oldTag
	newTag.Mentions = strings.Join(newMentions, " ")

	_, err = db.updateTag(ctx, oldTag, newTag)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) RenameTag(ctx context.Context, chatID int64, oldName, newName string) error {
	oldTag, err := db.getTag(chatID, oldName)
	if err != nil {
		return err
	}

	newTag := oldTag
	newTag.Name = newName

	_, err = db.updateTag(ctx, oldTag, newTag)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) ChangeDescriptionOfTag(ctx context.Context, chatID int64, name, description string) error {
	oldTag, err := db.getTag(chatID, name)
	if err != nil {
		return err
	}

	newTag := oldTag
	newTag.Description = description

	_, err = db.updateTag(ctx, oldTag, newTag)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) getTags(ctx context.Context, chatID int64) ([]Tag, error) {
	t, err := db.tagCache.Get(ctx, chatID)
	if err == nil {
		return t, nil
	}

	res := db.gormdb.Where(&Tag{ChatID: chatID}).Find(&t)
	if res.Error != nil {
		return t, res.Error
	}

	if err := db.tagCache.Set(ctx, chatID, t); err != nil {
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

func (db *DB) createTag(ctx context.Context, chatID int64, name, description string, mentions ...string) (Tag, error) {
	tags, err := db.getTags(ctx, chatID)
	if err != nil {
		return Tag{}, err
	}

	t := Tag{
		ChatID:      chatID,
		Name:        name,
		Description: description,
		Mentions:    strings.Join(mentions, " "),
	}

	res := db.gormdb.Create(&t)
	if res.Error != nil {
		return t, res.Error
	}

	if err := db.tagCache.Set(ctx, chatID, append(tags, t)); err != nil {
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

	if err := db.tagCache.Set(ctx, chatID, tags); err != nil {
		return err
	}

	return nil
}

func (db *DB) updateTag(ctx context.Context, oldTag, newTag Tag) (Tag, error) {
	tags, err := db.getTags(ctx, oldTag.ChatID)
	if err != nil {
		return newTag, err
	}

	res := db.gormdb.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&newTag)
	if res.Error != nil {
		return newTag, res.Error
	}

	i := slices.IndexFunc[[]Tag](tags, func(innerT Tag) bool {
		return oldTag.Name == innerT.Name
	})
	if i == -1 {
		return newTag, res.Error
	}

	tags[i] = newTag

	if err := db.tagCache.Set(ctx, newTag.ChatID, tags); err != nil {
		return newTag, err
	}

	return newTag, nil
}
