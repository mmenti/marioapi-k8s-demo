package main

type Resume struct {
	Name           string `json:"name"`
	ContactPhone   string `json:"contact_phone"`
	ContactEmail   string `json:"contact_email"`
	ContactAddress string `json:"contact_address"`
	Summary        string `json:"summary"`
	Projects       Projects
	Experience     Experience
}

type Summary struct {
	Text string `json:"text"`
}

type SearchResult struct {
	Projects   Projects
	Experience Experience
}

type Projects []struct {
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	Summary string   `json:"summary"`
	Url     string   `json:"url"`
	Tags    []string `json:"tags"`
}

type Experience []struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	Dates    string   `json:"dates"`
	Location string   `json:"location"`
	JobTitle string   `json:"job_title"`
	Summary  string   `json:"summary"`
	Tags     []string `json:"tags"`
}

type Tag struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type Tags []*Tag

type Contact struct {
	Name           string `json:"name"`
	ContactPhone   string `json:"contact_phone"`
	ContactEmail   string `json:"contact_email"`
	ContactAddress string `json:"contact_address"`
}

type ApiError struct {
	ErrorCode int    `json:"error_code"`
	ErrorText string `json:"error_text"`
}

type ApiSuccess struct {
	SuccessCode int    `json:"success_code"`
	SuccessText string `json:"success_text"`
}

type ApiToken struct {
	Token []struct {
		UserName  string `json:"user_name"`
		UserToken string `json:"user_token"`
	} `json:"token"`
}
