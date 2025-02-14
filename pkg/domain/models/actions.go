package models

import "time"

const (
	NopollNode = "none"
  GoogleSheets = "googlesheets"
  NotionToken = "notiontoken"
  NotionOAuth = "notionoauth"
)

var URIListByActionType = map[string]string{
	"googlesheets": "/api/actions/google/sheets",
	"notiontoken":  "/api/actions/notion",
	"notionoauth":  "/api/actions/notion",
	"":             "",
}

type RequestGoogleAction struct {
	ActionID       string `json:"actionid"`
	RequestID      string `json:"requestid"`
	Pollmode       string `json:"pollmode"`
	Selectdocument string `json:"selectdocument"`
	Document       string `json:"document"`
	NameDocument   string `json:"namedocument"`
	ResourceID     string `json:"resourceid"` // document id for example
	Operation      string `json:"operation"`
	Data           string `json:"data"`
	CredentialID   string `json:"credentialid"`
	Sub            string `json:"sub"`
	Type           string `json:"type" binding:"oneof=googlesheets notiontoken notionoauth"`
	WorkflowID     string `json:"workflowid"`
	NodeID         string `json:"nodeid"`
	RedirectURL    string `json:"redirecturl"`
	Status         string `json:"status"` // Default: 'pending'
	CreatedAt      string `json:"createdat"`
	Testmode       bool   `json:"testmode"`
}

type ActionData struct {
	ActionID string `json:"actioid"`
}

type ResponseGetGoogleSheetByID struct {
	Error  string `json:"error"`
	Data   string `json:"data"`
	Status int    `json:"status"`
}

type ActionsCommand struct {
	Timestamp time.Time            `json:"timestamp,omitempty"`
	Actions   *RequestGoogleAction `json:"actions"`
	Type      string               `json:"type,omitempty"`
}
