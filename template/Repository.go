package repository

import (
	"[project]/entity"
	"time"

	"gorm.io/gorm"
)

type [nameUpper]Repository interface {
	Save([name] entity.[nameUpper]) (entity.[nameUpper], error)
	Edit([name] entity.[nameUpper]) (entity.[nameUpper], error)
	[findBy]
	FindAll() ([]entity.[nameUpper], error)
	Delete(id int) (string, error)
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

[findByMethod]

func (r *[name]Repository) FindAll() ([]entity.[nameUpper], error) {
	var [name]s []entity.[nameUpper]

	err := r.db.Where("deleted_date = ?", nil).Find(&[name]s).Error

	if err != nil {
		return [name]s, err
	}

	return [name]s, nil
}

func (r *[name]Repository) Delete(id int) (string, error) {
	[name], err := r.FindById(id)

	if err != nil {
		return "Failed", err
	}

	[name].DeletedDate = time.Now()
	err = r.db.Save([name]).Error

	if err != nil {
		return "Failed", err
	}
	return "Success", nil
}
