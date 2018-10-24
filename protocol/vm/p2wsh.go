package vm

func P2WSHProgram(hash []byte)([]byte,error){
	builder := NewBuilder()
	builder.AddInt64(0)
	builder.AddData(hash)

	return builder.Build()
}
