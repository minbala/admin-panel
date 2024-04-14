package resources

type ContextKey string

var (
	LogData  ContextKey = "log data"
	Data     ContextKey = "data"
	UserData ContextKey = "user data"
	PayLoad  ContextKey = "payLoad"
)
