package lib

type target string

const (
	A      target = "A"
	B      target = "B"
	C      target = "C"
	D      target = "D"
	E      target = "E"
	F      target = "F"
	H      target = "H"
	L      target = "L"
	AF     target = "AF"
	BC     target = "BC"
	DE     target = "DE"
	HL     target = "HL"
	SP     target = "SP"
	e8     target = "e8"
	SPe8   target = "SP+e8"
	n      target = "n"
	nn     target = "nn"
	C_M    target = "(C)"
	BC_M   target = "(BC)"
	DE_M   target = "(DE)"
	HL_M   target = "(HL)"
	HLP_M  target = "(HL+)"
	HLM_M  target = "(HL-)"
	n_M    target = "(n)"
	nn_M   target = "(nn)"
	nn_M16 target = "(nn)16"
	None   target = "none"
)
