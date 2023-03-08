package service

import (
	"github.com/twinj/uuid"
	"log"
	"math"
	"sync"
	"time"
)

var mut sync.RWMutex

type inMemoryMessageRepository struct {
	messages     []StoredMessage
	messagesById map[string]*StoredMessage
	*sync.RWMutex
}

func (imr *inMemoryMessageRepository) GetMessage(id string) (StoredMessage, error) {
	imr.RLock()
	defer imr.RUnlock()

	return *imr.messagesById[id], nil
}

func (imr *inMemoryMessageRepository) AddMessage(message Message) (StoredMessage, error) {
	imr.Lock()
	defer imr.Unlock()

	id := uuid.NewV4().String()

	msg := StoredMessage{Id: id, Message: message, CreatedAt: time.Now().UTC()}
	imr.messages = append(imr.messages, msg)
	imr.messagesById[id] = &msg

	return msg, nil
}

func distance(loc1 Location, loc2 Location) float64 {
	lat1, lon1, lat2, lon2 := loc1.Lat, loc1.Long, loc2.Lat, loc2.Long
	const earthRadius = 6371000 // Earth's radius in meters
	phi1 := toRadians(lat1)
	phi2 := toRadians(lat2)
	deltaPhi := toRadians(lat2 - lat1)
	deltaLambda := toRadians(lon2 - lon1)

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	log.Printf("distance: %f", earthRadius*c)
	return earthRadius * c
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func (imr *inMemoryMessageRepository) GetMessagesForLocation(location Location, radiusMeters float64, limit int,
	after time.Time) ([]StoredMessage, error) {
	imr.RLock()
	defer imr.RUnlock()

	messages := make([]StoredMessage, 0)
	for i := len(imr.messages) - 1; i >= 0 && len(messages) < limit; i-- {
		msg := imr.messages[i]

		if !msg.CreatedAt.After(after) {
			break
		}

		if distance(location, msg.Location) < radiusMeters {
			messages = append(messages, msg)
		}
	}

	return messages, nil
}

func MakeInMemoryRepository(config Configuration) (MessageRepository, error) {
	return &inMemoryMessageRepository{
		make([]StoredMessage, 0),
		make(map[string]*StoredMessage),
		&mut,
	}, nil
}
