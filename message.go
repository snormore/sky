package skyd

//------------------------------------------------------------------------------
//
// Constants
//
//------------------------------------------------------------------------------

const (
	ExecuteMessageType  = "execute"
	ShutdownMessageType = "shutdown"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// A message used for communicating between channels.
type Message struct {
	messageType string
	channel     chan interface{}
	function    func() (interface{}, error)
}

//------------------------------------------------------------------------------
//
// Constructors
//
//------------------------------------------------------------------------------

func NewExecuteMessage(f func() (interface{}, error)) *Message {
	return &Message{
		messageType: ExecuteMessageType,
		channel:     make(chan interface{}),
		function:    f,
	}
}

func NewShutdownMessage() *Message {
	return &Message{messageType: ShutdownMessageType, channel: make(chan interface{})}
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

// Executes the message function and returns it to the channel.
func (m *Message) execute() {
	if m.function != nil {
		ret, err := m.function()
		if err != nil {
			m.channel <- err
		} else {
			m.channel <- ret
		}
	} else {
		m.channel <- nil
	}
}

// Waits for the message to complete.
func (m *Message) wait() (interface{}, error) {
	ret := <-m.channel

	close(m.channel)
	m.channel = nil

	if err, ok := ret.(error); ok {
		return nil, err
	}
	return ret, nil
}
