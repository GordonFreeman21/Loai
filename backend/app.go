package backend

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type App struct {
	ctx    context.Context
	db     *DB
	ollama *OllamaClient
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	userHome, _ := os.UserHomeDir()
	appDir := filepath.Join(userHome, ".ollama-chat")
	os.MkdirAll(appDir, 0755)
	dbPath := filepath.Join(appDir, "chat.db")

	db, err := InitDB(dbPath)
	if err != nil {
		fmt.Printf("Failed to init DB: %v\n", err)
	}
	a.db = db
	a.ollama = NewOllamaClient("http://localhost:11434")
}

func (a *App) GetModels() ([]string, error) {
	return a.ollama.GetModels()
}

func (a *App) GetChats() ([]Chat, error) {
	return a.db.GetChats()
}

func (a *App) GetMessages(chatID string) ([]Message, error) {
	return a.db.GetMessages(chatID)
}

func (a *App) CreateChat(title string) (string, error) {
	id := uuid.New().String()
	err := a.db.CreateChat(id, title)
	return id, err
}

func (a *App) SaveMessage(chatID, role, content string) error {
	id := uuid.New().String()
	return a.db.SaveMessage(id, chatID, role, content)
}

func (a *App) SendMessage(chatID, model string, messages []Message) error {
	go func() {
		err := a.ollama.StreamChat(a.ctx, model, messages)
		if err != nil {
			fmt.Printf("Error streaming chat: %v\n", err)
		}
	}()
	return nil
}

func (a *App) PullModel(model string) error {
	go a.ollama.PullModel(a.ctx, model)
	return nil
}

func (a *App) DeleteChat(id string) error {
	return a.db.DeleteChat(id)
}
