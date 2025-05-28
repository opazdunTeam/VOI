package service

import (
	"context"
	"encoding/json" // Импортируем для работы с JSON
	"errors"
	"fmt"
	"log"

	"github.com/jonathanhecl/gollama"
)

// Define a struct to unmarshal voiceDNA JSON into.
// Adjust fields based on the actual structure of your voiceDNA JSON.
type VoiceDNADetails struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	City    string `json:"city"`
	Hobbies []string `json:"hobbies"`
    // Добавьте сюда все поля, которые могут быть в вашем voiceDNA JSON
}


type GeneratorService interface {
	GenerateContent(ctx context.Context, prompt string, voiceDNA string) (string, error)
}

type generatorService struct {
	ollamaClient *gollama.Gollama
	modelName    string
}

func NewGeneratorService(modelName string) GeneratorService {
	client := gollama.New(modelName)
	// client.Verbose = true

	// Глобальный системный промпт (можно оставить, если нужен общий стиль ответов)
	// Этот промпт задает общую роль модели.
	defaultSystemPrompt := "Ты - опытный копирайтер, который генерирует тексты для социальных сетей. Учитывай предоставленную информацию о пользователе при создании поста. Не используй кавычки в начале или конце сгенерированного текста."
	client.SetSystemPrompt(defaultSystemPrompt)

	if err := client.PullIfMissing(context.Background()); err != nil {
		log.Printf("Внимание: Не удалось подтянуть/проверить модель Ollama '%s': %v. Убедитесь, что модель предварительно подтянута или Ollama доступна.", modelName, err)
	} else {
		log.Printf("Проверка/подтягивание модели Ollama '%s' успешно завершено.", modelName)
	}

	return &generatorService{
		ollamaClient: client,
		modelName:    modelName,
	}
}

func (s *generatorService) GenerateContent(ctx context.Context, prompt string, voiceDNA string) (string, error) {
	if prompt == "" {
		return "", errors.New("prompt не может быть пустым")
	}


	// --- НОВОЕ: Обработка voiceDNA как JSON-контекста ---
	var dnaDetails VoiceDNADetails
	if voiceDNA != "" {
		err := json.Unmarshal([]byte(voiceDNA), &dnaDetails)
		if err != nil {
			return "", fmt.Errorf("не удалось распарсить voiceDNA JSON: %w", err)
		}
	}

	// Формируем расширенный промпт, включая данные из voiceDNA
	// Здесь вы можете творчески подойти к тому, как представить эти данные модели.
	// Используйте форматирование, чтобы модель понимала, что это дополнительный контекст.
	fullPrompt := fmt.Sprintf("Основываясь на следующей информации о пользователе:\nИмя: %s\nВозраст: %d\nГород: %s\nУвлечения: %v\n\nНапиши пост на тему: %s",
		dnaDetails.Name,
		dnaDetails.Age,
		dnaDetails.City,
		dnaDetails.Hobbies,
		prompt, // Оригинальный запрос пользователя
	)

    // Если voiceDNA не пустой, то используем полный промпт, иначе - обычный.
    if voiceDNA == "" {
        fullPrompt = prompt
    }


	fmt.Printf("Генерируем контент с моделью '%s' для запроса: '%s'\n", s.modelName, fullPrompt)

	output, err := s.ollamaClient.Chat(ctx, fullPrompt) // Отправляем расширенный промпт
	if err != nil {
		return "", fmt.Errorf("не удалось сгенерировать контент из Ollama: %w", err)
	}

	return output.Content, nil
}