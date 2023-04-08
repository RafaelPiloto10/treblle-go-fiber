# treblle-go-fiber

The offical Trebble SDK for Go using Fiber

## Get Started

1. Sign in to [Treblle](https://app.treblle.com).
2. [Create a Treblle project](https://docs.treblle.com/en/dashboard/projects#creating-a-project).
3. [Setup the SDK](#install-the-SDK) for your platform.

### Install the SDK

```sh
go get github.com/treblle/treblle-go
```

Configure Treblle at the start of your `main()` function:

```go
import "github.com/treblle/treblle-go"

func main() {
	treblle.Configure(treblle.Configuration{
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
// TODO: Add Fiber code integration example
```

> See the [docs](https://docs.treblle.com/en/integrations/go) for this SDK to learn more.
