package repository

import (
	"[project]/entity"

	"gorm.io/gorm"
)

type [nameUpper]Repository interface {
	Save([name] entity.[nameUpper]) (entity.[nameUpper], error)
	Edit([name] entity.[nameUpper]) (entity.[nameUpper], error)
	FindById(id int) (entity.[nameUpper], error)
	FindAll() ([]entity.[nameUpper], error)
}

type [name]Repository struct {
	db *gorm.DB
}

func New[nameUpper]Repository(db *gorm.DB) *[name]Repository {
	return &[name]Repository{db}
}

func (r *[name]Repository) Save([name] entity.[nameUpper]) (entity.[nameUpper], error) {
	err := r.db.Create(&[name]).Error

	if err != nil {
		return [name], err
	}

	return [name], nil
}

func (r *[name]Repository) Edit([name] entity.[nameUpper]) (entity.[nameUpper], error) {
	err := r.db.Save([name]).Error

	if err != nil {
		return [name], err
	}

	return [name], nil
}

func (r *[name]Repository) FindById(id int) (entity.[nameUpper], error) {
	var [name] entity.[nameUpper]

	err := r.db.Where("id = ?", id).Find(&[name]).Error

	if err != nil {
		return [name], err
	}

	return [name], nil
}

func (r *[name]Repository) FindAll() ([]entity.[nameUpper], error) {
	var [name]s []entity.[nameUpper]

	err := r.db.Find(&[name]s).Error

	if err != nil {
		return [name]s, err
	}

	return [name]s, nil
}
