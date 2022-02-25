# Overview
Generate generate number by predefined value, which has the same trend

## Architecture
[README](https://blog.xstudio.mobi/a/203.html)

# Usage
```go
import (
	"fmt"

	generator "github.com/xstudio/pre_generator"
)

func main() {
	var pre int64 = 123
	id := generator.New().Generate(pre)
	fmt.Printf("ID stirng  ID: %s\n", id)

	parsedPre, err := generator.New().ParseString(id.String())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("parsed Pre is: %d\n", parsedPre)
}
```


# Team
* @xstudio
