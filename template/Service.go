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
	Get[nameUpper]ById(id string) (entity.[nameUpper], error)
	GetAll[nameUpper]() ([]entity.[nameUpper], error)
}

type [name]Service struct {
	[name]Repository repository.[nameUpper]Repository
}

func New[nameUpper]Service([name]Repository repository.[nameUpper]Repository) *[name]Service {
	return &[name]Service{[name]Repository}
}

func (s *[name]Service) Create[nameUpper](input input.[nameUpper]Input) (entity.[nameUpper], error) {
	[name] := entity.[nameUpper]{
		[createItem]
	}

	new[nameUpper], err := s.[name]Repository.Save([name])

	if err != nil {
		return [name], err
	}

	return new[nameUpper], nil
}

func (s *[name]Service) Edit[nameUpper](input input.[nameUpper]EditInput) (entity.[nameUpper], error) {
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

func (s *[name]Service) Get[nameUpper]ById(id string) (entity.[nameUpper], error) {
	idint, err := strconv.Atoi(id)

	if err != nil {
		return entity.[nameUpper]{}, err
	}

	[name], err := s.[name]Repository.FindById(idint)

	if err != nil {
		return [name], err
	}

	if [name].Id == 0 {
		return [name], errors.New("[name] not found")
	}

	return [name], nil
}

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
