package ws

type ThingUpdateMsg struct {
	Token   string `json:"token"`
	ThingId string `json:"thingId"`
}

const WsUnauthorized = "unauthorized"
