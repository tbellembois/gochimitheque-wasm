package types

type WelcomeAnnounce struct {
	WelcomeAnnounceID   int    `json:"welcomeannounce_id"`
	WelcomeAnnounceText string `json:"welcomeannounce_text"`
	WelcomeAnnounceHTML string `json:"welcomeannounce_html"`
}
