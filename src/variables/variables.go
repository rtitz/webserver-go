package variables

var LogPath string
var HttpPort int
var NotFound404 string

func init() {
	HttpPort = 8081
	// LogPath = "webserver.log"
	NotFound404 = "<body><center><b>404<br>Not found!</b></center></body>"
}
