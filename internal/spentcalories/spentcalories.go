package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)
//parseTraining парсит строку, переводит данные из строки в соответствующие типы и возвращает эти значения.
func parseTraining(data string) (int, string, time.Duration, error) {
	splitedString := strings.Split(data, ",")
	if len(splitedString) != 3 {
		return 0, "", time.Duration(0), errors.New("incorrect amount of transferred data")
	}

	steps, err := strconv.Atoi(splitedString[0])
	if err != nil {
		return 0, "", time.Duration(0), err
	}
	if steps <= 0 {
		return 0, "", time.Duration(0), errors.New("incorrect number of steps")
	}

	duration, err := time.ParseDuration(splitedString[2])
	if err != nil {
		return 0, "", time.Duration(0), err
	}
	if duration <= 0 {
		return 0, "", time.Duration(0), errors.New("incorrect duration of activity")
	}
	return steps, splitedString[1], duration, err
}

//distance принимает количество шагов и рост пользователя в метрах, а возвращает дистанцию в километрах.
func distance(steps int, height float64) float64 {
	lenghtOfStep := height * stepLengthCoefficient
	return float64(steps) * lenghtOfStep / float64(mInKm)
}

//meanSpeed принимает количество шагов, рост пользователя и продолжительность активности, после чего возвращает среднюю скорость.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distanceOfActivity := distance(steps, height)
	return distanceOfActivity / duration.Hours()
}

//TrainingInfo возвращает информацию о тренировке.
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, typeOfActivity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	switch typeOfActivity {
	case "Бег":
		distanceOfActivity, avgSpeed := distance(steps, height), meanSpeed(steps, height, duration)
		spentCalories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", typeOfActivity, duration.Hours(), distanceOfActivity, avgSpeed, spentCalories), nil

	case "Ходьба":
		distanceOfActivity, avgSpeed := distance(steps, height), meanSpeed(steps, height, duration)
		spentCalories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", typeOfActivity, duration.Hours(), distanceOfActivity, avgSpeed, spentCalories), nil
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
}

//RunningSpentCalories возвращает количество потраченных калорий при беге.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	//Мало ли на ребёнка повесят трекер)
	switch {
	case weight <= 2:
		return 0, errors.New("incorrect weight")
	case height <= 0.5:
		return 0, errors.New("incorrect height")
	case duration <= 0:
		return 0, errors.New("incorrect duration of activity")
	case steps <= 0:
		return 0, errors.New("incorrect number of steps")
	}
	avgSpeed := meanSpeed(steps, height, duration)

	return (weight * avgSpeed * duration.Minutes()) / minInH, nil
}

//WalkingSpentCalories возвращает количество потраченных калорий при прогулке.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	switch {
	case weight <= 2:
		return 0, errors.New("incorrect weight")
	case height <= 0.5:
		return 0, errors.New("incorrect height")
	case duration <= 0:
		return 0, errors.New("incorrect duration of activity")
	case steps <= 0:
		return 0, errors.New("incorrect number of steps")
	}
	avgSpeed := meanSpeed(steps, height, duration)

	return (weight * avgSpeed * duration.Minutes() * walkingCaloriesCoefficient) / minInH, nil
}
