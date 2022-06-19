package ws

type ThingAuthMsg struct {
	Token   string `json:"token"`
	ThingId string `json:"thingId"`
}

const WsUnauthorized = "unauthorized"
