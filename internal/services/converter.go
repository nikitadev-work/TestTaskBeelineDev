package services

import (
	"TestTaskBeelineDev/internal/models"
	"errors"
)

// Converter предоставляет методы для преобразования данных пользователей
type Converter struct {
}

// ConvertUser преобразует входные данные пользователя (UserInput)
// в выходной формат (UserOutput) с определением возрастной группы
func (c Converter) ConvertUser(user models.UserInput) (models.UserOutput, error) {
	ageGroup, err := determineAgeGroup(user.Age)
	if err != nil {
		return models.UserOutput{}, err
	}

	convertedUser := models.UserOutput{
		UserID:   user.UserID,
		FullName: user.Name,
		Email:    user.Email,
		AgeGroup: ageGroup,
	}

	return convertedUser, nil
}

// determineAgeGroup определяет возрастную группу по количеству лет
//
// Логика группировки:
//   - Менее 25 лет → "до 25"
//   - От 25 до 35 лет включительно → "от 25 до 35"
//   - Более 35 лет → "старше 35"
func determineAgeGroup(age int) (string, error) {
	if age < 1 || age > 150 {
		return "", errors.New("Invalid age")
	}

	if age < 25 {
		return "до 25", nil
	} else if age <= 35 {
		return "от 25 до 35", nil
	} else {
		return "старше 35", nil
	}
}
