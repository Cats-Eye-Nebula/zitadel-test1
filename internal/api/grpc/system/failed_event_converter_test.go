package system_test

import (
	"testing"
	"time"

	system_grpc "github.com/zitadel/zitadel/internal/api/grpc/system"
	"github.com/zitadel/zitadel/internal/test"
	"github.com/zitadel/zitadel/internal/view/model"
	system_pb "github.com/zitadel/zitadel/pkg/grpc/system"
)

func TestFailedEventsToPbFields(t *testing.T) {
	type args struct {
		failedEvents []*model.FailedEvent
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "all fields",
			args: args{
				failedEvents: []*model.FailedEvent{
					{
						Database:       "admin",
						ViewName:       "users",
						FailedSequence: 456,
						FailureCount:   5,
						LastFailed:     time.Now(),
						ErrMsg:         "some error",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := system_grpc.FailedEventsViewToPb(tt.args.failedEvents)
			for _, g := range got {
				test.AssertFieldsMapped(t, g)
			}
		})
	}
}

func TestFailedEventToPbFields(t *testing.T) {
	type args struct {
		failedEvent *model.FailedEvent
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"all fields",
			args{
				failedEvent: &model.FailedEvent{
					Database:       "admin",
					ViewName:       "users",
					FailedSequence: 456,
					FailureCount:   5,
					LastFailed:     time.Now(),
					ErrMsg:         "some error",
				},
			},
		},
	}
	for _, tt := range tests {
		converted := system_grpc.FailedEventViewToPb(tt.args.failedEvent)
		test.AssertFieldsMapped(t, converted)
	}
}

func TestRemoveFailedEventRequestToModelFields(t *testing.T) {
	type args struct {
		req *system_pb.RemoveFailedEventRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"all fields",
			args{
				req: &system_pb.RemoveFailedEventRequest{
					Database:       "admin",
					ViewName:       "users",
					FailedSequence: 456,
					InstanceId:     "instanceID",
				},
			},
		},
	}
	for _, tt := range tests {
		converted := system_grpc.RemoveFailedEventRequestToModel(tt.args.req)
		test.AssertFieldsMapped(t, converted, "FailureCount", "LastFailed", "ErrMsg")
	}
}
