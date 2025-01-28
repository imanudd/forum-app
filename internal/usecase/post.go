package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/imanudd/forum-app/config"
	"github.com/imanudd/forum-app/internal/domain"
	"github.com/imanudd/forum-app/internal/repository"
	"github.com/imanudd/forum-app/pkg/auth"
	"github.com/imanudd/forum-app/pkg/validator"
)

type PostUseCaseImpl interface {
	GetListPost(ctx context.Context, req *domain.GetListPostRequest) (*domain.GetListPostResponse, error)
	UpsertUserActivity(ctx context.Context, req *domain.UpsertUserActivity) error
	CreatePost(ctx context.Context, req *domain.CreatePostRequest) error
	CreateCommentOnPost(ctx context.Context, req *domain.CreateCommentRequest) error
}

type PostUseCase struct {
	cfg              *config.Config
	postRepo         repository.PostRepositoryImpl
	commentRepo      repository.CommentRepositoryImpl
	userActivityRepo repository.UserActivityRepositoryImpl
}

func NewPostUseCase(cfg *config.Config, postRepo repository.PostRepositoryImpl, commentRepo repository.CommentRepositoryImpl, userActivityRepo repository.UserActivityRepositoryImpl) PostUseCaseImpl {
	return &PostUseCase{cfg: cfg, postRepo: postRepo, commentRepo: commentRepo, userActivityRepo: userActivityRepo}
}

func (u *PostUseCase) GetListPost(ctx context.Context, req *domain.GetListPostRequest) (*domain.GetListPostResponse, error) {
	datas := domain.GetListPostResponse{
		Pagination: &domain.Pagination{
			Limit: req.Limit,
			Page:  req.Page,
		},
	}

	totalData, err := u.postRepo.CountPost(ctx)
	if err != nil {
		return &datas, err
	}
	datas.Pagination.TotalData = totalData
	datas.Pagination.TotalPage = datas.Pagination.GetTotalPage()

	listPost, err := u.postRepo.GetListPost(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(listPost) == 0 {
		datas.ListPost = make([]*domain.ListPost, 0)
		return &datas, nil
	}

	datas.ListPost = append(datas.ListPost, listPost...)

	return &datas, nil

}

func (u *PostUseCase) UpsertUserActivity(ctx context.Context, req *domain.UpsertUserActivity) error {
	post, err := u.postRepo.GetById(ctx, req.PostId)
	if err != nil {
		return err
	}

	user := auth.GetUserContext(ctx)

	userActivity, err := u.userActivityRepo.GetByUserIdAndPostId(ctx, user.Id, post.ID)
	if err != nil {
		return err
	}

	if userActivity == nil {
		if !req.IsLiked {
			return errors.New("user activity is not exist, so it can't be unliked")
		}

		return u.userActivityRepo.CreateUserActivity(ctx, &domain.UserActivity{
			UserId:    user.Id,
			PostId:    post.ID,
			IsLiked:   req.IsLiked,
			CreatedAt: time.Now(),
			CreatedBy: user.Username,
		})
	}

	if req.IsLiked == userActivity.IsLiked {
		return nil
	}

	userActivity.UpdatedBy = &user.Username
	userActivity.UpdatedAt = time.Now()
	userActivity.IsLiked = req.IsLiked

	return u.userActivityRepo.UpdateUserActivity(ctx, userActivity)
}

func (u *PostUseCase) CreatePost(ctx context.Context, req *domain.CreatePostRequest) error {
	user := auth.GetUserContext(ctx)

	if err := validator.ValidateStruct(req); err != nil {
		return err
	}

	createPost := domain.Post{
		UserID:       user.Id,
		PostTitle:    req.PostTitle,
		PostContent:  req.PostContent,
		PostHashtags: req.PostHashtags,
		CreatedAt:    time.Now(),
		CreatedBy:    user.Username,
	}

	return u.postRepo.CreatePost(ctx, &createPost)
}

func (u *PostUseCase) CreateCommentOnPost(ctx context.Context, req *domain.CreateCommentRequest) error {
	if err := validator.ValidateStruct(req); err != nil {
		return err
	}

	user := auth.GetUserContext(ctx)

	post, err := u.postRepo.GetById(ctx, req.PostId)
	if err != nil {
		return err
	}

	return u.commentRepo.CreateComment(ctx, &domain.Comment{
		UserId:         user.Id,
		PostId:         post.ID,
		CommentContent: req.CommentContent,
		CreatedAt:      time.Now(),
		CreatedBy:      user.Username,
	})
}
