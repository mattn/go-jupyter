package jupyter

type Error struct {
	Message string      `json:"message"`
	Reason  interface{} `json:"reason"`
}
