package service

import (
	"errors"
	"fmt"

	"github.com/fuvy/effmob-test/internal/storage"
	"gorm.io/gorm"
)

var (
	pageLimit = 7
)

var (
	ErrCarNotFound      = fmt.Errorf("car not found")
	ErrCarAlreadyExists = fmt.Errorf("car already exists")
)

func (s *Service) GetAllCars(page int) ([]storage.Car, error) {
	var cars []storage.Car
	offset := (page - 1) * pageLimit
	err := s.db.Limit(pageLimit).Offset(offset).Find(&cars).Error
	if err != nil {
		return []storage.Car{}, err
	}
	return cars, nil
}

func (s *Service) DeleteCar(car *storage.Car, id string) error {
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

func (s *Service) UpdateCar(updatedCar *storage.Car) (*storage.Car, error) {
	var car storage.Car
	tx := s.db.Take(&car, "id = ?", updatedCar.ID)
	err := tx.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCarNotFound
		}
		return nil, err
	}
	if updatedCar.RegNum != "" {
		car.RegNum = updatedCar.RegNum
	}
	if updatedCar.Mark != "" {
		car.Mark = updatedCar.Mark
	}
	if updatedCar.Model != "" {
		car.Model = updatedCar.Model
	}
	if updatedCar.Year != nil {
		car.Year = updatedCar.Year
	}
	if updatedCar.OwnerName != "" {
		car.OwnerName = updatedCar.OwnerName
	}
	if updatedCar.OwnerSurname != "" {
		car.OwnerSurname = updatedCar.OwnerSurname
	}
	if updatedCar.OwnerPatronymic != nil {
		car.OwnerPatronymic = updatedCar.OwnerPatronymic
	}
	return &car, tx.Save(&car).Error
}

func (s *Service) AddCars(cars []storage.Car) error {
	tx := s.db.Begin()
	defer tx.Rollback()
	if err := tx.Create(cars).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrCarAlreadyExists
		}
		return err
	}
	return tx.Commit().Error
}
