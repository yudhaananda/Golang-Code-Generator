package viewModel

type ViewModel struct {
	ProjectName      string              `json:"projectName"`
	Entity           map[string][]string `json:"entity"`
	Relation         []map[string]string `json:"relation"`
	Database         map[string]string   `json:"database"`
	IsUsingWebSocket bool                `json:"isUsingWebSocket"`
}
