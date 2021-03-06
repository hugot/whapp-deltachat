package whappdc

import (
	"sync"
	"time"

	"github.com/hugot/go-deltachat/deltachat"
	"github.com/hugot/whapp-deltachat/core"
)

func NewMessageTracker(
	DB *core.Database,
	flushInterval time.Duration,
	logger deltachat.Logger,
) *MessageTracker {
	tracker := &MessageTracker{
		DB:     DB,
		logger: logger,
	}

	tracker.FlushWithInterval(flushInterval)

	return tracker
}

// MessageTracker will keep track of encountered whatsapp messages to prevent sending them
// twice. It's storage is buffered to prevent continuous locks on the database. This means
// that calling WasSent immediately after calling MarkSent will most likely not return an
// up to date answer.
type MessageTracker struct {
	DB             *core.Database
	logger         deltachat.Logger
	delivered      [80]*string
	deliveredMutex sync.RWMutex
	deliveredIdx   int
}

func (t *MessageTracker) MarkSent(ID *string) error {
	t.deliveredMutex.Lock()
	defer t.deliveredMutex.Unlock()

	t.delivered[t.deliveredIdx] = ID

	if t.deliveredIdx == len(t.delivered)-1 {
		err := t.flush()

		if err != nil {
			return err
		}
	}

	t.deliveredIdx += 1

	return nil
}

// Flush without lock
func (t *MessageTracker) flush() error {
	err := t.DB.MarkWhappMessagesSent(t.delivered[:])
	t.deliveredIdx = 0

	return err
}

// Flush with lock
func (t *MessageTracker) Flush() error {
	t.deliveredMutex.Lock()
	defer t.deliveredMutex.Unlock()

	return t.flush()
}

func (t *MessageTracker) WasSent(ID string) (bool, error) {
	return t.DB.WhappMessageWasSent(ID)
}

func (t *MessageTracker) FlushWithInterval(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			err := t.Flush()

			if err != nil {
				t.logger.Println(err)
			}
		}
	}()
}
