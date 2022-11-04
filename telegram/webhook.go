package telegram // import "heytobi.dev/fuse/telegram"

// Webhook defines an endpoint for receiving telegram updates.
// See https://core.telegram.org/bots/api#setwebhook
type Webhook struct {
	Url                string   `json:"url"`
	IPAddress          string   `json:"ip_address"`
	MaxConnections     int      `json:"max_connections"`
	AllowedUpdates     []string `json:"allowed_updates"`
	DropPendingUpdates bool     `json:"drop_pending_updates"`
}

type deleteWebhookRequest struct {
	DropPendingUpdates bool `json:"drop_pending_updates"`
}

type webhookResponse struct {
	Ok          bool   `json:"ok"`
	Result      bool   `json:"result"`
	Description string `json:"description"`
}
