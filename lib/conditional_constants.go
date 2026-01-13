package lib

type conditional string

const (
	cond_None conditional = "None"
	cond_C    conditional = "C"
	cond_NC   conditional = "NC"
	cond_Z    conditional = "Z"
	cond_NZ   conditional = "NZ"
)
