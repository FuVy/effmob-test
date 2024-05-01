package someapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/fuvy/effmob-test/internal/storage"
	"github.com/google/uuid"
)

var apiUrl string

type Person struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

type CarResponse struct {
	RegNum string `json:"reg_num" example:"X123XX150"`
	Mark   string `json:"mark" example:"Lada"`
	Model  string `json:"model" example:"Vesta"`
	Year   int    `json:"year" example:"2002"`
	Owner  Person `json:"owner"`
}

type CarErr struct {
	RegNum string
	Error  error
}

func SetUrl(url string) {
	apiUrl = url
}

func GetCarInfo(regNum string, wg *sync.WaitGroup, ch chan *CarResponse, chErr chan *CarErr) {
	defer wg.Done()
	queryParam := url.Values{}
	queryParam.Add("regNum", regNum)
	urlWithParams := apiUrl + "/info?" + queryParam.Encode()
	req, err := http.NewRequest("GET", urlWithParams, nil)
	if err != nil {
		chErr <- &CarErr{regNum, fmt.Errorf("creting car info request: %w", err)}
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		chErr <- &CarErr{regNum, fmt.Errorf("sending car info request: %v", err)}
		return
	}
	defer resp.Body.Close()

	carData := &CarResponse{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		chErr <- &CarErr{regNum, fmt.Errorf("car info reading body: %w", err)}
		return
	}
	err = json.Unmarshal(body, carData)
	if err != nil {
		chErr <- &CarErr{regNum, fmt.Errorf("parsing car info: %w", err)}
		return
	}
	ch <- carData
}

func (resp *CarResponse) ToCarData() (*storage.Car, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	car := &storage.Car{}
	car.ID = id
	car.Mark = resp.Mark
	car.Model = resp.Model
	car.OwnerName = resp.Owner.Name
	if resp.Owner.Patronymic != "" {
		car.OwnerPatronymic = &resp.Owner.Patronymic
	}
	car.OwnerSurname = resp.Owner.Surname
	car.RegNum = resp.RegNum
	if resp.Year != 0 {
		car.Year = &resp.Year
	}
	return car, nil
}
