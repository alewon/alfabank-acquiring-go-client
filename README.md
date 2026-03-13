# Go-клиент для сервиса «Альфабанк Эквайринг»

Типизированный Go-клиент для Alfa Bank Acquiring REST API.

- Репозиторий: [github.com/alewon/alfabank-acquiring-go-client](https://github.com/alewon/alfabank-acquiring-go-client)
- Документация пакета: [pkg.go.dev/github.com/alewon/alfabank-acquiring-go-client](https://pkg.go.dev/github.com/alewon/alfabank-acquiring-go-client)
- Релизы: [github.com/alewon/alfabank-acquiring-go-client/releases](https://github.com/alewon/alfabank-acquiring-go-client/releases)
- Лицензия: [MIT](https://github.com/alewon/alfabank-acquiring-go-client/blob/main/LICENSE)

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
