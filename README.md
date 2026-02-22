# goappconfig
go application configuration builder

### Usage

```go
import "github.com/actofgod/goappconfig"

type AppConfig struct {
	// ...
}

func main() {
	builder := goappconfig.NewBuilder[AppConfig]()
	err := builder.Load("config.json")
	if err != nil {
		panic(err)
	}
	config, err := builder.Build()
	if err != nil {
		panic(err)
	}
	// now config initialized
}
```
