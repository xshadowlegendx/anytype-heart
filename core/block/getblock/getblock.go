package getblock

import (
	"context"
	"fmt"

	"github.com/anytypeio/go-anytype-middleware/core/block/editor/smartblock"
	"github.com/anytypeio/go-anytype-middleware/metrics"
)

type Picker interface {
	PickBlock(ctx context.Context, id string) (sb smartblock.SmartBlock, err error)
}

func Do[t any](p Picker, id string, apply func(sb t) error) error {
	sb, err := p.PickBlock(context.WithValue(context.TODO(), metrics.CtxKeyRequest, "do"), id)
	if err != nil {
		return err
	}

	bb, ok := sb.(t)
	if !ok {
		var dummy = new(t)
		return fmt.Errorf("the interface %T is not implemented in %T", dummy, sb)
	}

	sb.Lock()
	defer sb.Unlock()
	return apply(bb)
}