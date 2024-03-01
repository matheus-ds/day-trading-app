package middleware

// TODO: add authentication middleware here
import (
	"Not sure how to connect to other files? like front to end etc/couldnt figure it out"
)

type JwtCustomClaims struct { //not sure if this has to match our user from front end (password and username) or if theis is deteached and a new thing linked to each person??
	Name string `json:"name"`
	ID   string `json:"id"`
	jwt.StandardClaims
}

type JwtCustomRefreshClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}


type User struct {
	ID       string
	Name     string
}

//not sure which env file to put in the  secret key and expiry time for both the access token and refresh token in our configuration files such as .env or config.json.
//env docker had "secret key here" so i added it after that and also added it to env samle just not sure 100% which one needs it or both??
type Env struct {
	AccessTokenExpiryHour  int
	RefreshTokenExpiryHour int
	AccessTokenSecret      string
	RefreshTokenSecret     string
}