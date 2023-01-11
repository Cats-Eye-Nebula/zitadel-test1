package eventstore_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type benchmarkTest struct {
	name    string
	cmds    []eventstore.Command
	workers int
}

var (
	instanceID1 = "1"
	instanceID2 = "2"

	ctx = authz.SetCtxData(context.Background(), authz.CtxData{UserID: "adlerhurst", OrgID: "myorg"})

	ctxInstance1 = authz.WithInstanceID(ctx, instanceID1)
	ctxInstance2 = authz.WithInstanceID(ctx, instanceID2)

	agg1 = eventstore.NewAggregate(ctx, "ng5PD", "test", "v1")
	agg2 = eventstore.NewAggregate(ctx, "e4epE", "test", "v1")

	testCases = []benchmarkTest{
		{
			name: "without",
			cmds: []eventstore.Command{
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
			},
			workers: 1,
		},
		{
			name: "with",
			cmds: []eventstore.Command{
				cmdWithPayload(ctx, agg1),
				cmdWithPayload(ctx, agg2),
			},
			workers: 1,
		},
		{
			name: "without",
			cmds: []eventstore.Command{
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
			},
			workers: 3,
		},
		{
			name: "with",
			cmds: []eventstore.Command{
				cmdWithPayload(ctx, agg1),
				cmdWithPayload(ctx, agg2),
			},
			workers: 3,
		},
		{
			name: "without",
			cmds: []eventstore.Command{
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
			},
			workers: 5,
		},
		{
			name: "with",
			cmds: []eventstore.Command{
				cmdWithPayload(ctx, agg1),
				cmdWithPayload(ctx, agg2),
			},
			workers: 5,
		},
		{
			name: "without",
			cmds: []eventstore.Command{
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
			},
			workers: 7,
		},
		{
			name: "with",
			cmds: []eventstore.Command{
				cmdWithPayload(ctx, agg1),
				cmdWithPayload(ctx, agg2),
			},
			workers: 7,
		},
		{
			name: "without",
			cmds: []eventstore.Command{
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
			},
			workers: 9,
		},
		{
			name: "with",
			cmds: []eventstore.Command{
				cmdWithPayload(ctx, agg1),
				cmdWithPayload(ctx, agg2),
			},
			workers: 9,
		},
	}
)

func Benchmark_Eventstore_Push(b *testing.B) {
	for _, tt := range testCases {
		execTest(b, localClient, tt.workers, fmt.Sprintf("%d-workers_%s-payload", tt.workers, tt.name), tt.cmds, false)
	}
}

func Benchmark_Eventstore_Push_2_instances(b *testing.B) {
	for _, tt := range testCases {
		execTest(b, localClient, tt.workers, fmt.Sprintf("%d-workers_%s-payload", tt.workers, tt.name), tt.cmds, true)
	}
}

func execTest(b *testing.B, client *sql.DB, workers int, name string, commands []eventstore.Command, twoInstances bool) {
	b.Helper()

	// warmup sql connections
	var warmupWg sync.WaitGroup
	warmupWg.Add(40)
	for i := 0; i < 40; i++ {
		go func() {
			client.Ping()
			warmupWg.Done()
		}()
	}
	warmupWg.Wait()

	es, err := eventstore.Start(&eventstore.Config{Client: client})
	if err != nil {
		b.Fatal("unable to init eventstore: ", err)
	}
	if _, err = localClient.Exec("TRUNCATE eventstore.events;"); err != nil {
		b.Fatal("unable to truncate table: ", err)
	}

	if _, err = localClient.Exec("CREATE SEQUENCE IF NOT EXISTS eventstore.i_1_seq"); err != nil {
		b.Fatal("unable to create instance 1: ", err)
	}

	if _, err = localClient.Exec("CREATE SEQUENCE IF NOT EXISTS eventstore.i_2_seq"); err != nil {
		b.Fatal("unable to create instance 2: ", err)
	}

	b.Run(name, func(b *testing.B) {
		var wg sync.WaitGroup
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func(i int) {
				pushCtx := ctxInstance1
				if twoInstances && i%2 == 0 {
					pushCtx = ctxInstance2
				}
				for n := 0; n < b.N; n++ {
					if _, err := es.Push(pushCtx, commands...); err != nil {
						b.Error("push failed: ", err)
					}
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
	})
}

type benchCommand struct {
	eventstore.BaseEvent
	payload []byte
}

func commandWithoutPayload(ctx context.Context, agg *eventstore.Aggregate) *benchCommand {
	typ := eventstore.EventType("test")
	return &benchCommand{
		BaseEvent: *eventstore.NewBaseEventForPush(ctx, agg, typ),
	}
}

func cmdWithPayload(ctx context.Context, agg *eventstore.Aggregate) *benchCommand {
	cmd := commandWithoutPayload(ctx, agg)
	var err error

	cmd.payload, err = json.Marshal(struct {
		Username    string
		Firstname   string
		Lastname    string
		Email       string
		DisplayName string
		Gender      int8
	}{
		// The IT crowd S2.E4
		Username:    "peterfile",
		Firstname:   "Peter",
		Lastname:    "File",
		Email:       "peter.file@somemail.com",
		DisplayName: "Peter File",
		Gender:      10,
	})
	if err != nil {
		logging.Fatal("unable to create payload: ", err)
	}

	return cmd
}

func (cmd *benchCommand) Data() interface{} {
	if len(cmd.payload) == 0 {
		return nil
	}
	return cmd.payload
}

func (cmd *benchCommand) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

// func Benchmark_Eventstore_PushOneAggregate(b *testing.B) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	ctx = authz.SetCtxData(ctx, authz.CtxData{UserID: "adlerhurst", OrgID: "myorg"})
// 	defer cancel()

// 	agg := eventstore.NewAggregate(ctx, "S7boD", "test", "v1")

// 	tests := []struct {
// 		name    string
// 		client  *sql.DB
// 		cmds    []eventstore.Command
// 		workers int
// 	}{
// 		{
// 			name:   "1 event - no payload - 1 worker",
// 			client: localClient,
// 			cmds: []eventstore.Command{
// 				commandWithoutPayload(ctx, agg),
// 			},
// 			workers: 1,
// 		},
// 		{
// 			name:   "2 events - payload - 1 worker",
// 			client: localClient,
// 			cmds: []eventstore.Command{
// 				cmdWithPayload(ctx, agg),
// 				cmdWithPayload(ctx, agg),
// 			},
// 			workers: 1,
// 		},
// 		{
// 			name:   "5 events - no payload - 1 worker",
// 			client: localClient,
// 			cmds: []eventstore.Command{
// 				commandWithoutPayload(ctx, agg),
// 				commandWithoutPayload(ctx, agg),
// 				commandWithoutPayload(ctx, agg),
// 				commandWithoutPayload(ctx, agg),
// 				commandWithoutPayload(ctx, agg),
// 			},
// 			workers: 1,
// 		},
// 		{
// 			name:   "5 events - payload - 1 worker",
// 			client: localClient,
// 			cmds: []eventstore.Command{
// 				cmdWithPayload(ctx, agg),
// 				cmdWithPayload(ctx, agg),
// 				cmdWithPayload(ctx, agg),
// 				cmdWithPayload(ctx, agg),
// 				cmdWithPayload(ctx, agg),
// 			},
// 			workers: 1,
// 		},

// 		{
// 			name:   "1 events - no payload - 2 workers",
// 			client: localClient,
// 			cmds: []eventstore.Command{
// 				commandWithoutPayload(ctx, agg),
// 			},
// 			workers: 2,
// 		},
// 		{
// 			name:   "2 events - payload - 2 workers",
// 			client: localClient,
// 			cmds: []eventstore.Command{
// 				cmdWithPayload(ctx, agg),
// 				cmdWithPayload(ctx, agg),
// 			},
// 			workers: 2,
// 		},
// 		{
// 			name:   "5 event - no payload - 5 workers",
// 			client: localClient,
// 			cmds: []eventstore.Command{
// 				commandWithoutPayload(ctx, agg),
// 				commandWithoutPayload(ctx, agg),
// 				commandWithoutPayload(ctx, agg),
// 				commandWithoutPayload(ctx, agg),
// 				commandWithoutPayload(ctx, agg),
// 			},
// 			workers: 5,
// 		},
// 		{
// 			name:   "5 event - payload - 5 workers",
// 			client: localClient,
// 			cmds: []eventstore.Command{
// 				cmdWithPayload(ctx, agg),
// 				cmdWithPayload(ctx, agg),
// 				cmdWithPayload(ctx, agg),
// 				cmdWithPayload(ctx, agg),
// 				cmdWithPayload(ctx, agg),
// 			},
// 			workers: 5,
// 		},
// 	}
// 	for _, tt := range tests {
// 		execTest(b, tt.client, tt.workers, tt.name, tt.cmds)
// 	}
// }

// func Benchmark_Eventstore_PushMultipleAggregates(b *testing.B) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	ctx = authz.SetCtxData(ctx, authz.CtxData{UserID: "adlerhurst", OrgID: "myorg"})
// 	defer cancel()

// 	agg1 := eventstore.NewAggregate(ctx, "ng5PD", "test", "v1")
// 	agg2 := eventstore.NewAggregate(ctx, "e4epE", "test", "v1")
// 	agg3 := eventstore.NewAggregate(ctx, "vE0uJ", "test", "v1")

// 	tests := []struct {
// 		name    string
// 		client  *sql.DB
// 		cmds    []eventstore.Command
// 		workers int
// 	}{
// 		{
// 			name:   "2 events - no payload - 2 aggs - 1 worker",
// 			client: localClient,
// 			cmds: []eventstore.Command{
// 				commandWithoutPayload(ctx, agg1),
// 				commandWithoutPayload(ctx, agg2),
// 			},
// 			workers: 1,
// 		},
// 		{
// 			name:   "2 events - payload - 2 aggs - 1 worker",
// 			client: localClient,
// 			cmds: []eventstore.Command{
// 				cmdWithPayload(ctx, agg1),
// 				cmdWithPayload(ctx, agg2),
// 			},
// 			workers: 1,
// 		},
// 		{
// 			name:   "5 events - no payload - 3 aggs - 1 worker",
// 			client: localClient,
// 			cmds: []eventstore.Command{
// 				commandWithoutPayload(ctx, agg1),
// 				commandWithoutPayload(ctx, agg2),
// 				commandWithoutPayload(ctx, agg3),
// 				commandWithoutPayload(ctx, agg1),
// 				commandWithoutPayload(ctx, agg2),
// 			},
// 			workers: 1,
// 		},
// 		{
// 			name:   "5 events - payload - 3 aggs - 1 worker",
// 			client: localClient,
// 			cmds: []eventstore.Command{
// 				commandWithoutPayload(ctx, agg1),
// 				commandWithoutPayload(ctx, agg2),
// 				commandWithoutPayload(ctx, agg3),
// 				commandWithoutPayload(ctx, agg1),
// 				commandWithoutPayload(ctx, agg2),
// 			},
// 			workers: 1,
// 		},

// 		{
// 			name:   "2 events - no payload - 2 aggs - 2 workers",
// 			client: localClient,
// 			cmds: []eventstore.Command{
// 				commandWithoutPayload(ctx, agg1),
// 				commandWithoutPayload(ctx, agg2),
// 			},
// 			workers: 2,
// 		},
// 		{
// 			name:   "2 events - payload - 2 aggs - 2 workers",
// 			client: localClient,
// 			cmds: []eventstore.Command{
// 				cmdWithPayload(ctx, agg1),
// 				cmdWithPayload(ctx, agg2),
// 			},
// 			workers: 2,
// 		},
// 		{
// 			name:   "5 events - no payload - 3 aggs - 5 workers",
// 			client: localClient,
// 			cmds: []eventstore.Command{
// 				commandWithoutPayload(ctx, agg1),
// 				commandWithoutPayload(ctx, agg2),
// 				commandWithoutPayload(ctx, agg3),
// 				commandWithoutPayload(ctx, agg1),
// 				commandWithoutPayload(ctx, agg2),
// 			},
// 			workers: 5,
// 		},
// 		{
// 			name:   "5 events - payload - 3 aggs - 5 workers",
// 			client: localClient,
// 			cmds: []eventstore.Command{
// 				commandWithoutPayload(ctx, agg1),
// 				commandWithoutPayload(ctx, agg2),
// 				commandWithoutPayload(ctx, agg3),
// 				commandWithoutPayload(ctx, agg1),
// 				commandWithoutPayload(ctx, agg2),
// 			},
// 			workers: 5,
// 		},

// 		{
// 			name:   "2 events - no payload - 2 aggs - 10 workers",
// 			client: localClient,
// 			cmds: []eventstore.Command{
// 				commandWithoutPayload(ctx, agg1),
// 				commandWithoutPayload(ctx, agg2),
// 			},
// 			workers: 10,
// 		},
// 		{
// 			name:   "2 events - payload - 2 aggs - 10 wokers",
// 			client: localClient,
// 			cmds: []eventstore.Command{
// 				cmdWithPayload(ctx, agg1),
// 				cmdWithPayload(ctx, agg2),
// 			},
// 			workers: 10,
// 		},
// 		{
// 			name:   "6 events - no payload - 3 aggs - 10 wokers",
// 			client: localClient,
// 			cmds: []eventstore.Command{
// 				commandWithoutPayload(ctx, agg1),
// 				commandWithoutPayload(ctx, agg2),
// 				commandWithoutPayload(ctx, agg3),
// 				commandWithoutPayload(ctx, agg1),
// 				commandWithoutPayload(ctx, agg2),
// 				commandWithoutPayload(ctx, agg3),
// 			},
// 			workers: 10,
// 		},
// 		{
// 			name:   "6 events - payload - 3 aggs - 10 workers",
// 			client: localClient,
// 			cmds: []eventstore.Command{
// 				commandWithoutPayload(ctx, agg1),
// 				commandWithoutPayload(ctx, agg2),
// 				commandWithoutPayload(ctx, agg3),
// 				commandWithoutPayload(ctx, agg1),
// 				commandWithoutPayload(ctx, agg2),
// 				commandWithoutPayload(ctx, agg3),
// 			},
// 			workers: 10,
// 		},
// 	}
// 	for _, tt := range tests {
// 		execTest(b, tt.client, tt.workers, tt.name, tt.cmds)
// 	}
// }
