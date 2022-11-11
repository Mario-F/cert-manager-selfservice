package config

var (
	active   bool
	username string
	password string
)

func GetBasicAuth() (bool, string, string) {
	return active, username, password
}

func SetBasicAuth(user string, pass string) {
	active = true
	username = user
	password = pass
}
