package usecase

import (
	"context"
	"errors"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/imanudd/forum-app/config"
	"github.com/imanudd/forum-app/internal/domain"
	"github.com/imanudd/forum-app/internal/repository"
	"github.com/imanudd/forum-app/pkg/validator"
)

type AuthorUseCaseImpl interface {
	CreateAuthorAndBook(ctx context.Context, req *domain.CreateAuthorAndBookRequest) error
	DeleteBookByAuthor(ctx context.Context, id, bookId int) error
	GetListBookByAuthor(ctx context.Context, id int) ([]*domain.Book, error)
	AddAuthorBook(ctx context.Context, req *domain.AddAuthorBookRequest) error
	CreateAuthor(ctx context.Context, req *domain.CreateAuthorRequest) error
}

type authorUseCase struct {
	config     *config.MainConfig
	trx        repository.TransactionRepositoryImpl
	authorRepo repository.AuthorRepositoryImpl
	bookRepo   repository.BookRepositoryImpl
}

func NewAuthorUseCase(config *config.MainConfig, trx repository.TransactionRepositoryImpl, authorRepo repository.AuthorRepositoryImpl, bookRepo repository.BookRepositoryImpl) AuthorUseCaseImpl {
	return &authorUseCase{
		config:     config,
		trx:        trx,
		authorRepo: authorRepo,
		bookRepo:   bookRepo,
	}
}

func (u *authorUseCase) CreateAuthorAndBook(ctx context.Context, req *domain.CreateAuthorAndBookRequest) error {
	if err := validator.ValidateStruct(req); err != nil {
		return err
	}

	return u.trx.WithTransaction(ctx, func(txCtx context.Context) error {
		err := u.authorRepo.Create(txCtx, &domain.Author{
			ID:          req.Book.AuthorID,
			Name:        req.Author.Name,
			Email:       req.Author.Email,
			PhoneNumber: req.Author.PhoneNumber,
		})
		if err != nil {
			return err
		}

		return u.bookRepo.Create(txCtx, &domain.Book{
			AuthorID:  req.Book.AuthorID,
			BookName:  req.Book.BookName,
			Title:     req.Book.Title,
			Price:     req.Book.Price,
			CreatedAt: time.Now(),
		})
	})
}

func (u *authorUseCase) AddAuthorBook(ctx context.Context, req *domain.AddAuthorBookRequest) error {
	if err := validator.ValidateStruct(req); err != nil {
		return err
	}

	author, err := u.authorRepo.GetByID(ctx, req.AuthorID)
	if err != nil {
		return err
	}

	if author == nil {
		return errors.New("author not found")
	}

	return u.bookRepo.Create(ctx, &domain.Book{
		AuthorID: author.ID,
		BookName: req.BookName,
		Title:    req.Title,
		Price:    req.Price,
	})
}

func (u *authorUseCase) CreateAuthor(ctx context.Context, req *domain.CreateAuthorRequest) error {
	if err := validator.ValidateStruct(req); err != nil {
		return err
	}

	author, err := u.authorRepo.GetByName(ctx, req.Name)
	if err != nil {
		return err
	}

	if author != nil {
		return errors.New("author already exist")
	}

	return u.authorRepo.Create(ctx, &domain.Author{
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	})
}

func (u *authorUseCase) DeleteBookByAuthor(ctx context.Context, id int, bookId int) error {
	g, gCtx := errgroup.WithContext(ctx)

	var (
		book   *domain.Book
		author *domain.Author
		err    error
	)

	g.Go(func() error {
		book, err = u.bookRepo.GetByID(gCtx, bookId)
		if err != nil {
			return err
		}

		if book == nil {
			return errors.New("book not found")
		}

		return nil
	})

	g.Go(func() error {
		author, err = u.authorRepo.GetByID(gCtx, id)
		if err != nil {
			return err
		}

		if author == nil {
			return errors.New("author not found")
		}

		return nil

	})

	if err = g.Wait(); err != nil {
		return err
	}

	return u.bookRepo.DeleteBookByAuthorID(ctx, author.ID, book.ID)
}

func (u *authorUseCase) GetListBookByAuthor(ctx context.Context, id int) ([]*domain.Book, error) {
	author, err := u.authorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if author == nil {
		return nil, errors.New("author not found")
	}

	books, err := u.bookRepo.GetListBookByAuthorID(ctx, author.ID)
	if err != nil {
		return nil, err
	}

	return books, nil
}
