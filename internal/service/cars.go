package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/fuvy/effmob-test/internal/storage"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	pageLimit = 7
)

var (
	ErrCarNotFound            = fmt.Errorf("car not found")
	ErrCarAlreadyExists       = fmt.Errorf("car already exists")
	ErrSomeCarsAlreadyExisted = fmt.Errorf("some cars already existed")
)

func (s *Service) GetCars(q map[string]string) ([]storage.Car, error) {
	query := s.db
	page := 1
	if v, ok := q["reg_num"]; ok {
		query = query.Where("reg_num = ?", v)
	}
	if v, ok := q["mark"]; ok {
		query = query.Where("mark = ?", v)
	}
	if v, ok := q["model"]; ok {
		query = query.Where("model = ?", v)
	}
	if v, ok := q["year"]; ok {
		query = query.Where("year = ?", v)
	}
	if v, ok := q["owner_name"]; ok {
		query = query.Where("owner_name LIKE ?", v)
	}
	if v, ok := q["owner_surname"]; ok {
		query = query.Where("owner_surname LIKE ?", v)
	}
	if v, ok := q["owner_patronymic"]; ok {
		query = query.Where("owner_name LIKE ?", v)
	}
	if v, ok := q["page"]; ok {
		page, err := strconv.Atoi(v)
		if err != nil {
			page = 1
		}
		if page < 1 {
			page = 1
		}
	}

	var cars []storage.Car
	offset := (page - 1) * pageLimit
	err := query.Limit(pageLimit).Offset(offset).Find(&cars).Error
	if err != nil {
		return []storage.Car{}, err
	}
	return cars, nil
}

func (s *Service) DeleteCar(car *storage.Car, id uuid.UUID) error {
	tx := s.db.Take(car, "id = ?", id)
	err := tx.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCarNotFound
		}
		return err
	}
	return tx.Delete(car).Error
}

func (s *Service) UpdateCar(id uuid.UUID, m map[string]any) (*storage.Car, error) {
	var car storage.Car
	tx := s.db.Begin()
	defer tx.Rollback()
	err := tx.Take(&car, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCarNotFound
		}
		return nil, err
	}
	if err := tx.Model(&car).Updates(&m).Error; err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return &car, nil
}

func (s *Service) AddCars(cars []storage.Car) error {
	if err := s.db.Create(cars).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrCarAlreadyExists
		}
		return err
	}

	return nil
}
