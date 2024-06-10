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
	ArtistImage string `json:"image_url"`
	ArtistName  string `json:"name"`
	URL         string `json:"url"`
	AlbumImage  string `json:"cover_art_url"`
	AlbumName   string `json:"full_title"`
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
