package whappdc

import "log"

type ChatWorker struct {
	incomingHandlers chan MessageHandler
}

func (w *ChatWorker) Start() {
	go func() {
		for {
			select {
			case handler := <-w.incomingHandlers:
				log.Println("Chat worker executing action")
				err := handler.Action()

				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
}
