package message

import (
	"chatapi/internal/chatroom"
	apperrors "chatapi/internal/errors"
)

type Service struct {
	repo         *Repository
	chatroomRepo *chatroom.Repository
}

func NewService(repo *Repository, chatroomRepo *chatroom.Repository) *Service {
	return &Service{repo: repo, chatroomRepo: chatroomRepo}
}

func (s *Service) Create(chatroomID int, creator, content string) (Message, error) {
	chatroom, err := s.chatroomRepo.FindByID(chatroomID)
	if err != nil {
		return Message{}, apperrors.ErrChatroomNotFound
	}

	messageID, err := s.repo.Create(chatroom.ID, creator, content)
	if err != nil {
		return Message{}, err
	}

	return Message{
		ID:      messageID,
		Creator: creator,
		Content: content,
	}, nil
}

func (s *Service) ListByChatroom(chatroomID int) ([]Message, error) {
	_, err := s.chatroomRepo.FindByID(chatroomID)
	if err != nil {
		return nil, apperrors.ErrChatroomNotFound
	}

	return s.repo.ListByChatroom(chatroomID)
}

func (s *Service) Update(messageID int, username, newContent string) (Message, error) {
	message, err := s.repo.FindByID(messageID)
	if err != nil {
		return Message{}, apperrors.ErrMessageNotFound
	}

	if message.Creator != username {
		return Message{}, apperrors.ErrForbidden
	}

	if err = s.repo.UpdateContent(messageID, newContent); err != nil {
		return Message{}, err
	}

	return Message{
		ID:      messageID,
		Creator: username,
		Content: newContent,
	}, nil
}

func (s *Service) Delete(chatroomID, messageID int, username string) error {
	chatroom, err := s.chatroomRepo.FindByID(chatroomID)
	if err != nil {
		return apperrors.ErrChatroomNotFound
	}

	message, err := s.repo.FindByID(messageID)
	if err != nil {
		return apperrors.ErrMessageNotFound
	}

	if message.Creator != username && chatroom.Creator != username {
		return apperrors.ErrForbidden
	}

	return s.repo.Delete(messageID)
}
