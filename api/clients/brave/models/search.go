package models

type SearchResponse struct {
	Query Query  `json:"query"`
	Mixed Mixed  `json:"mixed"`
	Type  string `json:"type"`
	Web   Web    `json:"web"`
}
type Query struct {
	Original             string `json:"original"`
	ShowStrictWarning    bool   `json:"show_strict_warning"`
	IsNavigational       bool   `json:"is_navigational"`
	IsNewsBreaking       bool   `json:"is_news_breaking"`
	SpellcheckOff        bool   `json:"spellcheck_off"`
	Country              string `json:"country"`
	BadResults           bool   `json:"bad_results"`
	ShouldFallback       bool   `json:"should_fallback"`
	PostalCode           string `json:"postal_code"`
	City                 string `json:"city"`
	HeaderCountry        string `json:"header_country"`
	MoreResultsAvailable bool   `json:"more_results_available"`
	State                string `json:"state"`
}
type Main struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
	All   bool   `json:"all"`
}
type Mixed struct {
	Type string `json:"type"`
	Main []Main `json:"main"`
	Top  []any  `json:"top"`
	Side []any  `json:"side"`
}
type Profile struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	LongName string `json:"long_name"`
	Img      string `json:"img"`
}
type MetaURL struct {
	Scheme   string `json:"scheme"`
	Netloc   string `json:"netloc"`
	Hostname string `json:"hostname"`
	Favicon  string `json:"favicon"`
	Path     string `json:"path"`
}
type Thumbnail struct {
	Src      string `json:"src"`
	Original string `json:"original"`
	Logo     bool   `json:"logo"`
}
type Results struct {
	Title          string    `json:"title"`
	URL            string    `json:"url"`
	IsSourceLocal  bool      `json:"is_source_local"`
	IsSourceBoth   bool      `json:"is_source_both"`
	Description    string    `json:"description"`
	Profile        Profile   `json:"profile"`
	Language       string    `json:"language"`
	FamilyFriendly bool      `json:"family_friendly"`
	Type           string    `json:"type"`
	Subtype        string    `json:"subtype"`
	IsLive         bool      `json:"is_live"`
	MetaURL        MetaURL   `json:"meta_url"`
	PageAge        string    `json:"page_age,omitempty"`
	Thumbnail      Thumbnail `json:"thumbnail,omitempty"`
	Age            string    `json:"age,omitempty"`
	ContentType    string    `json:"content_type,omitempty"`
}
type Web struct {
	Type           string    `json:"type"`
	Results        []Results `json:"results"`
	FamilyFriendly bool      `json:"family_friendly"`
}
