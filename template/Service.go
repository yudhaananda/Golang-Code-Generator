package service

import (
	"[project]/entity"
	"[project]/input"
	"[project]/repository"
	"errors"
	"strconv"
	"time"
)

type [nameUpper]Service interface {
	Create[nameUpper](input input.[nameUpper]Input) (entity.[nameUpper], error)
	Edit[nameUpper](input input.[nameUpper]EditInput) (entity.[nameUpper], error)
	[getBy]
	GetAll[nameUpper]() ([]entity.[nameUpper], error)
	Delete[nameUpper](id string) (string, error)

}

type [name]Service struct {
	[name]Repository repository.[nameUpper]Repository
}

func New[nameUpper]Service([name]Repository repository.[nameUpper]Repository) *[name]Service {
	return &[name]Service{[name]Repository}
}

func (s *[name]Service) Create[nameUpper](input input.[nameUpper]Input, userName string) (entity.[nameUpper], error) {
	[name] := entity.[nameUpper]{
		[createItem]
	}

	new[nameUpper], err := s.[name]Repository.Save([name])

	if err != nil {
		return [name], err
	}

	return new[nameUpper], nil
}

func (s *[name]Service) Edit[nameUpper](input input.[nameUpper]EditInput, userName string) (entity.[nameUpper], error) {
	old[nameUpper], err := s.[name]Repository.FindById(input.Id)

	if err != nil {
		return old[nameUpper], err
	}

	[name] := entity.[nameUpper]{
		[editItem]
	}

	new[nameUpper], err := s.[name]Repository.Edit([name])

	if err != nil {
		return [name], err
	}

	return new[nameUpper], nil
}

[getByMethod]

func (s *[name]Service) GetAll[nameUpper]() ([]entity.[nameUpper], error) {
	[name]s, err := s.[name]Repository.FindAll()

	if err != nil {
		return [name]s, err
	}

	if len([name]s) <= 0 {
		return [name]s, errors.New("[name] not found")
	}

	return [name]s, nil
}

func (s *[name]Service) Delete[nameUpper](id int) (string, error) {

	result, err := s.[name]Repository.Delete(id)
	if err != nil {
		return result, err
	}
	return result, nil
}
