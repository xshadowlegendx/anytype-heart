package csv

import (
	"encoding/csv"
	"os"
	"path/filepath"

	"github.com/gogo/protobuf/types"
	"github.com/google/uuid"

	"github.com/anytypeio/go-anytype-middleware/core/block/collection"
	sb "github.com/anytypeio/go-anytype-middleware/core/block/editor/smartblock"
	"github.com/anytypeio/go-anytype-middleware/core/block/editor/state"
	"github.com/anytypeio/go-anytype-middleware/core/block/import/converter"
	"github.com/anytypeio/go-anytype-middleware/core/block/process"
	"github.com/anytypeio/go-anytype-middleware/pb"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/bundle"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/core/smartblock"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/pb/model"
	"github.com/anytypeio/go-anytype-middleware/util/pbtypes"
)

const (
	Name               = "Csv"
	rootCollectionName = "CSV Import"
)

type CSV struct {
	collectionService *collection.Service
}

func New(collectionService *collection.Service) converter.Converter {
	return &CSV{collectionService: collectionService}
}

func (c *CSV) Name() string {
	return Name
}

func (c *CSV) GetParams(req *pb.RpcObjectImportRequest) []string {
	if p := req.GetCsvParams(); p != nil {
		return p.Path
	}

	return nil
}

func (c *CSV) GetSnapshots(req *pb.RpcObjectImportRequest,
	progress *process.Progress) (*converter.Response, converter.ConvertError) {
	path := c.GetParams(req)
	if len(path) == 0 {
		return nil, nil
	}
	progress.SetProgressMessage("Start creating snapshots from files")
	snapshots := make([]*converter.Snapshot, 0)
	allRelations := make(map[string][]*converter.Relation, 0)
	allObjectsIDs := make([]string, 0)
	cErr := converter.NewError()
	for _, p := range path {
		if err := progress.TryStep(1); err != nil {
			cancelError := converter.NewFromError(p, err)
			return nil, cancelError
		}
		if filepath.Ext(p) != ".csv" {
			continue
		}
		csvTable, err := readCsvFile(p)
		if err != nil {
			cErr.Add(p, err)
			if req.Mode == pb.RpcObjectImportRequest_ALL_OR_NOTHING {
				return nil, cErr
			}
			continue
		}

		allObjectsIDs, snapshots, allRelations, err = c.handleCSVFile(p, csvTable, allObjectsIDs, snapshots, allRelations)
		if err != nil {
			cErr.Add(p, err)
			if req.Mode == pb.RpcObjectImportRequest_ALL_OR_NOTHING {
				return nil, cErr
			}
			continue
		}
	}

	rootCollection := converter.NewRootCollection(c.collectionService)
	rootCol, err := rootCollection.AddObjects(rootCollectionName, allObjectsIDs)
	if err != nil {
		cErr.Add(rootCollectionName, err)
		if req.Mode == pb.RpcObjectImportRequest_ALL_OR_NOTHING {
			return nil, cErr
		}
	}

	if rootCol != nil {
		snapshots = append(snapshots, rootCol)
	}
	progress.SetTotal(int64(len(allObjectsIDs)))
	if cErr.IsEmpty() {
		return &converter.Response{
			Snapshots: snapshots,
			Relations: allRelations,
		}, nil
	}

	return &converter.Response{
		Snapshots: snapshots,
		Relations: allRelations,
	}, cErr
}

func (c *CSV) handleCSVFile(path string,
	csvTable [][]string,
	allObjectsIDs []string,
	snapshots []*converter.Snapshot,
	allRelations map[string][]*converter.Relation) (
	[]string, []*converter.Snapshot, map[string][]*converter.Relation, error) {
	details := converter.GetDetails(path)
	details.GetFields()[bundle.RelationKeyLayout.String()] = pbtypes.Float64(float64(model.ObjectType_collection))
	_, _, st, err := c.collectionService.CreateCollection(details, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	relations := getDetailsFromCSVTable(csvTable)
	objectsSnapshots, objectsRelations := getEmptyObjects(csvTable, relations)
	targetIDs := make([]string, 0, len(objectsSnapshots))
	for _, objectsSnapshot := range objectsSnapshots {
		targetIDs = append(targetIDs, objectsSnapshot.Id)
	}
	allObjectsIDs = append(allObjectsIDs, targetIDs...)

	st.StoreSlice(sb.CollectionStoreKey, targetIDs)
	snapshot := c.getCollectionSnapshot(details, st, path)

	snapshots = append(snapshots, snapshot)
	snapshots = append(snapshots, objectsSnapshots...)
	allObjectsIDs = append(allObjectsIDs, snapshot.Id)

	allRelations[snapshot.Id] = relations
	allRelations = makeRelationsResultMap(allRelations, objectsRelations)
	return allObjectsIDs, snapshots, allRelations, nil
}

func (c *CSV) getCollectionSnapshot(details *types.Struct, st *state.State, p string) *converter.Snapshot {
	details = pbtypes.StructMerge(st.CombinedDetails(), details, false)
	sn := &model.SmartBlockSnapshotBase{
		Blocks:        st.Blocks(),
		Details:       details,
		ObjectTypes:   []string{bundle.TypeKeyCollection.URL()},
		Collections:   st.Store(),
		RelationLinks: st.GetRelationLinks(),
	}

	snapshot := &converter.Snapshot{
		Id:       uuid.New().String(),
		FileName: p,
		Snapshot: &pb.ChangeSnapshot{Data: sn},
		SbType:   smartblock.SmartBlockTypeCollection,
	}
	return snapshot
}

func readCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func getEmptyObjects(csvTable [][]string, relations []*converter.Relation) ([]*converter.Snapshot, map[string][]*converter.Relation) {
	snapshots := make([]*converter.Snapshot, 0, len(csvTable))
	objectsRelations := make(map[string][]*converter.Relation, len(csvTable))

	for i := 1; i < len(csvTable); i++ {
		details := &types.Struct{Fields: map[string]*types.Value{}}
		for j, value := range csvTable[i] {
			details.Fields[relations[j].Name] = pbtypes.String(value)
		}
		sn := &converter.Snapshot{
			Id:     uuid.New().String(),
			SbType: smartblock.SmartBlockTypePage,
			Snapshot: &pb.ChangeSnapshot{
				Data: &model.SmartBlockSnapshotBase{
					Details: details,
				},
			},
		}
		snapshots = append(snapshots, sn)

		objectsRelations[sn.Id] = relations
	}
	return snapshots, objectsRelations
}

func getDetailsFromCSVTable(csvTable [][]string) []*converter.Relation {
	if len(csvTable) == 0 {
		return nil
	}
	relations := make([]*converter.Relation, 0, len(csvTable[0]))
	for _, relation := range csvTable[0] {
		relations = append(relations, &converter.Relation{
			Relation: &model.Relation{
				Format: model.RelationFormat_longtext,
				Name:   relation,
			},
		})
	}
	return relations
}

func makeRelationsResultMap(rel1 map[string][]*converter.Relation, rel2 map[string][]*converter.Relation) map[string][]*converter.Relation {
	if len(rel1) != 0 {
		for id, relations := range rel2 {
			rel1[id] = relations
		}
		return rel1
	}
	if len(rel2) != 0 {
		for id, relations := range rel1 {
			rel2[id] = relations
		}
		return rel2
	}
	return map[string][]*converter.Relation{}
}