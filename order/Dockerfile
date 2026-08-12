# ======= Этап 1: Сборка приложения =======

# Используем официальный образ Go с конкретной версией для стабильности сборки
FROM golang:1.25.8 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod/go.sum модулей для кеширования зависимостей.
# go.work в контейнер НЕ копируем: модуль собирается сам,
# а связь с shared обеспечивает replace (order/go.mod: replace shared => ../shared)
COPY order/go.mod order/go.sum ./order/
COPY shared/go.mod shared/go.sum ./shared/

# Загружаем зависимости заранее
RUN cd order && go mod download

# Копируем исходный код модуля и shared
COPY order/ ./order/
COPY shared/ ./shared/

# ⚠️ ВАЖНО: Собираем статически слинкованный бинарник для Alpine
RUN cd order && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server ./cmd

# ======= Этап 2: Минимальный образ для запуска приложения =======

# Используем минимальный стабильный образ Alpine Linux для финального контейнера
FROM alpine:3.21.3

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Создаём непривилегированного пользователя и группу с фиксированными UID/GID
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Копируем готовый бинарник из стадии сборки
COPY --from=builder /server .

# ✅ Копируем папку с миграциями из стадии сборки в финальный образ
COPY --from=builder /app/order/migrations ./migrations

# Меняем владельца файлов на созданного пользователя
RUN chown -R appuser:appgroup /app

# Переключаемся на непривилегированного пользователя
USER appuser

# Order слушает HTTP на порту 8080
# ⚠️ EXPOSE — это лишь мета-информация для документации образа.
# Чтобы порт стал доступен снаружи, при запуске нужно явно указать -p
EXPOSE 8080

# Устанавливаем команду по умолчанию для запуска контейнера
CMD ["./server"]
