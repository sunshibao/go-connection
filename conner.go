package connection

import (
	"log"
)

var (
	connErrF = "连接外部程序出错%s"
)

type Connector interface {
	Connect() error
	Close() error
}

type ExternalProcedure struct {
	Connectors []Connector
}

func NewExternalProcedure(connectors ...Connector) *ExternalProcedure {
	e := &ExternalProcedure{}
	for _, connector := range connectors {
		err := connector.Connect()
		if err != nil {
			log.Printf(connErrF, err)
		}
		e.Connectors = append(e.Connectors, connector)
	}
	return e
}

func (e *ExternalProcedure) Close() {
	for _, connector := range e.Connectors {
		_ = connector.Close()
	}
}
