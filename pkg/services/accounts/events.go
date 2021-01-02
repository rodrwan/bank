package accounts

const (
	// CreatedAccountEvent ...
	CreatedAccountEvent = "created-account"
)

// IEvent ...
type IEvent interface {
	GetData() []byte
	GetName() string
}

// GetAccountEvent ...
type GetAccountEvent struct {
	name string
	data []byte
}

// NewGetAccountEvent ...
func NewGetAccountEvent(name string, data []byte) IEvent {
	return GetAccountEvent{
		name: name,
		data: data,
	}
}

// GetName ...
func (ga GetAccountEvent) GetName() string {
	return ga.name
}

// GetData ...
func (ga GetAccountEvent) GetData() []byte {
	return ga.data
}
