// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: zitadel/action.proto

package action

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	message "github.com/zitadel/zitadel/pkg/grpc/message"
	object "github.com/zitadel/zitadel/pkg/grpc/object"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ActionState int32

const (
	ActionState_ACTION_STATE_UNSPECIFIED ActionState = 0
	ActionState_ACTION_STATE_INACTIVE    ActionState = 1
	ActionState_ACTION_STATE_ACTIVE      ActionState = 2
)

// Enum value maps for ActionState.
var (
	ActionState_name = map[int32]string{
		0: "ACTION_STATE_UNSPECIFIED",
		1: "ACTION_STATE_INACTIVE",
		2: "ACTION_STATE_ACTIVE",
	}
	ActionState_value = map[string]int32{
		"ACTION_STATE_UNSPECIFIED": 0,
		"ACTION_STATE_INACTIVE":    1,
		"ACTION_STATE_ACTIVE":      2,
	}
)

func (x ActionState) Enum() *ActionState {
	p := new(ActionState)
	*p = x
	return p
}

func (x ActionState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ActionState) Descriptor() protoreflect.EnumDescriptor {
	return file_zitadel_action_proto_enumTypes[0].Descriptor()
}

func (ActionState) Type() protoreflect.EnumType {
	return &file_zitadel_action_proto_enumTypes[0]
}

func (x ActionState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ActionState.Descriptor instead.
func (ActionState) EnumDescriptor() ([]byte, []int) {
	return file_zitadel_action_proto_rawDescGZIP(), []int{0}
}

type ActionFieldName int32

const (
	ActionFieldName_ACTION_FIELD_NAME_UNSPECIFIED ActionFieldName = 0
	ActionFieldName_ACTION_FIELD_NAME_NAME        ActionFieldName = 1
	ActionFieldName_ACTION_FIELD_NAME_ID          ActionFieldName = 2
	ActionFieldName_ACTION_FIELD_NAME_STATE       ActionFieldName = 3
)

// Enum value maps for ActionFieldName.
var (
	ActionFieldName_name = map[int32]string{
		0: "ACTION_FIELD_NAME_UNSPECIFIED",
		1: "ACTION_FIELD_NAME_NAME",
		2: "ACTION_FIELD_NAME_ID",
		3: "ACTION_FIELD_NAME_STATE",
	}
	ActionFieldName_value = map[string]int32{
		"ACTION_FIELD_NAME_UNSPECIFIED": 0,
		"ACTION_FIELD_NAME_NAME":        1,
		"ACTION_FIELD_NAME_ID":          2,
		"ACTION_FIELD_NAME_STATE":       3,
	}
)

func (x ActionFieldName) Enum() *ActionFieldName {
	p := new(ActionFieldName)
	*p = x
	return p
}

func (x ActionFieldName) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ActionFieldName) Descriptor() protoreflect.EnumDescriptor {
	return file_zitadel_action_proto_enumTypes[1].Descriptor()
}

func (ActionFieldName) Type() protoreflect.EnumType {
	return &file_zitadel_action_proto_enumTypes[1]
}

func (x ActionFieldName) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ActionFieldName.Descriptor instead.
func (ActionFieldName) EnumDescriptor() ([]byte, []int) {
	return file_zitadel_action_proto_rawDescGZIP(), []int{1}
}

type FlowState int32

const (
	FlowState_FLOW_STATE_UNSPECIFIED FlowState = 0
	FlowState_FLOW_STATE_INACTIVE    FlowState = 1
	FlowState_FLOW_STATE_ACTIVE      FlowState = 2
)

// Enum value maps for FlowState.
var (
	FlowState_name = map[int32]string{
		0: "FLOW_STATE_UNSPECIFIED",
		1: "FLOW_STATE_INACTIVE",
		2: "FLOW_STATE_ACTIVE",
	}
	FlowState_value = map[string]int32{
		"FLOW_STATE_UNSPECIFIED": 0,
		"FLOW_STATE_INACTIVE":    1,
		"FLOW_STATE_ACTIVE":      2,
	}
)

func (x FlowState) Enum() *FlowState {
	p := new(FlowState)
	*p = x
	return p
}

func (x FlowState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FlowState) Descriptor() protoreflect.EnumDescriptor {
	return file_zitadel_action_proto_enumTypes[2].Descriptor()
}

func (FlowState) Type() protoreflect.EnumType {
	return &file_zitadel_action_proto_enumTypes[2]
}

func (x FlowState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FlowState.Descriptor instead.
func (FlowState) EnumDescriptor() ([]byte, []int) {
	return file_zitadel_action_proto_rawDescGZIP(), []int{2}
}

type Action struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            string                `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Details       *object.ObjectDetails `protobuf:"bytes,2,opt,name=details,proto3" json:"details,omitempty"`
	State         ActionState           `protobuf:"varint,3,opt,name=state,proto3,enum=zitadel.action.v1.ActionState" json:"state,omitempty"`
	Name          string                `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	Script        string                `protobuf:"bytes,5,opt,name=script,proto3" json:"script,omitempty"`
	Timeout       *durationpb.Duration  `protobuf:"bytes,6,opt,name=timeout,proto3" json:"timeout,omitempty"`
	AllowedToFail bool                  `protobuf:"varint,7,opt,name=allowed_to_fail,json=allowedToFail,proto3" json:"allowed_to_fail,omitempty"`
}

func (x *Action) Reset() {
	*x = Action{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zitadel_action_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Action) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Action) ProtoMessage() {}

func (x *Action) ProtoReflect() protoreflect.Message {
	mi := &file_zitadel_action_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Action.ProtoReflect.Descriptor instead.
func (*Action) Descriptor() ([]byte, []int) {
	return file_zitadel_action_proto_rawDescGZIP(), []int{0}
}

func (x *Action) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Action) GetDetails() *object.ObjectDetails {
	if x != nil {
		return x.Details
	}
	return nil
}

func (x *Action) GetState() ActionState {
	if x != nil {
		return x.State
	}
	return ActionState_ACTION_STATE_UNSPECIFIED
}

func (x *Action) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Action) GetScript() string {
	if x != nil {
		return x.Script
	}
	return ""
}

func (x *Action) GetTimeout() *durationpb.Duration {
	if x != nil {
		return x.Timeout
	}
	return nil
}

func (x *Action) GetAllowedToFail() bool {
	if x != nil {
		return x.AllowedToFail
	}
	return false
}

type ActionIDQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ActionIDQuery) Reset() {
	*x = ActionIDQuery{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zitadel_action_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActionIDQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActionIDQuery) ProtoMessage() {}

func (x *ActionIDQuery) ProtoReflect() protoreflect.Message {
	mi := &file_zitadel_action_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActionIDQuery.ProtoReflect.Descriptor instead.
func (*ActionIDQuery) Descriptor() ([]byte, []int) {
	return file_zitadel_action_proto_rawDescGZIP(), []int{1}
}

func (x *ActionIDQuery) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type ActionNameQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Method object.TextQueryMethod `protobuf:"varint,2,opt,name=method,proto3,enum=zitadel.v1.TextQueryMethod" json:"method,omitempty"`
}

func (x *ActionNameQuery) Reset() {
	*x = ActionNameQuery{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zitadel_action_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActionNameQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActionNameQuery) ProtoMessage() {}

func (x *ActionNameQuery) ProtoReflect() protoreflect.Message {
	mi := &file_zitadel_action_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActionNameQuery.ProtoReflect.Descriptor instead.
func (*ActionNameQuery) Descriptor() ([]byte, []int) {
	return file_zitadel_action_proto_rawDescGZIP(), []int{2}
}

func (x *ActionNameQuery) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ActionNameQuery) GetMethod() object.TextQueryMethod {
	if x != nil {
		return x.Method
	}
	return object.TextQueryMethod(0)
}

// ActionStateQuery is always equals
type ActionStateQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State ActionState `protobuf:"varint,1,opt,name=state,proto3,enum=zitadel.action.v1.ActionState" json:"state,omitempty"`
}

func (x *ActionStateQuery) Reset() {
	*x = ActionStateQuery{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zitadel_action_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActionStateQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActionStateQuery) ProtoMessage() {}

func (x *ActionStateQuery) ProtoReflect() protoreflect.Message {
	mi := &file_zitadel_action_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActionStateQuery.ProtoReflect.Descriptor instead.
func (*ActionStateQuery) Descriptor() ([]byte, []int) {
	return file_zitadel_action_proto_rawDescGZIP(), []int{3}
}

func (x *ActionStateQuery) GetState() ActionState {
	if x != nil {
		return x.State
	}
	return ActionState_ACTION_STATE_UNSPECIFIED
}

type Flow struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id of the flow type
	Type           *FlowType             `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Details        *object.ObjectDetails `protobuf:"bytes,2,opt,name=details,proto3" json:"details,omitempty"`
	State          FlowState             `protobuf:"varint,3,opt,name=state,proto3,enum=zitadel.action.v1.FlowState" json:"state,omitempty"`
	TriggerActions []*TriggerAction      `protobuf:"bytes,4,rep,name=trigger_actions,json=triggerActions,proto3" json:"trigger_actions,omitempty"`
}

func (x *Flow) Reset() {
	*x = Flow{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zitadel_action_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Flow) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Flow) ProtoMessage() {}

func (x *Flow) ProtoReflect() protoreflect.Message {
	mi := &file_zitadel_action_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Flow.ProtoReflect.Descriptor instead.
func (*Flow) Descriptor() ([]byte, []int) {
	return file_zitadel_action_proto_rawDescGZIP(), []int{4}
}

func (x *Flow) GetType() *FlowType {
	if x != nil {
		return x.Type
	}
	return nil
}

func (x *Flow) GetDetails() *object.ObjectDetails {
	if x != nil {
		return x.Details
	}
	return nil
}

func (x *Flow) GetState() FlowState {
	if x != nil {
		return x.State
	}
	return FlowState_FLOW_STATE_UNSPECIFIED
}

func (x *Flow) GetTriggerActions() []*TriggerAction {
	if x != nil {
		return x.TriggerActions
	}
	return nil
}

type FlowType struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// identifier of the type
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// key and name of the type
	Name *message.LocalizedMessage `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *FlowType) Reset() {
	*x = FlowType{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zitadel_action_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FlowType) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FlowType) ProtoMessage() {}

func (x *FlowType) ProtoReflect() protoreflect.Message {
	mi := &file_zitadel_action_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FlowType.ProtoReflect.Descriptor instead.
func (*FlowType) Descriptor() ([]byte, []int) {
	return file_zitadel_action_proto_rawDescGZIP(), []int{5}
}

func (x *FlowType) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *FlowType) GetName() *message.LocalizedMessage {
	if x != nil {
		return x.Name
	}
	return nil
}

type TriggerType struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// identifier of the type
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// key and name of the type
	Name *message.LocalizedMessage `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *TriggerType) Reset() {
	*x = TriggerType{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zitadel_action_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TriggerType) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TriggerType) ProtoMessage() {}

func (x *TriggerType) ProtoReflect() protoreflect.Message {
	mi := &file_zitadel_action_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TriggerType.ProtoReflect.Descriptor instead.
func (*TriggerType) Descriptor() ([]byte, []int) {
	return file_zitadel_action_proto_rawDescGZIP(), []int{6}
}

func (x *TriggerType) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TriggerType) GetName() *message.LocalizedMessage {
	if x != nil {
		return x.Name
	}
	return nil
}

type TriggerAction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id of the trigger type
	TriggerType *TriggerType `protobuf:"bytes,1,opt,name=trigger_type,json=triggerType,proto3" json:"trigger_type,omitempty"`
	Actions     []*Action    `protobuf:"bytes,2,rep,name=actions,proto3" json:"actions,omitempty"`
}

func (x *TriggerAction) Reset() {
	*x = TriggerAction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zitadel_action_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TriggerAction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TriggerAction) ProtoMessage() {}

func (x *TriggerAction) ProtoReflect() protoreflect.Message {
	mi := &file_zitadel_action_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TriggerAction.ProtoReflect.Descriptor instead.
func (*TriggerAction) Descriptor() ([]byte, []int) {
	return file_zitadel_action_proto_rawDescGZIP(), []int{7}
}

func (x *TriggerAction) GetTriggerType() *TriggerType {
	if x != nil {
		return x.TriggerType
	}
	return nil
}

func (x *TriggerAction) GetActions() []*Action {
	if x != nil {
		return x.Actions
	}
	return nil
}

var File_zitadel_action_proto protoreflect.FileDescriptor

var file_zitadel_action_proto_rawDesc = []byte{
	0x0a, 0x14, 0x7a, 0x69, 0x74, 0x61, 0x64, 0x65, 0x6c, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x7a, 0x69, 0x74, 0x61, 0x64, 0x65, 0x6c, 0x2e,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x1a, 0x14, 0x7a, 0x69, 0x74, 0x61, 0x64,
	0x65, 0x6c, 0x2f, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x15, 0x7a, 0x69, 0x74, 0x61, 0x64, 0x65, 0x6c, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e,
	0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e,
	0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xa2, 0x04, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x28, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x18, 0x92, 0x41, 0x15, 0x4a, 0x13, 0x22, 0x36, 0x39,
	0x36, 0x32, 0x39, 0x30, 0x32, 0x33, 0x39, 0x30, 0x36, 0x34, 0x38, 0x38, 0x33, 0x33, 0x34, 0x22,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x33, 0x0a, 0x07, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x7a, 0x69, 0x74, 0x61, 0x64, 0x65, 0x6c, 0x2e,
	0x76, 0x31, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73,
	0x52, 0x07, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x52, 0x0a, 0x05, 0x73, 0x74, 0x61,
	0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x7a, 0x69, 0x74, 0x61, 0x64,
	0x65, 0x6c, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x42, 0x1c, 0x92, 0x41, 0x19, 0x32, 0x17, 0x74,
	0x68, 0x65, 0x20, 0x73, 0x74, 0x61, 0x74, 0x65, 0x20, 0x6f, 0x66, 0x20, 0x74, 0x68, 0x65, 0x20,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x26, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x12, 0x92, 0x41, 0x0f,
	0x4a, 0x0d, 0x22, 0x6c, 0x6f, 0x67, 0x20, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x22, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x51, 0x0a, 0x06, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x42, 0x39, 0x92, 0x41, 0x36, 0x4a, 0x34, 0x22, 0x66, 0x75, 0x6e,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x20, 0x6c, 0x6f, 0x67, 0x28, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78,
	0x74, 0x2c, 0x20, 0x63, 0x61, 0x6c, 0x6c, 0x73, 0x29, 0x7b, 0x63, 0x6f, 0x6e, 0x73, 0x6f, 0x6c,
	0x65, 0x2e, 0x6c, 0x6f, 0x67, 0x28, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x29, 0x7d, 0x22,
	0x52, 0x06, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x12, 0x78, 0x0a, 0x07, 0x74, 0x69, 0x6d, 0x65,
	0x6f, 0x75, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x42, 0x43, 0x92, 0x41, 0x40, 0x32, 0x3e, 0x61, 0x66, 0x74, 0x65, 0x72,
	0x20, 0x77, 0x68, 0x69, 0x63, 0x68, 0x20, 0x74, 0x69, 0x6d, 0x65, 0x20, 0x74, 0x68, 0x65, 0x20,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x20, 0x77, 0x69, 0x6c, 0x6c, 0x20, 0x62, 0x65, 0x20, 0x74,
	0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x64, 0x20, 0x69, 0x66, 0x20, 0x6e, 0x6f, 0x74,
	0x20, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f,
	0x75, 0x74, 0x12, 0x70, 0x0a, 0x0f, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x5f, 0x74, 0x6f,
	0x5f, 0x66, 0x61, 0x69, 0x6c, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x42, 0x48, 0x92, 0x41, 0x45,
	0x32, 0x43, 0x77, 0x68, 0x65, 0x6e, 0x20, 0x74, 0x72, 0x75, 0x65, 0x2c, 0x20, 0x74, 0x68, 0x65,
	0x20, 0x6e, 0x65, 0x78, 0x74, 0x20, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x20, 0x77, 0x69, 0x6c,
	0x6c, 0x20, 0x62, 0x65, 0x20, 0x63, 0x61, 0x6c, 0x6c, 0x65, 0x64, 0x20, 0x65, 0x76, 0x65, 0x6e,
	0x20, 0x69, 0x66, 0x20, 0x74, 0x68, 0x69, 0x73, 0x20, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x20,
	0x66, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x0d, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x54, 0x6f,
	0x46, 0x61, 0x69, 0x6c, 0x22, 0x41, 0x0a, 0x0d, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44,
	0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x30, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x20, 0x92, 0x41, 0x15, 0x4a, 0x13, 0x22, 0x36, 0x39, 0x36, 0x32, 0x39, 0x30, 0x32,
	0x33, 0x39, 0x30, 0x36, 0x34, 0x38, 0x38, 0x33, 0x33, 0x34, 0x22, 0xfa, 0x42, 0x05, 0x72, 0x03,
	0x18, 0xc8, 0x01, 0x52, 0x02, 0x69, 0x64, 0x22, 0xa7, 0x01, 0x0a, 0x0f, 0x41, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x26, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x12, 0x92, 0x41, 0x07, 0x4a, 0x05,
	0x22, 0x6c, 0x6f, 0x67, 0x22, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0xc8, 0x01, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x6c, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x7a, 0x69, 0x74, 0x61, 0x64, 0x65, 0x6c, 0x2e, 0x76, 0x31,
	0x2e, 0x54, 0x65, 0x78, 0x74, 0x51, 0x75, 0x65, 0x72, 0x79, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x42, 0x37, 0x92, 0x41, 0x2c, 0x32, 0x2a, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x73, 0x20, 0x77,
	0x68, 0x69, 0x63, 0x68, 0x20, 0x74, 0x65, 0x78, 0x74, 0x20, 0x65, 0x71, 0x75, 0x61, 0x6c, 0x69,
	0x74, 0x79, 0x20, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x20, 0x69, 0x73, 0x20, 0x75, 0x73, 0x65,
	0x64, 0xfa, 0x42, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x22, 0x72, 0x0a, 0x10, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x5e, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x7a, 0x69, 0x74, 0x61, 0x64, 0x65, 0x6c, 0x2e, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x42, 0x28, 0x92, 0x41, 0x1d, 0x32, 0x1b, 0x63, 0x75, 0x72, 0x72, 0x65,
	0x6e, 0x74, 0x20, 0x73, 0x74, 0x61, 0x74, 0x65, 0x20, 0x6f, 0x66, 0x20, 0x74, 0x68, 0x65, 0x20,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0xfa, 0x42, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x05,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0xa4, 0x02, 0x0a, 0x04, 0x46, 0x6c, 0x6f, 0x77, 0x12, 0x4c,
	0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x7a,
	0x69, 0x74, 0x61, 0x64, 0x65, 0x6c, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31,
	0x2e, 0x46, 0x6c, 0x6f, 0x77, 0x54, 0x79, 0x70, 0x65, 0x42, 0x1b, 0x92, 0x41, 0x18, 0x32, 0x16,
	0x22, 0x74, 0x68, 0x65, 0x20, 0x74, 0x79, 0x70, 0x65, 0x20, 0x6f, 0x66, 0x20, 0x74, 0x68, 0x65,
	0x20, 0x66, 0x6c, 0x6f, 0x77, 0x22, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x33, 0x0a, 0x07,
	0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e,
	0x7a, 0x69, 0x74, 0x61, 0x64, 0x65, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x07, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c,
	0x73, 0x12, 0x4e, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x1c, 0x2e, 0x7a, 0x69, 0x74, 0x61, 0x64, 0x65, 0x6c, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x6c, 0x6f, 0x77, 0x53, 0x74, 0x61, 0x74, 0x65, 0x42, 0x1a,
	0x92, 0x41, 0x17, 0x32, 0x15, 0x74, 0x68, 0x65, 0x20, 0x73, 0x74, 0x61, 0x74, 0x65, 0x20, 0x6f,
	0x66, 0x20, 0x74, 0x68, 0x65, 0x20, 0x66, 0x6c, 0x6f, 0x77, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x12, 0x49, 0x0a, 0x0f, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x5f, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x7a, 0x69, 0x74,
	0x61, 0x64, 0x65, 0x6c, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x54,
	0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0e, 0x74, 0x72,
	0x69, 0x67, 0x67, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x4c, 0x0a, 0x08,
	0x46, 0x6c, 0x6f, 0x77, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x30, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x7a, 0x69, 0x74, 0x61, 0x64, 0x65, 0x6c,
	0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x64, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x4f, 0x0a, 0x0b, 0x54, 0x72,
	0x69, 0x67, 0x67, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x30, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x7a, 0x69, 0x74, 0x61, 0x64, 0x65,
	0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x64, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x87, 0x01, 0x0a, 0x0d,
	0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x41, 0x0a,
	0x0c, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x7a, 0x69, 0x74, 0x61, 0x64, 0x65, 0x6c, 0x2e, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x0b, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x33, 0x0a, 0x07, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x19, 0x2e, 0x7a, 0x69, 0x74, 0x61, 0x64, 0x65, 0x6c, 0x2e, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2a, 0x5f, 0x0a, 0x0b, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x12, 0x1c, 0x0a, 0x18, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x53,
	0x54, 0x41, 0x54, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44,
	0x10, 0x00, 0x12, 0x19, 0x0a, 0x15, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x53, 0x54, 0x41,
	0x54, 0x45, 0x5f, 0x49, 0x4e, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x01, 0x12, 0x17, 0x0a,
	0x13, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x41, 0x43,
	0x54, 0x49, 0x56, 0x45, 0x10, 0x02, 0x2a, 0x87, 0x01, 0x0a, 0x0f, 0x41, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x1d, 0x41, 0x43,
	0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x46, 0x49, 0x45, 0x4c, 0x44, 0x5f, 0x4e, 0x41, 0x4d, 0x45, 0x5f,
	0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1a, 0x0a,
	0x16, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x46, 0x49, 0x45, 0x4c, 0x44, 0x5f, 0x4e, 0x41,
	0x4d, 0x45, 0x5f, 0x4e, 0x41, 0x4d, 0x45, 0x10, 0x01, 0x12, 0x18, 0x0a, 0x14, 0x41, 0x43, 0x54,
	0x49, 0x4f, 0x4e, 0x5f, 0x46, 0x49, 0x45, 0x4c, 0x44, 0x5f, 0x4e, 0x41, 0x4d, 0x45, 0x5f, 0x49,
	0x44, 0x10, 0x02, 0x12, 0x1b, 0x0a, 0x17, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x46, 0x49,
	0x45, 0x4c, 0x44, 0x5f, 0x4e, 0x41, 0x4d, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x10, 0x03,
	0x2a, 0x57, 0x0a, 0x09, 0x46, 0x6c, 0x6f, 0x77, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x1a, 0x0a,
	0x16, 0x46, 0x4c, 0x4f, 0x57, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50,
	0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x17, 0x0a, 0x13, 0x46, 0x4c, 0x4f,
	0x57, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x49, 0x4e, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45,
	0x10, 0x01, 0x12, 0x15, 0x0a, 0x11, 0x46, 0x4c, 0x4f, 0x57, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45,
	0x5f, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x02, 0x42, 0xb6, 0x01, 0x0a, 0x15, 0x63, 0x6f,
	0x6d, 0x2e, 0x7a, 0x69, 0x74, 0x61, 0x64, 0x65, 0x6c, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x76, 0x31, 0x42, 0x0b, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x7a,
	0x69, 0x74, 0x61, 0x64, 0x65, 0x6c, 0x2f, 0x7a, 0x69, 0x74, 0x61, 0x64, 0x65, 0x6c, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0xa2, 0x02,
	0x03, 0x5a, 0x41, 0x58, 0xaa, 0x02, 0x11, 0x5a, 0x69, 0x74, 0x61, 0x64, 0x65, 0x6c, 0x2e, 0x41,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x11, 0x5a, 0x69, 0x74, 0x61, 0x64,
	0x65, 0x6c, 0x5c, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1d, 0x5a,
	0x69, 0x74, 0x61, 0x64, 0x65, 0x6c, 0x5c, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5c, 0x56, 0x31,
	0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x13, 0x5a,
	0x69, 0x74, 0x61, 0x64, 0x65, 0x6c, 0x3a, 0x3a, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x3a, 0x3a,
	0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_zitadel_action_proto_rawDescOnce sync.Once
	file_zitadel_action_proto_rawDescData = file_zitadel_action_proto_rawDesc
)

func file_zitadel_action_proto_rawDescGZIP() []byte {
	file_zitadel_action_proto_rawDescOnce.Do(func() {
		file_zitadel_action_proto_rawDescData = protoimpl.X.CompressGZIP(file_zitadel_action_proto_rawDescData)
	})
	return file_zitadel_action_proto_rawDescData
}

var file_zitadel_action_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_zitadel_action_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_zitadel_action_proto_goTypes = []interface{}{
	(ActionState)(0),                 // 0: zitadel.action.v1.ActionState
	(ActionFieldName)(0),             // 1: zitadel.action.v1.ActionFieldName
	(FlowState)(0),                   // 2: zitadel.action.v1.FlowState
	(*Action)(nil),                   // 3: zitadel.action.v1.Action
	(*ActionIDQuery)(nil),            // 4: zitadel.action.v1.ActionIDQuery
	(*ActionNameQuery)(nil),          // 5: zitadel.action.v1.ActionNameQuery
	(*ActionStateQuery)(nil),         // 6: zitadel.action.v1.ActionStateQuery
	(*Flow)(nil),                     // 7: zitadel.action.v1.Flow
	(*FlowType)(nil),                 // 8: zitadel.action.v1.FlowType
	(*TriggerType)(nil),              // 9: zitadel.action.v1.TriggerType
	(*TriggerAction)(nil),            // 10: zitadel.action.v1.TriggerAction
	(*object.ObjectDetails)(nil),     // 11: zitadel.v1.ObjectDetails
	(*durationpb.Duration)(nil),      // 12: google.protobuf.Duration
	(object.TextQueryMethod)(0),      // 13: zitadel.v1.TextQueryMethod
	(*message.LocalizedMessage)(nil), // 14: zitadel.v1.LocalizedMessage
}
var file_zitadel_action_proto_depIdxs = []int32{
	11, // 0: zitadel.action.v1.Action.details:type_name -> zitadel.v1.ObjectDetails
	0,  // 1: zitadel.action.v1.Action.state:type_name -> zitadel.action.v1.ActionState
	12, // 2: zitadel.action.v1.Action.timeout:type_name -> google.protobuf.Duration
	13, // 3: zitadel.action.v1.ActionNameQuery.method:type_name -> zitadel.v1.TextQueryMethod
	0,  // 4: zitadel.action.v1.ActionStateQuery.state:type_name -> zitadel.action.v1.ActionState
	8,  // 5: zitadel.action.v1.Flow.type:type_name -> zitadel.action.v1.FlowType
	11, // 6: zitadel.action.v1.Flow.details:type_name -> zitadel.v1.ObjectDetails
	2,  // 7: zitadel.action.v1.Flow.state:type_name -> zitadel.action.v1.FlowState
	10, // 8: zitadel.action.v1.Flow.trigger_actions:type_name -> zitadel.action.v1.TriggerAction
	14, // 9: zitadel.action.v1.FlowType.name:type_name -> zitadel.v1.LocalizedMessage
	14, // 10: zitadel.action.v1.TriggerType.name:type_name -> zitadel.v1.LocalizedMessage
	9,  // 11: zitadel.action.v1.TriggerAction.trigger_type:type_name -> zitadel.action.v1.TriggerType
	3,  // 12: zitadel.action.v1.TriggerAction.actions:type_name -> zitadel.action.v1.Action
	13, // [13:13] is the sub-list for method output_type
	13, // [13:13] is the sub-list for method input_type
	13, // [13:13] is the sub-list for extension type_name
	13, // [13:13] is the sub-list for extension extendee
	0,  // [0:13] is the sub-list for field type_name
}

func init() { file_zitadel_action_proto_init() }
func file_zitadel_action_proto_init() {
	if File_zitadel_action_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_zitadel_action_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Action); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_zitadel_action_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActionIDQuery); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_zitadel_action_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActionNameQuery); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_zitadel_action_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActionStateQuery); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_zitadel_action_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Flow); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_zitadel_action_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FlowType); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_zitadel_action_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TriggerType); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_zitadel_action_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TriggerAction); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_zitadel_action_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_zitadel_action_proto_goTypes,
		DependencyIndexes: file_zitadel_action_proto_depIdxs,
		EnumInfos:         file_zitadel_action_proto_enumTypes,
		MessageInfos:      file_zitadel_action_proto_msgTypes,
	}.Build()
	File_zitadel_action_proto = out.File
	file_zitadel_action_proto_rawDesc = nil
	file_zitadel_action_proto_goTypes = nil
	file_zitadel_action_proto_depIdxs = nil
}
