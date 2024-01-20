package variables

var HttpPortDefault int
var NotFound404 string
var LogRequests bool

func init() {
	HttpPortDefault = 8081
	NotFound404 = "<body><center><b>404<br>Not found!</b></center></body>"
	LogRequests = true
}
