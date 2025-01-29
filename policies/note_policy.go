package policies

type NotePolicy interface {
	CanCreate(authenticatedUserId uint) bool
	CanUpdate(userID, authenticatedUserId uint) bool
}

type NotePolicyImpl struct {
}

func NewNotePolicy() NotePolicy {
	return &NotePolicyImpl{}
}

func (s *NotePolicyImpl) CanCreate(authenticatedUserId uint) bool {
	return true
}

func (s *NotePolicyImpl) CanView(noteUserId, authenticatedUserId uint) bool {
	return noteUserId == authenticatedUserId
}

func (s *NotePolicyImpl) CanUpdate(noteUserId, authenticatedUserId uint) bool {
	return noteUserId == authenticatedUserId
}
