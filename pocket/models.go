package pocket

// MappedArticle represents an article
type MappedArticle struct {
	ID    string   `json:"id"`
	URL   string   `json:"url"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}
