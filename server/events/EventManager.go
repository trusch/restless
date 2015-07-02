package events

import (
	"math/rand"
	"path/filepath"
	"sync"
)

type EventHandler func(topic string, payload interface{})

type eventHandlerContainer struct {
	Id      int64
	handler EventHandler
}

type EventManager struct {
	handlers map[string][]*eventHandlerContainer
	mutex    sync.Mutex
}

func NewEventManager() *EventManager {
	mgr := new(EventManager)
	mgr.handlers = make(map[string][]*eventHandlerContainer)
	return mgr
}

func (mgr *EventManager) Emit(topic string, payload interface{}) {
	handlers := make([]*eventHandlerContainer, 0)
	mgr.mutex.Lock()
	for key, handlerSlice := range mgr.handlers {
		if matched, _ := filepath.Match(key, topic); matched {
			handlers = append(handlers, handlerSlice...)
		}
	}
	mgr.mutex.Unlock()
	for _, handler := range handlers {
		go handler.handler(topic, payload)
	}
}

func (mgr *EventManager) On(topic string, handler EventHandler) int64 {
	id := rand.Int63()
	mgr.mutex.Lock()
	mgr.handlers[topic] = append(mgr.handlers[topic], &eventHandlerContainer{id, handler})
	mgr.mutex.Unlock()
	return id
}

func (mgr *EventManager) Off(id int64) bool {
	mgr.mutex.Lock()
	for topic, containers := range mgr.handlers {
		for idx, val := range containers {
			if val.Id == id {
				mgr.handlers[topic] = append(containers[:idx], containers[idx+1:]...)
				return true
			}
		}
	}
	mgr.mutex.Unlock()
	return false
}

func (mgr *EventManager) Many(topic string, num int, handler EventHandler) {
	id := int64(0)
	counter := 0
	id = mgr.On(topic, func(topic string, payload interface{}) {
		counter++
		if counter > num {
			mgr.Off(id)
			return
		}
		handler(topic, payload)
	})
}

func (mgr *EventManager) Once(topic string, handler EventHandler) {
	mgr.Many(topic, 1, handler)
}
