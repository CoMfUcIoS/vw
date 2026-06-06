package model

type Item struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Type     int    `json:"type,omitempty"`
	Login    *Login `json:"login,omitempty"`
	Notes    string `json:"notes,omitempty"`
	FolderID string `json:"folderId,omitempty"`
}

type Login struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Totp     string `json:"totp,omitempty"`
	URIs     []URI  `json:"uris,omitempty"`
}

type URI struct {
	URI string `json:"uri,omitempty"`
}
