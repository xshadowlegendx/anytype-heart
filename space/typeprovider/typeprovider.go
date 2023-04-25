package typeprovider

import (
	"context"
	"errors"
	"github.com/anytypeio/any-sync/app"
	"github.com/anytypeio/any-sync/commonspace/object/tree/treechangeproto"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/core/smartblock"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/logging"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/pb/model"
	"github.com/anytypeio/go-anytype-middleware/space"
	"github.com/gogo/protobuf/proto"
	"go.uber.org/zap"
	"sync"
)

const CName = "space.typeprovider"

var log = logging.Logger(CName)

var ErrUnknownSmartBlockType = errors.New("error unknown smartblock type")

type ObjectTypeProvider interface {
	app.Component
	Type(id string) (smartblock.SmartBlockType, error)
}

func New() ObjectTypeProvider {
	return &objectTypeProvider{}
}

type objectTypeProvider struct {
	sync.Mutex
	spaceService space.Service
	cache        map[string]smartblock.SmartBlockType
}

func (o *objectTypeProvider) Init(a *app.App) (err error) {
	o.spaceService = a.MustComponent(space.CName).(space.Service)
	o.cache = map[string]smartblock.SmartBlockType{}
	return
}

func (o *objectTypeProvider) Name() (name string) {
	return CName
}

func (o *objectTypeProvider) Type(id string) (tp smartblock.SmartBlockType, err error) {
	tp, err = smartblock.SmartBlockTypeFromID(id)
	if err != nil || tp != smartblock.SmartBlockTypePage {
		return
	}
	return o.objectTypeFromSpace(id)
}

func (o *objectTypeProvider) objectTypeFromSpace(id string) (tp smartblock.SmartBlockType, err error) {
	o.Lock()
	tp, exists := o.cache[id]
	if exists {
		o.Unlock()
		return
	}
	o.Unlock()

	sp, err := o.spaceService.AccountSpace(context.Background())
	if err != nil {
		return
	}

	store := sp.Storage()
	rawRoot, err := store.TreeRoot(id)
	if err != nil {
		return
	}

	root, err := o.unmarshallRoot(rawRoot)
	if err != nil {
		return
	}

	tp, err = o.objectType(root.ChangeType)
	if err != nil {
		return
	}
	o.Lock()
	defer o.Unlock()
	o.cache[id] = tp
	return
}

func (o *objectTypeProvider) objectType(changeType string) (smartblock.SmartBlockType, error) {
	log.With(zap.String("changeType", changeType)).Warn("getting change type")
	if v, exists := model.SmartBlockType_value[changeType]; exists {
		return smartblock.SmartBlockType(v), nil
	}

	return smartblock.SmartBlockTypePage, nil
}

func (o *objectTypeProvider) unmarshallRoot(rawRoot *treechangeproto.RawTreeChangeWithId) (root *treechangeproto.RootChange, err error) {
	raw := &treechangeproto.RawTreeChange{}
	err = proto.Unmarshal(rawRoot.GetRawChange(), raw)
	if err != nil {
		return
	}

	root = &treechangeproto.RootChange{}
	err = proto.Unmarshal(raw.Payload, root)
	if err != nil {
		return
	}
	return
}