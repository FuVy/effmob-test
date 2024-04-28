package handler

import (
	"errors"
	"strconv"
	"sync"

	"github.com/fuvy/effmob-test/internal/service"
	"github.com/fuvy/effmob-test/internal/storage"
	someapi "github.com/fuvy/effmob-test/pkg/someAPI"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type RegisterRequest struct {
	RegNums []string `json:"regNums"`
}

func (h *Handler) GetCars(ctx *fiber.Ctx) error {
	pageString := ctx.Query("page", "1")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}
	cars, err := h.service.GetAllCars(page)
	if err != nil {
		return internalError(ctx, err)
	}
	return respondJSON(ctx, 200, "received cars successfully", "", cars)
}

func (h *Handler) DeleteCar(ctx *fiber.Ctx) error {
	var car storage.Car
	id := ctx.Params("id")
	err := h.service.DeleteCar(&car, id)
	if err != nil {
		if errors.Is(err, service.ErrCarNotFound) {
			return writeApiError(ctx, errCarNotFound)
		}
		return internalError(ctx, err)
	}
	return respondMessage(ctx, "car deleted successfully")
}

func (h *Handler) RegisterCars(ctx *fiber.Ctx) error {
	req := &RegisterRequest{}
	if err := ctx.BodyParser(req); err != nil {
		return writeInvalidRequest(ctx)
	}
	var wg sync.WaitGroup
	wg.Add(len(req.RegNums))
	carResponses := make([]*someapi.CarResponse, 0, len(req.RegNums))
	for _, regNum := range req.RegNums {
		go func(regNum string) {
			defer wg.Done()
			car, err := someapi.GetCarInfo(regNum)
			if err != nil {
				log.Errorf("Failed to get car info for regNum %s: %v\n", regNum, err)
				return
			}
			carResponses = append(carResponses, car)
		}(regNum)
	}

	wg.Wait()
	carsData := make([]storage.Car, 0, len(carResponses))
	for _, v := range carResponses {
		res, err := v.ToCarData()
		if err != nil {
			carsData = append(carsData, *res)
		}
	}

	err := h.service.AddCars(carsData)
	if err != nil {
		if errors.Is(err, service.ErrCarAlreadyExists) {
			return writeApiError(ctx, errCarAlreadyAdded)
		}
		return internalError(ctx, err)
	}
	return respondJSON(ctx, 200, "added cars successfully", "", carsData)
}

func (h *Handler) UpdateCar(ctx *fiber.Ctx) error {
	car := &storage.Car{}
	if err := ctx.BodyParser(car); err != nil {
		return writeInvalidRequest(ctx)
	}
	car, err := h.service.UpdateCar(car)
	if err != nil {
		if errors.Is(err, service.ErrCarNotFound) {
			return writeApiError(ctx, errCarNotFound)
		}
		return internalError(ctx, err)
	}
	return respondJSON(ctx, 200, "car info updated successfully", "", car)
}
