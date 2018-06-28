package cdn

import (
	// "log"
	// "fmt"
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("cdn")
var format = logging.MustStringFormatter(
	`%{color}[-%{program}-][%{module}]-%{time:15:04:05} %{shortfile} %{shortfunc} >>> %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func init() {

	// log := logging.MustGetLogger(ns)
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")

	logging.SetBackend(backend2Formatter)
	// log.S
}
