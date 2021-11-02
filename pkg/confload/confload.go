package confload

import (
	"log"
	"os"

	"github.com/max-gui/logagent/pkg/logsets"
)

func Load() []byte {
	bytes, err := os.ReadFile(*logsets.Apppath + string(os.PathSeparator) + "application-" + *logsets.Appenv + ".yml")
	if err != nil {
		log.Panic(err)
	}

	return bytes
}
