package whappdc

// MessageWorker receives structs of type MessageHandler and distributes them across chat
// workers. Each whatsapp chat should have its own worker. If a message is encountered for
// a chat that has no worker yet, it is created.
type MessageWorker struct {
	incomingHandlers chan MessageHandler
	chatWorkers      map[string]*ChatWorker
	quit             chan bool
}

func NewMessageWorker() *MessageWorker {
	return &MessageWorker{
		incomingHandlers: make(chan MessageHandler),
		chatWorkers:      make(map[string]*ChatWorker),
		quit:             make(chan bool),
	}
}

func (w *MessageWorker) HandleMessage(m MessageHandler) {
	w.incomingHandlers <- m
}

func (w *MessageWorker) Stop() {
	w.quit <- true

	for _, worker := range w.chatWorkers {
		worker.Stop()
	}
}

func (w *MessageWorker) Start() {
	go func() {
		for {
			select {
			case <-w.quit:
				return
			case handler := <-w.incomingHandlers:
				worker, ok := w.chatWorkers[handler.Jid]

				if !ok {
					worker = NewChatWorker()
					worker.Start()

					w.chatWorkers[handler.Jid] = worker
				}

				worker.HandleMessage(handler)
			}
		}
	}()
}
