package libs

type ContextKey int

const (
	RequestId ContextKey = iota
	IpAddress
	UserID
	UserEmail
)