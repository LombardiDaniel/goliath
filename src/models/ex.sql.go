package models

import (
	"database/sql"
)

type Chunk struct {
	DocumentId int64
	PageNum    int
	Text       string
	PageOffset int
}

func (m *Chunk) Insert(db *sql.DB) error {
	_, err := db.Exec(`
		INSERT INTO chunks
		(document_id, page_num, text, page_offset, embedding)
		VALUES ($1, $2, $3, $4, $5);
	`,
		m.DocumentId,
		m.PageNum,
		m.Text,
		m.PageNum,
	)
	return err
}
