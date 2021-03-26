package main

// QuoteData for quote Request
type QuoteData struct {
	PFP        string `json:"pfp"`
	NoPFP      string `json:"no_pfp"`
	Username   string `json:"username"`
	RawText    string `json:"raw_text"`
	Colour     string `json:"colour"`
	AdminTitle string `json:"admintitle"`
	Style      string `json:"style"`
	APIToken   string `json:"token"`
}

// Success aye
type Success struct {
	File  string
	Token string
}

// QuoteAnswer for parsing quote URL
type QuoteAnswer struct {
	Success         Success
	TokenInvalid    string `json:"202"`
	InvalidTemplate string `json:"402"`
	NoText          string `json:"401"`
	AccessDenied    string `json:"707"`
}
