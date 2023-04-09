# treblle-go-fiber

The offical Trebble SDK for Go using Fiber

## Get Started

1. Sign in to [Treblle](https://app.treblle.com).
2. [Create a Treblle project](https://docs.treblle.com/en/dashboard/projects#creating-a-project).

### Install the Go-Fiber SDK

```sh
go get github.com/RafaelPiloto10/treblle-go-fiber/trebble_fiber@latest
```

Configure Treblle at the start of your `main()` function:

```go
package main

import (
	treblle_fiber "github.com/RafaelPiloto10/treblle-go-fiber/trebble_fiber"
)

func main() {
	treblle_fiber.Configure(treblle_fiber.Configuration{
		APIKey:     "YOUR API KEY HERE",
		ProjectID:  "YOUR PROJECT ID HERE",
		KeysToMask: []string{"password", "card_number"}, // optional, mask fields you don't want sent to Treblle
		ServerURL:  "https://rocknrolla.treblle.com",    // optional, don't use default server URL
	}

    // rest of your program.
}
```

After that, just use the middleware with any of your handlers:

```go
app := fiber.New(fiber.Config{})
app.Use(treblle_fiber.Middleware())
app.Listen("localhost:3000")
```

> See the [docs](https://docs.treblle.com/en/integrations/go) to learn more.
