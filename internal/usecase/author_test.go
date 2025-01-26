package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/imanudd/forum-app/config"
	"github.com/imanudd/forum-app/internal/domain"
	repositoryMock "github.com/imanudd/forum-app/shared/mock/repository"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateAuthorAndBook(t *testing.T) {
	Convey("Test create author and book", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &config.MainConfig{}
		authorRepo := repositoryMock.NewMockAuthorRepositoryImpl(ctrl)
		bookRepo := repositoryMock.NewMockBookRepositoryImpl(ctrl)
		trx := repositoryMock.NewMockTransactionRepositoryImpl(ctrl)

		authorUseCase := NewAuthorUseCase(config, trx, authorRepo, bookRepo)

		var (
			ctx     = context.Background()
			errResp = errors.New("error")
			req     = &domain.CreateAuthorAndBookRequest{
				Author: domain.CreateAuthorRequest{
					Name:        "jamil",
					Email:       "jamil@mail.com",
					PhoneNumber: "082833",
				},
				Book: domain.CreateBookRequest{
					AuthorID:  1,
					BookName:  "book test",
					Title:     "testing",
					Price:     10000,
					CreatedAt: time.Now(),
				},
			}
		)

		Convey("resp err validator", func() {
			req.Author.Name = ""
			err := authorUseCase.CreateAuthorAndBook(ctx, req)
			So(err, ShouldNotBeNil)
		})

		Convey("transaction schema", func() {
			Convey("error when create author", func() {
				trx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, fn func(txCtx context.Context) error) {
					authorRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errResp)
					err := fn(ctx)
					So(err, ShouldNotBeNil)
				}).Return(errResp)
				err := authorUseCase.CreateAuthorAndBook(ctx, req)
				So(err, ShouldNotBeNil)
			})

			Convey("error when create book", func() {
				trx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, fn func(txCtx context.Context) error) {
					authorRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
					bookRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errResp)
					err := fn(ctx)
					So(err, ShouldNotBeNil)
				}).Return(errResp)
				err := authorUseCase.CreateAuthorAndBook(ctx, req)
				So(err, ShouldNotBeNil)
			})

			Convey("commit", func() {
				trx.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, fn func(txCtx context.Context) error) {
					authorRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
					bookRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
					err := fn(ctx)
					So(err, ShouldBeNil)
				}).Return(nil)
				err := authorUseCase.CreateAuthorAndBook(ctx, req)
				So(err, ShouldBeNil)
			})

		})
	})
}
func TestAddAuthorBook(t *testing.T) {
	Convey("Test add author book", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &config.MainConfig{}
		authorRepo := repositoryMock.NewMockAuthorRepositoryImpl(ctrl)
		bookRepo := repositoryMock.NewMockBookRepositoryImpl(ctrl)
		trx := repositoryMock.NewMockTransactionRepositoryImpl(ctrl)

		var (
			ctx = context.Background()
			req = &domain.AddAuthorBookRequest{
				AuthorID: 1,
				BookName: "petualangan sherina",
				Title:    "adventure",
				Price:    10000,
			}

			author = &domain.Author{
				ID:          1,
				Name:        "jamil",
				Email:       "jamil@mail.com",
				PhoneNumber: "098827392",
			}

			errResp = errors.New("error")
		)

		authorUseCase := NewAuthorUseCase(config, trx, authorRepo, bookRepo)

		Convey("resp err validator", func() {
			req.BookName = ""
			err := authorUseCase.AddAuthorBook(ctx, req)
			So(err, ShouldNotBeNil)
		})

		Convey("resp err when get author by id", func() {
			authorRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, errResp)
			err := authorUseCase.AddAuthorBook(ctx, req)
			So(err, ShouldNotBeNil)
		})

		Convey("resp err when author not found ", func() {
			authorRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, nil)
			err := authorUseCase.AddAuthorBook(ctx, req)
			So(err, ShouldNotBeNil)
		})

		Convey("resp err when create book by author ", func() {
			authorRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(author, nil)
			bookRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errResp)
			err := authorUseCase.AddAuthorBook(ctx, req)
			So(err, ShouldNotBeNil)
		})

		Convey("resp success add book by author ", func() {
			authorRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(author, nil)
			bookRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			err := authorUseCase.AddAuthorBook(ctx, req)
			So(err, ShouldBeNil)
		})

	})
}

func TestCreateAuthor(t *testing.T) {
	Convey("Test create author", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &config.MainConfig{}
		authorRepo := repositoryMock.NewMockAuthorRepositoryImpl(ctrl)
		bookRepo := repositoryMock.NewMockBookRepositoryImpl(ctrl)
		trx := repositoryMock.NewMockTransactionRepositoryImpl(ctrl)

		authorUseCase := NewAuthorUseCase(config, trx, authorRepo, bookRepo)

		var (
			ctx = context.Background()
			req = &domain.CreateAuthorRequest{
				Name:        "jamil",
				Email:       "jamil@mail.com",
				PhoneNumber: "08129872979372",
			}

			author = &domain.Author{
				Name:        "jamil",
				Email:       "jamil@mail.com",
				PhoneNumber: "087648638832",
			}

			errResp = errors.New("error")
		)
		Convey("resp err validator", func() {
			req.Name = ""
			err := authorUseCase.CreateAuthor(ctx, req)
			So(err, ShouldNotBeNil)
		})

		Convey("resp err when get author by name", func() {
			authorRepo.EXPECT().GetByName(gomock.Any(), gomock.Any()).Return(nil, errResp)
			err := authorUseCase.CreateAuthor(ctx, req)
			So(err, ShouldNotBeNil)
		})

		Convey("resp err author is already exist", func() {
			authorRepo.EXPECT().GetByName(gomock.Any(), gomock.Any()).Return(author, nil)
			err := authorUseCase.CreateAuthor(ctx, req)
			So(err, ShouldNotBeNil)
		})

		Convey("resp err when create author ", func() {
			authorRepo.EXPECT().GetByName(gomock.Any(), gomock.Any()).Return(nil, nil)
			authorRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errResp)
			err := authorUseCase.CreateAuthor(ctx, req)
			So(err, ShouldNotBeNil)
		})

		Convey("resp success create author ", func() {
			authorRepo.EXPECT().GetByName(gomock.Any(), gomock.Any()).Return(nil, nil)
			authorRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			err := authorUseCase.CreateAuthor(ctx, req)
			So(err, ShouldBeNil)
		})

	})
}

func TestDeleteBookByAuthor(t *testing.T) {
	Convey("Test delete book by author", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &config.MainConfig{}
		authorRepo := repositoryMock.NewMockAuthorRepositoryImpl(ctrl)
		bookRepo := repositoryMock.NewMockBookRepositoryImpl(ctrl)
		trx := repositoryMock.NewMockTransactionRepositoryImpl(ctrl)

		authorUseCase := NewAuthorUseCase(config, trx, authorRepo, bookRepo)

		var (
			ctx      = context.Background()
			errResp  = errors.New("error")
			authorID = 1
			bookID   = 123

			book = &domain.Book{
				ID:        bookID,
				AuthorID:  authorID,
				BookName:  "buku tulis",
				Title:     "anak anak",
				Price:     10000,
				CreatedAt: time.Now(),
			}

			author = &domain.Author{
				ID:          authorID,
				Name:        "jamil",
				Email:       "jamil@mail.com",
				PhoneNumber: "0884782629363",
			}
		)

		Convey("resp err when get book by id", func() {
			bookRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, errResp).AnyTimes()
			authorRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(author, nil).AnyTimes()
			err := authorUseCase.DeleteBookByAuthor(ctx, authorID, bookID)
			So(err, ShouldNotBeNil)
		})

		Convey("resp err when book doesnt exist", func() {
			bookRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
			authorRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(author, nil).AnyTimes()
			err := authorUseCase.DeleteBookByAuthor(ctx, authorID, bookID)
			So(err, ShouldNotBeNil)
		})

		Convey("resp err when get author by id", func() {
			bookRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(book, nil).AnyTimes()
			authorRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, errResp).AnyTimes()
			err := authorUseCase.DeleteBookByAuthor(ctx, authorID, bookID)
			So(err, ShouldNotBeNil)
		})

		Convey("resp err when author doesnt exist", func() {
			bookRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(book, nil).AnyTimes()
			authorRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
			err := authorUseCase.DeleteBookByAuthor(ctx, authorID, bookID)
			So(err, ShouldNotBeNil)
		})

		Convey("resp err when delete book by author", func() {
			bookRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(book, nil).AnyTimes()
			authorRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(author, nil).AnyTimes()
			bookRepo.EXPECT().DeleteBookByAuthorID(gomock.Any(), gomock.Any(), gomock.Any()).Return(errResp)
			err := authorUseCase.DeleteBookByAuthor(ctx, authorID, bookID)
			So(err, ShouldNotBeNil)
		})

		Convey("resp success delete book by author", func() {
			bookRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(book, nil).AnyTimes()
			authorRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(author, nil).AnyTimes()
			bookRepo.EXPECT().DeleteBookByAuthorID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			err := authorUseCase.DeleteBookByAuthor(ctx, authorID, bookID)
			So(err, ShouldBeNil)
		})
	})
}

func TestGetListBookByAuthor(t *testing.T) {
	Convey("Test get list book by author", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := &config.MainConfig{}
		authorRepo := repositoryMock.NewMockAuthorRepositoryImpl(ctrl)
		bookRepo := repositoryMock.NewMockBookRepositoryImpl(ctrl)
		trx := repositoryMock.NewMockTransactionRepositoryImpl(ctrl)

		authorUseCase := NewAuthorUseCase(config, trx, authorRepo, bookRepo)

		var (
			ctx      = context.Background()
			errResp  = errors.New("error")
			authorID = 1

			author = &domain.Author{
				ID:          authorID,
				Name:        "jamil",
				Email:       "jamil@mail",
				PhoneNumber: "092173820",
			}

			books = []*domain.Book{
				{
					ID:       1,
					AuthorID: authorID,
					BookName: "buku tulis",
					Title:    "anak anak",
					Price:    7000,
				},
				{
					ID:       2,
					AuthorID: authorID,
					BookName: "buku tulis gambar",
					Title:    "anak anak",
					Price:    7000,
				},
			}
		)

		Convey("resp err when get author by id", func() {
			authorRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, errResp)
			resp, err := authorUseCase.GetListBookByAuthor(ctx, authorID)
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("resp err when author doesnt exist", func() {
			authorRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, nil)
			resp, err := authorUseCase.GetListBookByAuthor(ctx, authorID)
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("resp err when get books by author", func() {
			authorRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(author, nil)
			bookRepo.EXPECT().GetListBookByAuthorID(gomock.Any(), gomock.Any()).Return(nil, errResp)
			resp, err := authorUseCase.GetListBookByAuthor(ctx, authorID)
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("resp success", func() {
			authorRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(author, nil)
			bookRepo.EXPECT().GetListBookByAuthorID(gomock.Any(), gomock.Any()).Return(books, nil)
			resp, err := authorUseCase.GetListBookByAuthor(ctx, authorID)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}
