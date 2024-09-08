# Используем более новый образ Go (версия 1.21)
FROM golang:1.21

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем исходные файлы приложения
COPY ai_service.proto .
COPY main.go .
COPY go.mod .
RUN go mod download
# Установка protoc (Protobuf компилятор)
RUN apt-get update && apt-get install -y protobuf-compiler
COPY . ./
# Установка плагинов для генерации Go-классов из .proto файлов
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Удалите старую директорию, если существует, и создайте новую
RUN rm -rf /app/generated \
    && mkdir -p /app/generated

# Убедись, что файл ai-service.proto действительно существует
RUN ls -la /app

# Генерация Go-классов для обоих .proto файлов
RUN protoc --go_out=paths=source_relative:./generated --go-grpc_out=paths=source_relative:./generated ./ai_service.proto
    #\
    #&& protoc --go_out=paths=source_relative:./generated --go-grpc_out=paths=source_relative:./generated ./ai_service_interaction.proto




# Проверяем сгенерированные файлы
RUN echo "Generated files:" && ls -la /app/generated

# Устанавливаем зависимости Go
RUN go mod tidy

# Собираем приложение
#RUN go build -o ai_service ./main.go && echo "Сборка успешна" || (echo "Ошибка сборки!" && exit 1)
RUN CGO_ENABLED=0 GOOS=linux go build -o /ai_service
# Убедись, что файл ai-service.proto действительно существует
RUN ls -la /app

# Открываем порт для gRPC
EXPOSE 50051
WORKDIR /app
RUN chmod +x /ai_service
# Запускаем gRPC сервер
CMD ["/ai_service"]
