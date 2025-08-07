package services

import (
	"TestTaskBeelineDev/internal/models"
	"reflect"
	"testing"
)

// Пример unit тестов для функций
func TestConverter_ConvertUser(t *testing.T) {

	input := models.UserInput{
		UserID: "1",
		Name:   "Иван Иванов",
		Email:  "ivan@example.com",
		Age:    30,
	}

	output := models.UserOutput{
		UserID:   "1",
		FullName: "Иван Иванов",
		Email:    "ivan@example.com",
		AgeGroup: "от 25 до 35",
	}

	converter := Converter{}
	result, err := converter.ConvertUser(input)
	if err != nil {
		t.Errorf("Convert user input error: %v", err)
		return
	}

	if result != output {
		t.Errorf("Convert user output error: expected %v, got %v", output, result)
		return
	}
}

// Unit тест для функции determineAgeGroup
func TestConverter_DetermineAgeGroup(t *testing.T) {
	input := 30
	output := "от 25 до 35"

	result, err := determineAgeGroup(input)
	if err != nil {
		t.Errorf("DetermineAgeGroup error: %v", err)
		return
	}

	if !reflect.DeepEqual(result, output) {
		t.Errorf("DetermineAgeGroup output error: expected %v, got %v", output, result)
		return
	}
}
