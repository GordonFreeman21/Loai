package backend

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type Chat struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID        string    `json:"id"`
	ChatID    string    `json:"chat_id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type DB struct {
	Conn *sql.DB
}

func InitDB(path string) (*DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS chats (
		id TEXT PRIMARY KEY,
		title TEXT,
		created_at DATETIME
	);
	CREATE TABLE IF NOT EXISTS messages (
		id TEXT PRIMARY KEY,
		chat_id TEXT,
		role TEXT,
		content TEXT,
		created_at DATETIME,
		FOREIGN KEY(chat_id) REFERENCES chats(id) ON DELETE CASCADE
	);`

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return &DB{Conn: db}, nil
}

func (d *DB) CreateChat(id, title string) error {
	_, err := d.Conn.Exec("INSERT INTO chats (id, title, created_at) VALUES (?, ?, ?)", id, title, time.Now())
	return err
}

func (d *DB) GetChats() ([]Chat, error) {
	rows, err := d.Conn.Query("SELECT id, title, created_at FROM chats ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []Chat
	for rows.Next() {
		var c Chat
		if err := rows.Scan(&c.ID, &c.Title, &c.CreatedAt); err != nil {
			return nil, err
		}
		chats = append(chats, c)
	}
	return chats, nil
}

func (d *DB) SaveMessage(id, chatID, role, content string) error {
	_, err := d.Conn.Exec("INSERT INTO messages (id, chat_id, role, content, created_at) VALUES (?, ?, ?, ?, ?)",
		id, chatID, role, content, time.Now())
	return err
}

func (d *DB) GetMessages(chatID string) ([]Message, error) {
	rows, err := d.Conn.Query("SELECT id, chat_id, role, content, created_at FROM messages WHERE chat_id = ? ORDER BY created_at ASC", chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.ChatID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func (d *DB) DeleteChat(id string) error {
	_, err := d.Conn.Exec("DELETE FROM chats WHERE id = ?", id)
	return err
}
