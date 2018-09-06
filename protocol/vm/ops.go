package vm

type Op uint8

const(
	OP_0     Op = 0x00 // synonym
	OP_1    Op = 0x51

	OP_DATA_1  Op = 0x01
	OP_TRUE Op = 0x51 // synonym

	OP_PUSHDATA1 Op = 0x4c
	OP_PUSHDATA2 Op = 0x4d
	OP_PUSHDATA4 Op = 0x4e


	OP_FAIL           Op = 0x6a
)