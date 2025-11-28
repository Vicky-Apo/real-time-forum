package models

type HomePageData struct {
	Posts      []*Post        `json:"posts"`
	Categories []Category     `json:"categories"`
	Pagination PaginationInfo `json:"pagination"`
	User       *User          `json:"user,omitempty"`
	Sort       string         `json:"sort,omitempty"`
}

// RegisterPageData - Data for registration page template
type RegisterPageData struct {
	Error    string            `json:"error,omitempty"`
	Success  string            `json:"success,omitempty"`
	FormData *UserRegistration `json:"form_data,omitempty"` // this is to keep the form data in case of validation errors
}

// LoginPageData - Data for login page template
type LoginPageData struct {
	Error    string     `json:"error,omitempty"`
	Success  string     `json:"success,omitempty"`
	FormData *UserLogin `json:"form_data,omitempty"` // this is to keep the form data in case of validation errors
}

type PostPageData struct {
	Post               *Post          `json:"post"`
	Comments           []Comment      `json:"comments"`
	CommentsPagination PaginationInfo `json:"comments_pagination"`
	User               *User          `json:"user,omitempty"`
	Sort               string         `json:"sort,omitempty"`
}

// ProfilePageData - Data for user profile page template
type ProfilePageData struct {
	Profile        *UserProfile            `json:"profile"`
	Posts          *PaginatedPostsResponse `json:"posts,omitempty"`
	LikedPosts     *PaginatedPostsResponse `json:"liked_posts,omitempty"`
	CommentedPosts *PaginatedPostsResponse `json:"commented_posts,omitempty"`
	User           *User                   `json:"user,omitempty"`
}
