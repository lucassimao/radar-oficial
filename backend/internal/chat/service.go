package chat

type ChatService struct{}

func NewChatService() *ChatService {
	return &ChatService{}
}

func (srv *ChatService) WelcomeMessage() (string, error) {

	return "", nil
}
