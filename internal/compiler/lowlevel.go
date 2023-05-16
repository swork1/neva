package compiler

type LLProgram struct {
	Funcs  []LLFunc             // what functions to spawn and how
	Net    []LLConnection       // how ports are connected to each other
	Consts map[string]LLMsg     // predefined data that can be referred by funcs
	Ports  map[LLPortAddr]uint8 // ports and their buffers size
}

type LLPortAddr struct {
	Path, Name string
	Idx        uint8
}

type LLConnection struct {
	SenderSide    LLConnectionSide
	ReceiverSides []LLConnectionSide
}

type LLConnectionSide struct {
	PortAddr  LLPortAddr
	Selectors []LLSelector
}

type LLSelector struct {
	RecField string
	ArrIdx   int
}

// LLFunc is a instantiation object that runtime will use to spawn a function
type LLFunc struct {
	Ref    LLFuncRef // runtime will use this reference to find the function to spawn
	IO     LLFuncIO  // this is the ports function will use to receive and send data
	MsgRef string    // function can receive predefined message at instantiation time
}

type LLFuncRef struct {
	Pkg, Name string
}

type LLFuncIO struct {
	In, Out []LLPortAddr
}

type LLMsg struct {
	Type  LLMsgType
	Bool  bool
	Int   int
	Float float64
	Str   string
	Vec   []string          // ordered list of msg refs
	Map   map[string]string // key -> msg ref
}

type LLMsgType uint8

const (
	LLBoolMsg LLMsgType = iota + 1
	LLIntMsg
	LLFloatMsg
	LLStrMsg
	LLVecMsg
	LLMapMsg
)
