package cookieutils

type cookieStore interface {
	Cookie(name string) (string, error)
	SetCookie(name string, value string, maxAge int, path string, domain string, secure bool, httpOnly bool)
}

var (
	OneYear  = 365 * 24 * 60 * 60
	OneMonth = 31 * 24 * 60 * 60
	OneWeek  = 7 * 60 * 60
	OneDay   = 24 * 60 * 60
	OneHour  = 60 * 60
)

// set cookie if not exists
// returns existing cookie or new value
func SetIfNotExists(store cookieStore, name string, value string, maxAge int, path string, domain string, secure bool, httpOnly bool) string {
	cookie, err := store.Cookie(name)

	if err != nil {
		// cookie does not exists
		store.SetCookie(
			name,
			value,
			maxAge,
			path,
			domain,
			secure,
			httpOnly,
		)

		return value
	} else {
		return cookie
	}
}

func Delete(store cookieStore, name string) {
	store.SetCookie(name, "", -1, "/", "", true, true)
}
