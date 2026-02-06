package daysteps

import (
	"time"
	"strconv"
	"strings"
	"fmt"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// Разделить строку на слайс строк
	parts := strings.Split(data, ",")
	
	// Проверить, чтобы длина слайса была равна 2
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неправильный формат данных: ожидается 'шаги,время'")
	}

	stepsStr := strings.TrimSpace(parts[0])
	durationStr := strings.TrimSpace(parts[1])

	// Преобразовать первый элемент слайса в тип int
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования шагов: %v", err)
	}

	// Проверить: количество шагов должно быть больше 0
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше 0")
	}

	// Преобразовать второй элемент слайса в time.Duration
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования времени: %v", err)
	}

	if duration <= 0 {
		return 0, 0, fmt.Errorf("продолжительность должна быть положительной")
	}
	// Возврат количества шагов, продолжительности и nil
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// Получить данные о количестве шагов и продолжительности прогулки
	steps, duration, err := parsePackage(data)
	if err != nil {
		return ""
	}

	// Проверить, чтобы количество шагов было больше 0
	if steps <= 0 {
		return ""
	}

	// Вычислить дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// Перевести дистанцию в километры
	distanceKm := distanceMeters / mInKm

	// Вычислить количество калорий
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	
	if err != nil {
		return ""
	}

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.",
		steps, distanceKm, calories)
	
	return result
}
