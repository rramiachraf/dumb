package data

type SearchResponse struct {
	Response struct {
		Sections sections
	}
}

type result struct {
	ArtistNames string `json:"artist_names"`
	Title       string
	Path        string
	Thumbnail   string `json:"song_art_image_thumbnail_url"`
}

type hits []struct {
	Result result
}

type sections []struct {
	Type string
	Hits hits
}

type SearchResults struct {
	Query    string
	Sections sections
}
