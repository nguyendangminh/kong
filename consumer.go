package kong

type Consumer struct {
	Id string `json:"id"`
	CustomId string `json:"custom_id"`
	CreatedAt int64 `json:"created_at"`
}