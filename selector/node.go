package selector

type Node struct {
	ID              string `json:"id"`
	RemoteAddr      string `json:"remote_addr" form:"remote_addr" binding:"required"`
	weight          int
	currentWeight   int
	effectiveWeight int
	keys            []uint32
}
