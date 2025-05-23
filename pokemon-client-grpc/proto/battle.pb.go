// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v4.25.7
// source: proto/battle.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PokemonUpdate struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	X             float64                `protobuf:"fixed64,3,opt,name=x,proto3" json:"x,omitempty"`
	Y             float64                `protobuf:"fixed64,4,opt,name=y,proto3" json:"y,omitempty"`
	Hp            int32                  `protobuf:"varint,5,opt,name=hp,proto3" json:"hp,omitempty"`
	Action        string                 `protobuf:"bytes,6,opt,name=action,proto3" json:"action,omitempty"`
	TargetId      string                 `protobuf:"bytes,7,opt,name=target_id,json=targetId,proto3" json:"target_id,omitempty"`
	Timestamp     int64                  `protobuf:"varint,8,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PokemonUpdate) Reset() {
	*x = PokemonUpdate{}
	mi := &file_proto_battle_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PokemonUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PokemonUpdate) ProtoMessage() {}

func (x *PokemonUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_proto_battle_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PokemonUpdate.ProtoReflect.Descriptor instead.
func (*PokemonUpdate) Descriptor() ([]byte, []int) {
	return file_proto_battle_proto_rawDescGZIP(), []int{0}
}

func (x *PokemonUpdate) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PokemonUpdate) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PokemonUpdate) GetX() float64 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *PokemonUpdate) GetY() float64 {
	if x != nil {
		return x.Y
	}
	return 0
}

func (x *PokemonUpdate) GetHp() int32 {
	if x != nil {
		return x.Hp
	}
	return 0
}

func (x *PokemonUpdate) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

func (x *PokemonUpdate) GetTargetId() string {
	if x != nil {
		return x.TargetId
	}
	return ""
}

func (x *PokemonUpdate) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

type BattleStatus struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Pokemons      []*PokemonUpdate       `protobuf:"bytes,1,rep,name=pokemons,proto3" json:"pokemons,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BattleStatus) Reset() {
	*x = BattleStatus{}
	mi := &file_proto_battle_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BattleStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BattleStatus) ProtoMessage() {}

func (x *BattleStatus) ProtoReflect() protoreflect.Message {
	mi := &file_proto_battle_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BattleStatus.ProtoReflect.Descriptor instead.
func (*BattleStatus) Descriptor() ([]byte, []int) {
	return file_proto_battle_proto_rawDescGZIP(), []int{1}
}

func (x *BattleStatus) GetPokemons() []*PokemonUpdate {
	if x != nil {
		return x.Pokemons
	}
	return nil
}

var File_proto_battle_proto protoreflect.FileDescriptor

const file_proto_battle_proto_rawDesc = "" +
	"\n" +
	"\x12proto/battle.proto\"\xb2\x01\n" +
	"\rPokemonUpdate\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\f\n" +
	"\x01x\x18\x03 \x01(\x01R\x01x\x12\f\n" +
	"\x01y\x18\x04 \x01(\x01R\x01y\x12\x0e\n" +
	"\x02hp\x18\x05 \x01(\x05R\x02hp\x12\x16\n" +
	"\x06action\x18\x06 \x01(\tR\x06action\x12\x1b\n" +
	"\ttarget_id\x18\a \x01(\tR\btargetId\x12\x1c\n" +
	"\ttimestamp\x18\b \x01(\x03R\ttimestamp\":\n" +
	"\fBattleStatus\x12*\n" +
	"\bpokemons\x18\x01 \x03(\v2\x0e.PokemonUpdateR\bpokemons2B\n" +
	"\rBattleService\x121\n" +
	"\fBattleStream\x12\x0e.PokemonUpdate\x1a\r.BattleStatus(\x010\x01B\tZ\a./protob\x06proto3"

var (
	file_proto_battle_proto_rawDescOnce sync.Once
	file_proto_battle_proto_rawDescData []byte
)

func file_proto_battle_proto_rawDescGZIP() []byte {
	file_proto_battle_proto_rawDescOnce.Do(func() {
		file_proto_battle_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_battle_proto_rawDesc), len(file_proto_battle_proto_rawDesc)))
	})
	return file_proto_battle_proto_rawDescData
}

var file_proto_battle_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_battle_proto_goTypes = []any{
	(*PokemonUpdate)(nil), // 0: PokemonUpdate
	(*BattleStatus)(nil),  // 1: BattleStatus
}
var file_proto_battle_proto_depIdxs = []int32{
	0, // 0: BattleStatus.pokemons:type_name -> PokemonUpdate
	0, // 1: BattleService.BattleStream:input_type -> PokemonUpdate
	1, // 2: BattleService.BattleStream:output_type -> BattleStatus
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_battle_proto_init() }
func file_proto_battle_proto_init() {
	if File_proto_battle_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_battle_proto_rawDesc), len(file_proto_battle_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_battle_proto_goTypes,
		DependencyIndexes: file_proto_battle_proto_depIdxs,
		MessageInfos:      file_proto_battle_proto_msgTypes,
	}.Build()
	File_proto_battle_proto = out.File
	file_proto_battle_proto_goTypes = nil
	file_proto_battle_proto_depIdxs = nil
}
