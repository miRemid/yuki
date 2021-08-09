package response

type StatusCode uint32

const (
	StatusOK        StatusCode = 0
	StatusBindError StatusCode = 100
)
