package chatroom

import apperrors "chatapi/internal/errors"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(name, creator string) (Summary, error) {
	exists, err := s.repo.Exists(name)
	if err != nil {
		return Summary{}, err
	}

	if exists {
		return Summary{}, apperrors.ErrAlreadyExists
	}

	chatroomID, createdAt, err := s.repo.Create(name, creator)
	if err != nil {
		return Summary{}, err
	}

	return Summary{
		ID:        chatroomID,
		Name:      name,
		CreatedAt: createdAt,
	}, nil
}

func (s *Service) List() ([]Summary, error) {
	return s.repo.List()
}

func (s *Service) Update(chatroomID int, username, newName string) (Summary, error) {
	chatroom, err := s.repo.FindByID(chatroomID)
	if err != nil {
		return Summary{}, apperrors.ErrChatroomNotFound
	}

	if chatroom.Creator != username {
		return Summary{}, apperrors.ErrForbidden
	}

	if err = s.repo.UpdateName(chatroomID, newName); err != nil {
		return Summary{}, err
	}

	return Summary{
		ID:        chatroomID,
		Name:      newName,
		CreatedAt: chatroom.CreatedAt,
	}, nil
}

func (s *Service) Delete(chatroomID int, username string) error {
	chatroom, err := s.repo.FindByID(chatroomID)
	if err != nil {
		return apperrors.ErrChatroomNotFound
	}

	if chatroom.Creator != username {
		return apperrors.ErrForbidden
	}

	return s.repo.Delete(chatroomID)
}
