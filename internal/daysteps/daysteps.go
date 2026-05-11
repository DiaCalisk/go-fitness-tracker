package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"log"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// parsePackage парсит строку с шагами и временем активности и возвращает эти значения.
func parsePackage(data string) (int, time.Duration, error) {
	splitedString := strings.Split(data, ",")
	if len(splitedString) != 2 {
		return 0, time.Duration(0), errors.New("Неверное количество переданных данных")
	}
	steps, err := strconv.Atoi(splitedString[0])
	if err != nil {
		return 0, time.Duration(0), err
	}
	if steps <= 0 {
		return 0, time.Duration(0), errors.New("Неверное количество шагов")
	}
	duration, err := time.ParseDuration(splitedString[1])
	if err != nil {
		return 0, time.Duration(0), err
	}
	if duration <= 0 {
		return 0, time.Duration(0), errors.New("Неверная длительность активности")
	}
	return steps, duration, err
}

// DayActionInfo парсит строку с данными с помощью parsePackage() и вычисляет дистанцию
// в километрах и количество потраченных калорий, после его возвращает строку с данными.
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	// Не стал проверять количество шагов на 0, так как в функции parsePackage и так есть эта проверка, и при её срабатывании выводится ошибка выше

	numOfKm := float64(steps) * stepLength / mInKm
	spentCalories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, numOfKm, spentCalories)
}
