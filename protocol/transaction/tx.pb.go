package transaction

import (
	"io"
	"github.com/golang/protobuf/proto"
)

type ValueSource struct {
	Ref      *Hash        `protobuf:"bytes,1,opt,name=ref" json:"ref,omitempty"`
	Value    *AssetAmount `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	Position uint64       `protobuf:"varint,3,opt,name=position" json:"position,omitempty"`
}


type ValueDestination struct {
	Ref      *Hash        `protobuf:"bytes,1,opt,name=ref" json:"ref,omitempty"`
	Value    *AssetAmount `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	Position uint64       `protobuf:"varint,3,opt,name=position" json:"position,omitempty"`
}
/****Program****/
type Program struct {
	VmVersion uint64 `protobuf:"varint,1,opt,name=vm_version,json=vmVersion" json:"vm_version,omitempty"`
	Code      []byte `protobuf:"bytes,2,opt,name=code,proto3" json:"code,omitempty"`
}
func (m *Program) Reset()                    { *m = Program{} }
func (m *Program) String() string            { return proto.CompactTextString(m) }
func (*Program) ProtoMessage()               {}
func (*Program) Descriptor() ([]byte, []int) { return nil, []int{1} }

/****Output****/
type Output struct {
	Source         *ValueSource `protobuf:"bytes,1,opt,name=source" json:"source,omitempty"`
	ControlProgram *Program     `protobuf:"bytes,2,opt,name=control_program,json=controlProgram" json:"control_program,omitempty"`
	Ordinal        uint64       `protobuf:"varint,3,opt,name=ordinal" json:"ordinal,omitempty"`
}

func (Output) typ() string { return "output1" }
func (o *Output) writeForHash(w io.Writer) {
}
// NewOutput creates a new Output.
func NewOutput(source *ValueSource, controlProgram *Program, ordinal uint64) *Output {
	return &Output{
		Source:         source,
		ControlProgram: controlProgram,
		Ordinal:        ordinal,
	}
}
func (m *Output) Reset()                    { *m = Output{} }
func (m *Output) String() string            { return proto.CompactTextString(m) }
func (*Output) ProtoMessage()               {}
func (*Output) Descriptor() ([]byte, []int) { return nil, []int{13} }



/****Spend****/
type Spend struct {
	SpentOutputId      *Hash             `protobuf:"bytes,1,opt,name=spent_output_id,json=spentOutputId" json:"spent_output_id,omitempty"`
	WitnessDestination *ValueDestination `protobuf:"bytes,2,opt,name=witness_destination,json=witnessDestination" json:"witness_destination,omitempty"`
	WitnessArguments   [][]byte          `protobuf:"bytes,3,rep,name=witness_arguments,json=witnessArguments,proto3" json:"witness_arguments,omitempty"`
	Ordinal            uint64            `protobuf:"varint,4,opt,name=ordinal" json:"ordinal,omitempty"`
}

func (Spend) typ() string { return "spend1" }
func (s *Spend) writeForHash(w io.Writer) {
}
// SetDestination will link the spend to the output
func (s *Spend) SetDestination(id *Hash, val *AssetAmount, pos uint64) {
	s.WitnessDestination = &ValueDestination{
		Ref:      id,
		Value:    val,
		Position: pos,
	}
}

func (m *Spend) Reset()                    { *m = Spend{} }
func (m *Spend) String() string            { return proto.CompactTextString(m) }
func (*Spend) ProtoMessage()               {}
func (*Spend) Descriptor() ([]byte, []int) { return nil, []int{16} }


// NewSpend creates a new Spend.
func NewSpend(spentOutputID *Hash, ordinal uint64) *Spend {
	return &Spend{
		SpentOutputId: spentOutputID,
		Ordinal:       ordinal,
	}
}


/****Mux****/
type Mux struct {
	Sources             []*ValueSource      `protobuf:"bytes,1,rep,name=sources" json:"sources,omitempty"`
	Program             *Program            `protobuf:"bytes,2,opt,name=program" json:"program,omitempty"`
	WitnessDestinations []*ValueDestination `protobuf:"bytes,3,rep,name=witness_destinations,json=witnessDestinations" json:"witness_destinations,omitempty"`
	WitnessArguments    [][]byte            `protobuf:"bytes,4,rep,name=witness_arguments,json=witnessArguments,proto3" json:"witness_arguments,omitempty"`
}

func (Mux) typ() string { return "mux1" }
func (m *Mux) writeForHash(w io.Writer) {
}
func (m *Mux) Reset()                    { *m = Mux{} }
func (m *Mux) String() string            { return proto.CompactTextString(m) }
func (*Mux) ProtoMessage()               {}
func (*Mux) Descriptor() ([]byte, []int) { return nil, []int{16} }


// NewMux creates a new Mux.
func NewMux(sources []*ValueSource, program *Program) *Mux {
	return &Mux{
		Sources: sources,
		Program: program,
	}
}

/****TxHeader****/
type TxHeader struct {
	Version        uint64  `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	SerializedSize uint64  `protobuf:"varint,2,opt,name=serialized_size,json=serializedSize" json:"serialized_size,omitempty"`
	TimeRange      uint64  `protobuf:"varint,3,opt,name=time_range,json=timeRange" json:"time_range,omitempty"`
	ResultIds      []*Hash `protobuf:"bytes,4,rep,name=result_ids,json=resultIds" json:"result_ids,omitempty"`
}

func (m *TxHeader) Reset()                    { *m = TxHeader{} }
func (m *TxHeader) String() string            { return proto.CompactTextString(m) }
func (*TxHeader) ProtoMessage()               {}
func (*TxHeader) Descriptor() ([]byte, []int) { return nil, []int{8} }

func (TxHeader) typ() string { return "txheader" }
func (h *TxHeader) writeForHash(w io.Writer) {
}

// NewTxHeader creates an new TxHeader.
func NewTxHeader(version, serializedSize, timeRange uint64, resultIDs []*Hash) *TxHeader {
	return &TxHeader{
		Version:        version,
		SerializedSize: serializedSize,
		TimeRange:      timeRange,
		ResultIds:      resultIDs,
	}
}
