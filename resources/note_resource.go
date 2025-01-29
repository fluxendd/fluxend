package resources

import "myapp/models"

type NoteResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NoteResource(note *models.Note) NoteResponse {
	return NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: note.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func NoteResourceCollection(notes []models.Note) []NoteResponse {
	resourceNotes := make([]NoteResponse, len(notes))
	for i, note := range notes {
		resourceNotes[i] = NoteResource(&note)
	}

	return resourceNotes
}
