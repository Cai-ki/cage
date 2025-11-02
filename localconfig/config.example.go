package localconfig

import "os"

var (
	AnyAPIToken = "your-api-token"
)

func init() {
	os.Setenv("your-key", "your-value")
}
