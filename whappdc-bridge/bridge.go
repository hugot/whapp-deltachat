package bridge

import (
	"time"

	"github.com/hugot/go-deltachat/deltabot"
	"github.com/hugot/whapp-deltachat/botcommands"
	"github.com/hugot/whapp-deltachat/whappdc"
	core "github.com/hugot/whapp-deltachat/whappdc-core"
)

type Bridge struct {
	Ctx          *core.BridgeContext
	whappHandler *whappdc.WhappHandler
}

func (b *Bridge) Init(config *core.Config) error {
	ctx := core.NewBridgeContext(config)
	whappHandler := whappdc.NewWhappHandler(ctx, 15*time.Minute)

	err := ctx.Init(
		whappHandler,
		[]deltabot.Command{
			&botcommands.Echo{},
			botcommands.NewWhappBridge(ctx),
		},
	)

	if err != nil {
		return err
	}

	whappHandler.Start()

	b.whappHandler = whappHandler
	b.Ctx = ctx

	return nil
}

func (b *Bridge) Close() error {
	err := b.whappHandler.Stop()

	if err != nil {
		return err
	}

	return b.Ctx.Close()
}
