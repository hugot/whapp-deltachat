package main

import "log"

type MessageWorker struct {
	incomingHandlers chan MessageHandler
	chatWorkers      map[string]chan MessageHandler
}

func NewMessageWorker() *MessageWorker {
	return &MessageWorker{
		incomingHandlers: make(chan MessageHandler),
		chatWorkers:      make(map[string]chan MessageHandler),
	}
}

func (w *MessageWorker) HandleMessage(m MessageHandler) {
	w.incomingHandlers <- m
}

func (w *MessageWorker) Start() {
	go func() {
		for {
			select {
			case handler := <-w.incomingHandlers:
				log.Println("Got Handler for " + handler.Jid)
				workerChan, ok := w.chatWorkers[handler.Jid]

				if !ok {
					workerChan = make(chan MessageHandler)

					worker := &ChatWorker{
						incomingHandlers: workerChan,
					}

					worker.Start()
					w.chatWorkers[handler.Jid] = workerChan
				}

				workerChan <- handler
			}
		}
	}()
}
