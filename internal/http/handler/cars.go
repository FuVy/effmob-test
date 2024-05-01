package handler

import (
	"errors"
	"net/http"
	"sync"

	"github.com/fuvy/effmob-test/internal/service"
	"github.com/fuvy/effmob-test/internal/storage"
	someapi "github.com/fuvy/effmob-test/pkg/someAPI"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type RegisterRequest struct {
	RegNums []string `json:"regNums"`
} // @name RegisterRequest

// @Summary Get cars
// @Description Retrieve a list of cars based on provided queries
// @ID get-cars
// @Tags Car
// @Produce json
// @Param reg_num query string false "Registration number of the car"
// @Param mark query string false "Car manufacturer mark"
// @Param model query string false "Car model"
// @Param year query string false "Car manufacturing year"
// @Param owner_name query string false "Owner's name"
// @Param owner_surname query string false "Owner's surname"
// @Param owner_patronymic query string false "Owner's patronymic"
// @Success 200 {object} ApiResponse
// @Failure 500
// @Router /v1/cars [get]
func (h *Handler) GetCars(ctx *fiber.Ctx) error {
	queries := ctx.Queries()
	cars, err := h.service.GetCars(queries)
	if err != nil {
		return internalError(ctx, err)
	}
	return respondJSON(ctx, http.StatusOK, "received cars successfully", "", cars)
}

// @Summary Delete car
// @Description Deletes car by id
// @Tags Car
// @ID delete-car
// @Produce json
// @Param id path string true "Car ID"
// @Success 200 {object} ApiResponse
// @Failure 400 {string} ApiResponse{error=ApiCodeResponse}
// @Failure 404 {string} ApiResponse{error=ApiCodeResponse}
// @Failure 500
// @Router	/v1/cars/{id} [delete]
func (h *Handler) DeleteCar(ctx *fiber.Ctx) error {
	var car storage.Car
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return writeApiError(ctx, http.StatusBadRequest, errCarInvalidID)
	}
	if err := h.service.DeleteCar(&car, id); err != nil {
		if errors.Is(err, service.ErrCarNotFound) {
			return writeApiError(ctx, http.StatusNotFound, errCarNotFound)
		}
		return internalError(ctx, err)
	}
	return respondMessage(ctx, "car deleted successfully")
}

// @Summary Register cars
// @Description Register multiple cars by fetching their information from an external API
// @ID register-cars
// @Tags Car
// @Accept json
// @Produce json
// @Param	json body RegisterRequest true "Request body"
// @Success 200 {object} ApiResponse
// @Failure 400 {string} ApiResponse{error=ApiCodeResponse}
// @Failure 404 {string} ApiResponse{error=ApiCodeResponse}
// @Failure 500
// @Router /v1/cars/ [post]
func (h *Handler) RegisterCars(ctx *fiber.Ctx) error {
	req := &RegisterRequest{}
	if err := ctx.BodyParser(req); err != nil {
		return writeInvalidRequest(ctx)
	}
	if len(req.RegNums) == 0 {
		return writeApiError(ctx, http.StatusBadRequest, errNoCarsProvided)
	}
	var wg sync.WaitGroup
	wg.Add(len(req.RegNums))
	ch := make(chan *someapi.CarResponse, len(req.RegNums))
	chErr := make(chan *someapi.CarErr, len(req.RegNums))
	carResponses := make([]*someapi.CarResponse, 0, len(req.RegNums))
	for _, regNum := range req.RegNums {
		go someapi.GetCarInfo(regNum, &wg, ch, chErr)
	}

	go func() {
		wg.Wait()
		close(ch)
		close(chErr)
	}()

	var mu sync.Mutex

	for carResponse := range ch {
		mu.Lock()
		carResponses = append(carResponses, carResponse)
		mu.Unlock()
	}

	for err := range chErr {
		log.Errorf("Failed to get car info for %q : %v", err.RegNum, err.Error)
	}

	carsData := make([]storage.Car, 0, len(carResponses))
	for _, v := range carResponses {
		res, err := v.ToCarData()
		if err == nil {
			carsData = append(carsData, *res)
		}
	}

	if len(carsData) == 0 {
		return writeApiError(ctx, http.StatusBadRequest, errUnknownCars)
	}

	err := h.service.AddCars(carsData)
	if err != nil {
		if errors.Is(err, service.ErrCarAlreadyExists) {
			return writeApiError(ctx, http.StatusBadRequest, errCarAlreadyAdded)
		}
		return internalError(ctx, err)
	}
	return respondJSON(ctx, http.StatusOK, "added cars successfully", "", carsData)
}

// @Summary		Update Car Info
// @Description	Changes one or several car fields
// @ID update-car
// @Tags		Car
// @Accept		json
// @Produce		json
// @Param 		id path string true "Car ID"
// @Param 		body body CarObject true "Car object to update"
// @Success		200 {object}	ApiResponse
// @Failure		400	{object}	ApiResponse{error=ApiCodeResponse}
// @Failure		404	{object}	ApiResponse{error=ApiCodeResponse}
// @Failure		500
// @Router		/v1/cars/{id} [patch]
func (h *Handler) UpdateCar(ctx *fiber.Ctx) error {
	var car storage.Car
	if err := ctx.BodyParser(&car); err != nil {
		return writeInvalidRequest(ctx)
	}
	var fields map[string]any
	toChange := make(map[string]any)
	if err := ctx.BodyParser(&fields); err != nil {
		return writeInvalidRequest(ctx)
	}
	if _, ok := fields["reg_num"]; ok {
		toChange["reg_num"] = car.RegNum
	}
	if _, ok := fields["mark"]; ok {
		toChange["mark"] = car.Mark
	}
	if _, ok := fields["model"]; ok {
		toChange["model"] = car.Model
	}
	if _, ok := fields["year"]; ok {
		toChange["year"] = car.Year
	}
	if _, ok := fields["owner_name"]; ok {
		toChange["owner_name"] = car.OwnerName
	}
	if _, ok := fields["owner_surname"]; ok {
		toChange["owner_surname"] = car.OwnerSurname
	}
	if _, ok := fields["owner_patronymic"]; ok {
		toChange["owner_patronymic"] = car.OwnerPatronymic
	}
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return writeApiError(ctx, http.StatusBadRequest, errCarInvalidID)
	}
	newCar, err := h.service.UpdateCar(id, toChange)
	if err != nil {
		if errors.Is(err, service.ErrCarNotFound) {
			return writeApiError(ctx, http.StatusNotFound, errCarNotFound)
		}
		return internalError(ctx, err)
	}
	return respondJSON(ctx, http.StatusOK, "car updated successfully", "", *newCar)
}
