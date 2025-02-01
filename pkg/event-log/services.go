package eventlog

type Service interface{}

type eventLogSvc struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &eventLogSvc{repo: r}
}
