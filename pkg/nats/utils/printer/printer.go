package printer

import (
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
	"github.com/pseudoincorrect/bariot/pkg/utils/logger"
)

// printNatsMsg print a nats message
func PrintNatsMsg(m *nats.Msg) {
	str := fmt.Sprintf("NATS Message Received on [%s] Queue[%s] Pid[%d]", m.Subject, m.Sub.Queue, os.Getpid())
	logger.Info(str)
	str = fmt.Sprintf("NATS Message Payload %s", m.Data)
	logger.Info(str)
}
