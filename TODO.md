[ ] - JWT implementation
[ ] - sessions implementation
[ ] - add 2fa

-> auth middleware -> validate each jwt, user_session encoded into token should not be soft deleted
-> valid -> next()
-> invalid -> return 403

-> /sign-in -> if user exists -> create session -> new tokens
-> if user does not exists -> create user -> create session -> new tokens
-> /refresh -> token valid and token_session_id is not soft deleted -> update -> new tokens
-> /sign-out -> soft delete user_sessions

-> /terminate-session -> soft delete session
