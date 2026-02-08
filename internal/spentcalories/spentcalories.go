package spentcalories

import (
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

func parseTraining(data string) (int, string, time.Duration, error) {
	// Разделить строку на слайс строк
	parts := strings.Split(data, ",")

	// Проверить, чтобы длина слайса была равна 3
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("invalid format")
	}

	firstElement := parts[0]

	// Преобразовать первый элемент слайса в тип int
	steps, err := strconv.Atoi(firstElement)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования шагов: %v", err)
	}

	if steps <= 0  {
		return 0, "", 0, fmt.Errorf("ошибка преобразования шагов: %v", err)
	}
	
	activity := parts[1]

	// Преобразовать третий элемент слайса в time.Duration
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования продолжительности: %v", err)
	}

	// Проверка на нулевую продолжительность
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("продолжительность не может быть равна нулю")
	}

	// Или более детальная проверка в минутах и часах
	if duration.Minutes() <= 0 {
		return 0, "", 0, fmt.Errorf("продолжительность не может быть 0 минут")
	}

	// Возврат количества шагов, вида активности, продолжительности и nil 
	return steps, activity, duration, nil
}


func distance(steps int, height float64) float64 {
	// Длина шага
	stepLength := height * stepLengthCoefficient

	// Вычисление дистанции
	distanceInMeters := float64(steps) * stepLength
	
	distanceInKilometers := distanceInMeters / mInKm

	return distanceInKilometers
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// Проверить, что продолжительность duration больше 0
	if duration <= 0 {
		return 0
	}

	// Вычислить дистанцию с помощью distance()
	dist := distance(steps, height)

	// Вычислить продолжительность в часах
	durationInHours := duration.Hours()

	// Вычислить и вернуть среднюю скорость
	speed := dist / durationInHours
	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// Получить значения из строки данных с помощью функции parseTraining()
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64
	var caloriesErr error

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	// Проверить, какой вид тренировки был передан в строке
	switch activity {
	case "Бег", "бег":
		calories, caloriesErr = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба", "ходьба":
		calories, caloriesErr = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	// Проверяем ошибки при расчете калорий
	if caloriesErr != nil {
		log.Println(caloriesErr)
		return "", caloriesErr
	}

	// Резльтат
	result := fmt.Sprintf(
			"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		activity,
		duration.Hours(),
		dist,
		speed,
		calories,
	)
	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверить входные параметры на корректность
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть положительным")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть положительной")
	}

	// Рассчитать среднюю скорость с помощью meanSpeed()
	meanSpeedValue := meanSpeed(steps, height, duration)

	// Рассчитать и вернуть количество калорий
	durationInMinutes := duration.Minutes()
	calories := (weight * meanSpeedValue * durationInMinutes) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверить входные параметры на корректность
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть положительным")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть положительной")
	}

	// Рассчитать среднюю скорость с помощью meanSpeed()
	meanSpeedValue := meanSpeed(steps, height, duration)

	// Рассчитать количество калорий
	durationInMinutes := duration.Minutes()
	calories := (weight * meanSpeedValue * durationInMinutes) / minInH

	// Возврат количества потраченных калорий
	calories *= walkingCaloriesCoefficient
	return calories, nil
}
