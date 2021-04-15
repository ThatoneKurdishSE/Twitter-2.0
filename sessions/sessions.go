package sessions

import(
	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("p4l0"))