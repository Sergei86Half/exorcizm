package main

import (
	"bufio"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"log"
	"os"
	"strings"
	"time"
)

const (
	sampleRate = 44100 // Стандартная частота дискретизации для аудио
	// Пути к аудиофайлам (убедитесь, что они находятся в той же директории, что и исполняемый файл)
	agonySoundFile = "agony.wav"
	roarSoundFile  = "roar.wav"
)

var (
	audioContext *audio.Context
	successSequence = []string{"свечи", "кресты", "свечи"} // Правильная последовательность предметов
)

// playSound воспроизводит указанный WAV-файл.
func playSound(filePath string) {
	if audioContext == nil {
		log.Printf("Аудио контекст не инициализирован, не могу воспроизвести %s", filePath)
		return
	}

	f, err := os.Open(filePath)
	if err != nil {
		log.Printf("Не удалось открыть аудио файл %s: %v. Убедитесь, что файл существует.", filePath, err)
		return
	}
	defer f.Close()

	// Декодируем WAV-файл
	d, err := wav.DecodeWithSampleRate(sampleRate, f)
	if err != nil {
		log.Printf("Не удалось декодировать WAV файл %s: %v. Убедитесь, что это корректный WAV.", filePath, err)
		return
	}

	// Создаем новый проигрыватель и воспроизводим звук
	player, err := audioContext.NewPlayer(d)
	if err != nil {
		log.Printf("Не удалось создать аудио проигрыватель для %s: %v", filePath, err)
		return
	}

	player.Play()
	log.Printf("Воспроизведение %s...", filePath)
	// Ждем, пока звук закончится
	for player.IsPlaying() {
		time.Sleep(time.Millisecond * 100)
	}
	player.Close() // Закрываем проигрыватель после воспроизведения
}

// displayScene отрисовывает текущую сцену в консоли.
func displayScene(pentagramActive, entityAggressive bool, itemsUsed []string) {
	fmt.Print("\033[H\033[2J") // Очистка консоли (работает в большинстве терминалов)

	fmt.Println("========================================")
	fmt.Println("           РИТУАЛ ЭКЗОРЦИЗМА            ")
	fmt.Println("========================================")

	if pentagramActive {
		fmt.Println("    /\\")
		fmt.Println("   /  \\")
		fmt.Println("  / -- \\    * Активная пентаграмма *")
		fmt.Println(" <------->")
		fmt.Println("  \\ -- /")
		fmt.Println("   \\  /")
		fmt.Println("    \\/")
	} else {
		fmt.Println("    /\\")
		fmt.Println("   /  \\")
		fmt.Println("  /    \\    * Пентаграмма исчезла *")
		fmt.Println(" <      >")
		fmt.Println("  \\    /")
		fmt.Println("   \\  /")
		fmt.Println("    \\/")
	}

	fmt.Println("\n----------------------------------------")
	fmt.Println(" Статус Сущности:")
	if entityAggressive {
		fmt.Println("   --> СТАНОВИТСЯ АГРЕССИВНОЙ! <--")
	} else if pentagramActive {
		fmt.Println("   --> Прикована к пентаграмме.")
	} else {
		fmt.Println("   --> Изгнана! Покойся с миром...")
	}
	fmt.Println("----------------------------------------")

	fmt.Printf(" Используемые предметы: [%s]\n", strings.Join(itemsUsed, ", "))
	fmt.Println("----------------------------------------")
}

// getPlayerInput считывает ввод пользователя.
func getPlayerInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input))
}

func main() {
	log.Println("Запуск ритуала изгнания...")

	// Инициализация аудио контекста
	audioContext = audio.NewContext(sampleRate)
	if audioContext == nil {
		log.Fatal("Не удалось инициализировать аудио контекст. Проверьте аудио-драйверы.")
	}

	var (
		pentagramActive  = true
		entityAggressive = false
		itemsUsed        []string
		attempts         = 0
		maxAttempts      = len(successSequence) + 2 // Даем немного больше попыток, чем нужно для успеха
	)

	// Основной игровой цикл
	for pentagramActive && !entityAggressive && attempts < maxAttempts {
		displayScene(pentagramActive, entityAggressive, itemsUsed)

		fmt.Println("\nЧто используем? ( свечи / кресты)")
		fmt.Print("> ")
		input := getPlayerInput()

		attempts++

		// Проверяем ввод
		if input != "свечи" && input != "кресты" {
			fmt.Println("Неизвестный предмет. Попробуйте 'свечи' или 'кресты'.")
			time.Sleep(time.Second * 2)
			continue // Не считать это за попытку использования предмета
		}

		itemsUsed = append(itemsUsed, input)

		// Проверяем последовательность
		correctSoFar := true
		for i, item := range itemsUsed {
			if i >= len(successSequence) || item != successSequence[i] {
				correctSoFar = false
				break
			}
		}

		if !correctSoFar {
			entityAggressive = true
			fmt.Println("\nО нет! Последовательность нарушена! Сущность становится агрессивной!")
			playSound(roarSoundFile)
			time.Sleep(time.Second * 3)
			break // Выход из цикла, так как сущность агрессивна
		}

		if len(itemsUsed) == len(successSequence) && correctSoFar {
			pentagramActive = false
			fmt.Println("\nДА! Правильная последовательность! Пентаграмма исчезает, сущность издает крик агонии!")
			playSound(agonySoundFile)
			time.Sleep(time.Second * 3)
			break // Успешное изгнание, выход из цикла
		}

		fmt.Println("Вы использовали: " + input + ". Что дальше?")
		time.Sleep(time.Second * 2) // Небольшая пауза для чтения
	}

	// Конец игры
	displayScene(pentagramActive, entityAggressive, itemsUsed) // Обновить финальную сцену
	fmt.Println("\n========================================")
	if !pentagramActive && !entityAggressive {
		fmt.Println("      РИТУАЛ ЗАВЕРШЕН УСПЕШНО!         ")
		fmt.Println("        Сущность изгнана.             ")
	} else if entityAggressive {
		fmt.Println("      РИТУАЛ ПРОВАЛЕН!                 ")
		fmt.Println("        Сущность свободна и зла.      ")
	} else {
		fmt.Println("      РИТУАЛ ПРЕРВАН.                 ")
		fmt.Println("        Недостаточно попыток.         ")
	}
	fmt.Println("========================================")
	fmt.Println("Игра завершена.")
}
