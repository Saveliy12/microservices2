package path

import (
	"os"
	"path/filepath"
)

// GetProjectRoot ищет корневую директорию проекта по наличию go.work файла
// Нужен, чтобы из тестов указывать абсолютные пути к Dockerfile и миграциям
func GetProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic("не удалось получить рабочую директорию: " + err.Error())
	}

	for {
		_, err = os.Stat(filepath.Join(dir, "go.work"))
		if err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			panic("не удалось найти корень проекта (go.work)")
		}

		dir = parent
	}
}
