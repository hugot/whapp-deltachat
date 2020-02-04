package whappdc

import (
	"github.com/hugot/go-deltachat/deltachat"
)

// ChatWorker receives structs of type MessageHandler and executes them sequentially. By
// executing the handlers sequentially we try to make sure that the messages are sent
// through deltachat in the right order.
type ChatWorker struct {
	logger           deltachat.Logger
	incomingHandlers chan MessageHandler
	quit             chan bool
}

func NewChatWorker(logger deltachat.Logger) *ChatWorker {
	return &ChatWorker{
		incomingHandlers: make(chan MessageHandler, 3),
		quit:             make(chan bool),
		logger:           logger,
	}
}

func (w *ChatWorker) HandleMessage(m MessageHandler) {
	w.incomingHandlers <- m
}

func (w *ChatWorker) Stop() {
	w.quit <- true
}

func (w *ChatWorker) Start() {
	go func() {
		for {
			select {
			case <-w.quit:
				return
			case handler := <-w.incomingHandlers:
				err := handler.Action()

				if err != nil {
					w.logger.Println(err)
				}
			}
		}
	}()
}
