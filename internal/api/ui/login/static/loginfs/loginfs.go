package loginfs

import (
	"fmt"
	"net/http"

	"github.com/rakyll/statik/fs"
)

var loadedloginFS http.FileSystem

func MustLoad() http.FileSystem {
	if loadedloginFS != nil {
		return loadedloginFS
	}
	statikLoginFS, err := fs.NewWithNamespace("login")
	if err != nil {
		panic(fmt.Errorf("unable to start login statik dir: %w", err))
	}
	loadedloginFS = statikLoginFS
	return loadedloginFS
}
