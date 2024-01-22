package personapi

import (
	"encoding/json"
	"go-jun/internal/person/entity"
	"net/http"
)

func CallExternalAPI(name string) (age int, gender string, nationality string, err error) {
	ageUrl := "https://api.agify.io/?name=" + name
	resp, err := http.Get(ageUrl)
	if err != nil {
		return 0, "", "", err
	}
	var ageResponse entity.AgeResponse
	err = json.NewDecoder(resp.Body).Decode(&ageResponse)
	if err != nil {
		return 0, "", "", err
	}

	err = resp.Body.Close()
	if err != nil {
		return 0, "", "", err
	}

	genderUrl := "https://api.genderize.io/?name=" + name
	resp, err = http.Get(genderUrl)
	if err != nil {
		return 0, "", "", err
	}

	var genderResponse entity.GenderResponse
	err = json.NewDecoder(resp.Body).Decode(&genderResponse)
	if err != nil {
		return 0, "", "", err
	}
	err = resp.Body.Close()
	if err != nil {
		return 0, "", "", err
	}

	nationalityUrl := "https://api.nationalize.io/?name=" + name
	resp, err = http.Get(nationalityUrl)
	if err != nil {
		return 0, "", "", err
	}
	var nationalityResponse entity.NationalizeResponse
	err = json.NewDecoder(resp.Body).Decode(&nationalityResponse)
	if err != nil {
		return 0, "", "", err
	}
	err = resp.Body.Close()
	if err != nil {
		return 0, "", "", err
	}

	return ageResponse.Age, genderResponse.Gender, nationalityResponse.Country[0].CountryID, nil
}
