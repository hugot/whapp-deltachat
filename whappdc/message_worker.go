package whappdc

import "log"

type MessageWorker struct {
	incomingHandlers chan MessageHandler
	chatWorkers      map[string]chan MessageHandler
	quit             chan bool
}

func NewMessageWorker() *MessageWorker {
	return &MessageWorker{
		incomingHandlers: make(chan MessageHandler),
		chatWorkers:      make(map[string]chan MessageHandler),
		quit:             make(chan bool),
	}
}

func (w *MessageWorker) HandleMessage(m MessageHandler) {
	w.incomingHandlers <- m
}

func (w *MessageWorker) Stop() {
	w.quit <- true
}

func (w *MessageWorker) Start() {
	go func() {
		for {
			select {
			case <-w.quit:
				return
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
