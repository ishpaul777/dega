package post

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/service/core/action/author"
	"github.com/factly/dega-server/service/core/model"
	factCheckModel "github.com/factly/dega-server/service/fact-check/model"
	"github.com/factly/dega-server/util"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

// details - Get post by id
// @Summary Show a post by id
// @Description Get post by ID
// @Tags Post
// @ID get-post-by-id
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Space header string true "Space ID"
// @Param post_id path string true "Post ID"
// @Success 200 {object} postData
// @Router /core/posts/{post_id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	sID, err := util.GetSpace(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.Unauthorized()))
		return
	}

	postID := chi.URLParam(r, "post_id")
	id, err := strconv.Atoi(postID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &postData{}
	result.Authors = make([]model.Author, 0)
	result.Claims = make([]factCheckModel.Claim, 0)

	postAuthors := []model.PostAuthor{}
	postClaims := []factCheckModel.PostClaim{}
	result.ID = uint(id)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	defer cancel()

	err = config.DB.Session(&gorm.Session{Context: ctx}).Model(&model.Post{}).Preload("Medium").Preload("Format").Preload("Tags").Preload("Categories").Where(&model.Post{
		SpaceID: uint(sID),
	}).First(&result.Post).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	if result.Format.Slug == "fact-check" {
		config.DB.Session(&gorm.Session{Context: ctx}).Model(&factCheckModel.PostClaim{}).Where(&factCheckModel.PostClaim{
			PostID: uint(id),
		}).Preload("Claim").Preload("Claim.Rating").Preload("Claim.Rating.Medium").Preload("Claim.Claimant").Preload("Claim.Claimant.Medium").Find(&postClaims)

		// appending all post claims
		for _, postClaim := range postClaims {
			result.Claims = append(result.Claims, postClaim.Claim)
		}
	}

	// fetch all authors
	config.DB.Session(&gorm.Session{Context: ctx}).Model(&model.PostAuthor{}).Where(&model.PostAuthor{
		PostID: uint(id),
	}).Find(&postAuthors)

	// Adding author
	authors, err := author.All(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}
	for _, postAuthor := range postAuthors {
		aID := fmt.Sprint(postAuthor.AuthorID)
		if author, found := authors[aID]; found {
			result.Authors = append(result.Authors, author)
		}
	}

	renderx.JSON(w, http.StatusOK, result)
}
