package handler

import (
	"golang-sql/model"

	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
)

type UserUpdateRequest struct {
	User struct {
		Username string `json:"username"`
		Email    string `json:"email" validate:"email"'`
		Password string `json:"password"`
		Bio      string `json:"bio"`
		Image    string `json:"image"`
	} `json:"user"`
}

func newUserUpdateRequest() *UserUpdateRequest {
	return new(UserUpdateRequest)
}

func (r *UserUpdateRequest) populate(u *model.User) {
	r.User.Username = u.Username
	r.User.Email = u.Email
	r.User.Password = u.Password
	if u.Bio != nil {
		r.User.Bio = *u.Bio
	}
	if u.Image != nil {
		r.User.Image = *u.Image
	}
}

func (r *UserUpdateRequest) bind(c echo.Context, u *model.User) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	u.Username = r.User.Username
	u.Email = r.User.Email
	if r.User.Password != u.Password {
		h, err := u.HashPassword(r.User.Password)
		if err != nil {
			u.Password = h
		}
	}
	u.Bio = &r.User.Bio
	u.Image = &r.User.Image
	return nil
}

type UserRegisterRequest struct {
	User struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
	} `json:"user"`
}

func (r *UserRegisterRequest) bind(c echo.Context, u *model.User) error {
	if err := c.Bind(r); err != nil {
		return nil
	}
	if err := c.Validate(r); err != nil {
		return nil
	}
	u.Username = r.User.Username
	u.Email = r.User.Email
	h, err := u.HashPassword(r.User.Password)
	if err != nil {
		return err
	}
	u.Password = h
	return nil
}

type UserLoginRequest struct {
	User struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	} `json:"user"`
}

func (r *UserLoginRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}

type ArticleCreateRequest struct {
	Article struct {
		Title       string   `json:"title" validate:"required"`
		Description string   `json:"description" validate:"required"`
		Body        string   `json:"body" validate:"required"`
		Tags        []string `json:"tagList, omitempty"`
	} `json:"article"`
}

func (r *ArticleCreateRequest) bind(c echo.Context, a *model.Article) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	a.Title = r.Article.Title
	a.Slug = slug.Make(r.Article.Title)
	a.Description = r.Article.Description
	a.Body = r.Article.Body
	if r.Article.Tags != nil {
		for _, t := range r.Article.Tags {
			a.Tags = append(a.Tags, model.Tag{Tag: t})
		}
	}
	return nil
}

type ArticleUpdateRequest struct {
	Article struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Body        string   `json:"body"`
		Tags        []string `json:"tagList"`
	} `json:"article"`
}

func (r *ArticleUpdateRequest) populate(a *model.Article) {
	r.Article.Title = a.Title
	r.Article.Description = a.Description
	r.Article.Body = a.Body
}

func (r *ArticleUpdateRequest) bind(c echo.Context, a *model.Article) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	a.Title = r.Article.Title
	a.Slug = slug.Make(a.Title)
	a.Description = r.Article.Description
	a.Body = r.Article.Body
	return nil
}

type CreateCommentRequest struct {
	Comment struct {
		Body string `json:"body" validate:"required"`
	} `json:"comment"`
}

//func (r *CreateCommentRequest) bind(c echo.Context, cm *model.Comment) error {
//	if err := c.Bind(r); err != nil {
//		return err
//	}
//	if err := c.Validate(r); err != nil {
//		return err
//	}
//	cm.Body = r.Comment.Body
//	cm.UserID = userIDFromToken(c)
//	return nil
//}
