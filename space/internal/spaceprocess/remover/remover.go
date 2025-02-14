package remover

import (
	"context"

	"github.com/anyproto/any-sync/app"

	"github.com/anyproto/anytype-heart/space/clientspace"
	"github.com/anyproto/anytype-heart/space/internal/components/aclobjectmanager"
	"github.com/anyproto/anytype-heart/space/internal/components/builder"
	"github.com/anyproto/anytype-heart/space/internal/components/spaceloader"
	"github.com/anyproto/anytype-heart/space/internal/components/spacestatus"
	"github.com/anyproto/anytype-heart/space/internal/spaceprocess/loader"
	"github.com/anyproto/anytype-heart/space/internal/spaceprocess/mode"
)

type remover struct {
	app *app.App
}

type Remover interface {
	mode.Process
	loader.LoadWaiter
}

type Params struct {
	SpaceId             string
	Status              spacestatus.SpaceStatus
	StopIfMandatoryFail bool
	OwnerMetadata       []byte
}

func New(app *app.App, params Params) Remover {
	child := app.ChildApp()
	child.Register(params.Status).
		Register(builder.New()).
		Register(spaceloader.New(params.StopIfMandatoryFail, true)).
		Register(aclobjectmanager.New(params.OwnerMetadata))
	return &remover{
		app: child,
	}
}

func (r *remover) Start(ctx context.Context) error {
	return r.app.Start(ctx)
}

func (r *remover) Close(ctx context.Context) error {
	return r.app.Close(ctx)
}

func (r *remover) CanTransition(next mode.Mode) bool {
	return true
}

func (r *remover) WaitLoad(ctx context.Context) (sp clientspace.Space, err error) {
	spaceLoader := app.MustComponent[spaceloader.SpaceLoader](r.app)
	return spaceLoader.WaitLoad(ctx)
}
