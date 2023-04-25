package status

//
// import (
// 	"context"
// 	"testing"
// 	"time"
//
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/require"
//
// 	"github.com/anytypeio/go-anytype-middleware/core/anytype/config"
// 	"github.com/anytypeio/go-anytype-middleware/core/filestorage/filesync/mock_filesync"
// 	"github.com/anytypeio/go-anytype-middleware/pb"
// 	"github.com/anytypeio/go-anytype-middleware/pkg/lib/core/mock_core"
// 	"github.com/anytypeio/go-anytype-middleware/space/mock_space"
// 	"github.com/anytypeio/go-anytype-middleware/space/typeprovider/mock_typeprovider"
// )
//
// type fixture struct {
// 	ctrl          *gomock.Controller
// 	statusWatcher *mock_filesync.MockStatusWatcher
// 	statusService Service
// }
//
// func newFixture(t *testing.T) *fixture {
// 	fx := fixture{
// 		ctrl: gomock.NewController(t),
// 	}
//
// 	fileSync := mock_filesync.NewMockFileSync(fx.ctrl)
// 	statusWatcher := mock_filesync.NewMockStatusWatcher(fx.ctrl)
// 	eventReceiver := func(event *pb.Event) {
//
// 	}
//
// 	spaceService := mock_space.NewMockService(fx.ctrl)
// 	typeProvider := mock_typeprovider.NewMockSmartBlockTypeProvider(fx.ctrl)
// 	coreService := mock_core.NewMockService(fx.ctrl)
//
// 	statusService := New(typeProvider, &config.Config{}, eventReceiver, spaceService, coreService, fileSync)
//
// 	ctx := context.Background()
// 	fileSync.EXPECT().NewStatusWatcher(statusService, 5*time.Second).Return(statusWatcher)
//
// 	spc := mock_space.NewMockSpace(fx.ctrl)
//
// 	spaceService.EXPECT().AccountSpace(ctx).Return(spc, nil)
//
// 	err := statusService.Run(ctx)
// 	require.NoError(t, err)
// 	return &fx
// }
//
// func (f *fixture) finish() {
// 	f.ctrl.Finish()
// }
//
// func TestUpdateTree(t *testing.T) {
// 	fx := newFixture(t)
// 	defer fx.finish()
//
// 	fx.statusService.Watch("kek", nil)
// }