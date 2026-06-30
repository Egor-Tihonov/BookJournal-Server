package services

import (
	"book-journal/internal/models"
	"context"
	"net/http"
)

// LibraryService owns the global book catalog (books table). It resolves a book
// by looking in our own catalog first and falling back to openlibrary. It knows
// nothing about users or their shelves.
type LibraryService struct {
	librRepo LibraryRepository
	client   *http.Client
}

func NewLibraryService(repo LibraryRepository, client *http.Client) *LibraryService {
	return &LibraryService{librRepo: repo, client: client}
}

// SearchBook выполняет поиск книги по ISBN или названию, сначала в базе данных, затем через внешний API OpenLibrary
func (libraryService *LibraryService) SearchBook(ctx context.Context, query models.Book) (*models.Book, error) {
	book, err := libraryService.dbSearch(ctx, query)
	if err != nil {
		return nil, err // реальная ошибка БД
	}
	if book != nil {
		return book, nil // нашли в БД
	}

	return libraryService.openLibrarySearch(ctx, query) // фоллбэк на внешний API
}

// dbSearch выполняет поиск книги в базе данных по ISBN или названию
func (libraryService *LibraryService) dbSearch(ctx context.Context, book models.Book) (*models.Book, error) {
	if book.ISBN != "" {
		return libraryService.librRepo.FindByISBN(ctx, book.ISBN)
	}
	return libraryService.librRepo.FindByName(ctx, book.Name)
}

// openLibrarySearch выполняет поиск книги через внешний API OpenLibrary
func (libraryService *LibraryService) openLibrarySearch(ctx context.Context, book models.Book) (*models.Book, error) {
	// https://openlibrary.org/api/books?bibkeys=ISBN:9781250319180&format=json&jscmd=data
	// тут мы делаем HTTP-запрос и парсим ответ из JSON и добавляем книгу в бд(в библиотеку)
	return nil, nil
}
