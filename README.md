# alfabank-acquiring-go-client

Типизированный Go-клиент для Alfa Bank Acquiring REST API.

Пакет ориентирован на:
- явные структуры запросов и ответов;
- прозрачные вложенные модели вместо сырых JSON-блоков;
- поддержку как `application/x-www-form-urlencoded`, так и JSON-endpoint-ов API.

## Статус

Пакет собран на основе локального снимка документации Alfa REST API, сохранённого в [doc.md](/home/user/Projects/clients/alfabank-acquiring-go-client/doc.md).

Поддерживаемая версия модуля:
- Go 1.18+

## Установка

```bash
go get github.com/alewon/alfabank-acquiring-go-client
```

## Быстрый старт

```go
package main

import (
	"context"
	"fmt"
	"net/http"

	alfabank "github.com/alewon/alfabank-acquiring-go-client"
)

func main() {
	client := alfabank.NewClient("api-user", "api-password", "", http.DefaultClient)

	result, err := client.Register(context.Background(), &alfabank.RegisterRequest{
		OrderNumber: "order-1001",
		Amount:      10000,
		Currency:    "643",
		ReturnURL:   "https://merchant.example/success",
		FailURL:     "https://merchant.example/fail",
		Description: "Test order",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(result.Response.OrderID)
	fmt.Println(result.Response.FormURL)
}
```

## Тестирование

```bash
go test ./...
```

## Документация

Основное покрытие API в репозитории включает:
- регистрацию и оплату заказов;
- связки;
- методы СБП;
- статусы заказов и чеков;
- методы шаблонов.

Локальный снимок документации намеренно хранится в репозитории для аудита и проверки моделей.

## Участие в разработке

См. [CONTRIBUTING.md](/home/user/Projects/clients/alfabank-acquiring-go-client/CONTRIBUTING.md).

## Безопасность

См. [SECURITY.md](/home/user/Projects/clients/alfabank-acquiring-go-client/SECURITY.md).

## Лицензия

[MIT](/home/user/Projects/clients/alfabank-acquiring-go-client/LICENSE)
