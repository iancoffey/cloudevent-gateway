package sources

import (
	"github.com/iancoffey/cloudevent-gateway/pkg/types"
	log "github.com/sirupsen/logrus"
)

func NewServer(events []string, receivers []types.EventReceiver, logger *log.Entry) *types.EventServer {
	return &types.EventServer{
		Events:    events,
		Receivers: receivers,
		Logger:    logger,
	}
}
