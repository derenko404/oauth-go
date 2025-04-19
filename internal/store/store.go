package store

type Store struct {
	User    UserStore
	Session SessionStore
}
