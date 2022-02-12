package store

import (
	"golang-sql/model"
	"gorm.io/gorm"
)

type ArticleStore struct {
	db *gorm.DB
}

func NewArticleStore(db *gorm.DB) *ArticleStore {
	return &ArticleStore{
		db: db,
	}
}

func (as *ArticleStore) GetBySlug(slug string) (*model.Article, error) {
	var m model.Article

	err := as.db.Where(&model.Article{Slug: slug}).Preload("Favourites").Preload("Tags").Preload("Author").Find(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (as *ArticleStore) GetUserArticleBySlug(userId uint, slug string) (*model.Article, error) {
	var m model.Article

	err := as.db.Where(&model.Article{Slug: slug, AuthorID: userId}).Find(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (as *ArticleStore) CreateArticle(a *model.Article) error {
	tags := a.Tags

	tx := as.db.Begin()
	if err := tx.Create(&a).Error; err != nil {
		return err
	}

	for _, t := range a.Tags {
		err := tx.Where(&model.Tag{Tag: t.Tag}).First(&t).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Model(&a).Association("Tags").Append(t); err != nil {
			tx.Rollback()
			return err
		}
	}

	a.Tags = tags
	return tx.Commit().Error
}

func (as *ArticleStore) UpdateArticle(a *model.Article, tagList []string) error {
	tx := as.db.Begin()
	if err := tx.Model(a).Updates(a).Error; err != nil {
		tx.Rollback()
		return err
	}

	tags := make([]model.Tag, 0)

	for _, t := range tagList {
		tag := model.Tag{Tag: t}

		err := tx.Where(&tag).First(&tag).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		tags = append(tags, tag)
	}

	if err := tx.Model(a).Association("Tags").Replace(tags); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where(a.ID).Preload("Favorites").Preload("Tags").Preload("Author").Find(a).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (as *ArticleStore) DeleteArticle(a *model.Article) error {
	return as.db.Delete(a).Error
}