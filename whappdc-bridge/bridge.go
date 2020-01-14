package bridge

import (
	"time"

	"github.com/hugot/go-deltachat/deltabot"
	"github.com/hugot/whapp-deltachat/botcommands"
	"github.com/hugot/whapp-deltachat/whappdc"
	core "github.com/hugot/whapp-deltachat/whappdc-core"
)

type Bridge struct {
	Ctx *core.BridgeContext
}

func (b *Bridge) Init(config *core.Config) error {
	messageWorker := whappdc.NewMessageWorker()
	ctx := core.NewBridgeContext(config, 15*time.Minute)

	err := ctx.Init(
		whappdc.NewWhappHandler(ctx, messageWorker),
		[]deltabot.Command{
			&botcommands.Echo{},
			botcommands.NewWhappBridge(ctx),
		},
	)

	if err != nil {
		return err
	}

	messageWorker.Start()

	b.Ctx = ctx

	return nil
}

func (b *Bridge) Close() error {
	return b.Ctx.Close()
}
