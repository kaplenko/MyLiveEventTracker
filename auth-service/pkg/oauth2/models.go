package oauth2Models

type User struct {
	ID           string
	Login        string
	Email        string
	Provider     string
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64
}
