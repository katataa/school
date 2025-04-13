package graph

import (
	"context"
	"fmt"
	"match-me/config"
	"match-me/graph/model"
	"match-me/models"

	"github.com/gin-gonic/gin"
)

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	var user models.User
	if err := config.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	age := int32(user.Age)
	gqlBio := &model.Bio{
		ID:              fmt.Sprint(user.ID),
		Interests:       &user.Interests,
		Age:             &age,
		Gender:          &user.Gender,
		Location:        &user.Location,
		PreferredRadius: &user.PreferredRadius,
		Info:            &user.Info,

		User: &model.User{
			ID:    fmt.Sprint(user.ID),
			Name:  user.Name,
			Email: user.Email,
		},
	}

	var profile models.Profile
	profileErr := config.DB.First(&profile, "user_id = ?", user.ID).Error
	var gqlProfile *model.Profile
	if profileErr != nil {
		gqlProfile = &model.Profile{
			ID:   fmt.Sprint(user.ID),
			User: &model.User{ID: fmt.Sprint(user.ID), Name: user.Name},
		}
	} else {
		gqlProfile = &model.Profile{
			ID:   fmt.Sprint(profile.ID),
			User: &model.User{ID: fmt.Sprint(user.ID), Name: user.Name},
		}
	}

	return &model.User{
		ID:             fmt.Sprint(user.ID),
		Name:           user.Name,
		Email:          user.Email,
		ProfilePicture: &user.ProfilePicture,
		Bio:            gqlBio,
		Profile:        gqlProfile,
	}, nil
}

func (r *queryResolver) Bio(ctx context.Context, id string) (*model.Bio, error) {
	var bio models.Bio
	err := config.DB.First(&bio, "user_id = ?", id).Error
	if err != nil {
		var user models.User
		if err2 := config.DB.First(&user, "id = ?", id).Error; err2 != nil {
			return nil, fmt.Errorf("bio not found for user id %s: %v", id, err2)
		}
		age := int32(user.Age)
		return &model.Bio{
			ID:              fmt.Sprint(user.ID),
			Interests:       &user.Interests,
			Age:             &age,
			Gender:          &user.Gender,
			Location:        &user.Location,
			PreferredRadius: &user.PreferredRadius,
			Info:            &user.Info,
			User: &model.User{
				ID:    fmt.Sprint(user.ID),
				Name:  user.Name,
				Email: user.Email,
			},
		}, nil
	}

	var user models.User
	if err := config.DB.First(&user, "id = ?", bio.UserID).Error; err != nil {
		return nil, fmt.Errorf("user not found for bio: %v", err)
	}
	age := int32(bio.Age)
	return &model.Bio{
		ID:              fmt.Sprint(bio.ID),
		User:            &model.User{ID: fmt.Sprint(user.ID), Name: user.Name, Email: user.Email},
		Interests:       &bio.Interests,
		Age:             &age,
		Gender:          &bio.Gender,
		Location:        &bio.Location,
		PreferredRadius: &bio.PreferredRadius,
		Info:            &bio.Info,
	}, nil
}

func (r *queryResolver) Profile(ctx context.Context, id string) (*model.Profile, error) {
	var user models.User
	if err := config.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	var profile models.Profile
	if err := config.DB.First(&profile, "user_id = ?", id).Error; err != nil {

		return &model.Profile{
			ID:   fmt.Sprint(user.ID),
			User: &model.User{ID: fmt.Sprint(user.ID), Name: user.Name, Email: user.Email},
		}, nil
	}

	return &model.Profile{
		ID:   fmt.Sprint(profile.ID),
		User: &model.User{ID: fmt.Sprint(user.ID), Name: user.Name},
	}, nil
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	ginCtx, ok := ctx.Value("GinContextKey").(*gin.Context)
	if !ok {
		return nil, fmt.Errorf("unauthorized: missing Gin context")
	}

	userID, exists := ginCtx.Get("user_id")
	if !exists {
		return nil, fmt.Errorf("unauthorized: user_id not found")
	}

	return r.User(ctx, fmt.Sprint(userID))
}

func (r *queryResolver) MyBio(ctx context.Context) (*model.Bio, error) {
	ginCtx, ok := ctx.Value("GinContextKey").(*gin.Context)
	if !ok {
		return nil, fmt.Errorf("unauthorized: missing Gin context")
	}

	userID, exists := ginCtx.Get("user_id")
	if !exists {
		return nil, fmt.Errorf("unauthorized: user_id not found")
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	age := int32(user.Age)
	return &model.Bio{
		ID:              fmt.Sprint(user.ID),
		Interests:       &user.Interests,
		Age:             &age,
		Gender:          &user.Gender,
		Location:        &user.Location,
		PreferredRadius: &user.PreferredRadius,
		Info:            &user.Info,

		User: &model.User{
			ID:    fmt.Sprint(user.ID),
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

func (r *queryResolver) MyProfile(ctx context.Context) (*model.Profile, error) {
	ginCtx, ok := ctx.Value("GinContextKey").(*gin.Context)
	if !ok {
		return nil, fmt.Errorf("unauthorized: missing Gin context")
	}

	userID, exists := ginCtx.Get("user_id")
	if !exists {
		return nil, fmt.Errorf("unauthorized: user_id not found")
	}

	return r.Profile(ctx, fmt.Sprint(userID))
}

func (r *queryResolver) Recommendations(ctx context.Context) ([]*model.User, error) {
	var users []models.User
	if err := config.DB.Limit(104).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch recommendations: %v", err)
	}

	var gqlUsers []*model.User
	for _, user := range users {
		age := int32(user.Age)
		gqlBio := &model.Bio{
			ID:              fmt.Sprint(user.ID),
			Interests:       &user.Interests,
			Age:             &age,
			Gender:          &user.Gender,
			Location:        &user.Location,
			PreferredRadius: &user.PreferredRadius,
			Info:            &user.Info,
		}

		var profile models.Profile
		profileErr := config.DB.First(&profile, "user_id = ?", user.ID).Error
		var gqlProfile *model.Profile
		if profileErr != nil {
			gqlProfile = &model.Profile{
				ID:   fmt.Sprint(user.ID),
				User: &model.User{ID: fmt.Sprint(user.ID), Name: user.Name, Email: user.Email},
			}
		} else {
			gqlProfile = &model.Profile{
				ID:   fmt.Sprint(profile.ID),
				User: &model.User{ID: fmt.Sprint(user.ID), Name: user.Name, Email: user.Email},
			}
		}

		gqlUsers = append(gqlUsers, &model.User{
			ID:             fmt.Sprint(user.ID),
			Name:           user.Name,
			Email:          user.Email,
			ProfilePicture: &user.ProfilePicture,
			Bio:            gqlBio,
			Profile:        gqlProfile,
		})
	}
	return gqlUsers, nil
}

func (r *queryResolver) Connections(ctx context.Context) ([]*model.User, error) {
	ginCtx, ok := ctx.Value("GinContextKey").(*gin.Context)
	if !ok {
		return nil, fmt.Errorf("unauthorized: missing Gin context")
	}

	userIDInterface, exists := ginCtx.Get("user_id")
	if !exists {
		return nil, fmt.Errorf("unauthorized: user_id not found")
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		return nil, fmt.Errorf("unauthorized: user_id has invalid type")
	}

	var connections []models.Connection
	if err := config.DB.Where("(sender_id = ? OR receiver_id = ?) AND status = ?", userID, userID, "accepted").Find(&connections).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch connections: %v", err)
	}

	var userIDs []uint
	for _, conn := range connections {
		if conn.SenderID == userID {
			userIDs = append(userIDs, conn.ReceiverID)
		} else {
			userIDs = append(userIDs, conn.SenderID)
		}
	}

	var users []models.User
	if err := config.DB.Where("id IN (?)", userIDs).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch user connections: %v", err)
	}

	var gqlUsers []*model.User
	for _, user := range users {
		age := int32(user.Age)
		gqlBio := &model.Bio{
			ID:              fmt.Sprint(user.ID),
			Interests:       &user.Interests,
			Age:             &age,
			Gender:          &user.Gender,
			Location:        &user.Location,
			PreferredRadius: &user.PreferredRadius,
			Info:            &user.Info,
		}

		gqlUsers = append(gqlUsers, &model.User{
			ID:             fmt.Sprint(user.ID),
			Name:           user.Name,
			Email:          user.Email,
			ProfilePicture: &user.ProfilePicture,
			Bio:            gqlBio,
		})
	}
	return gqlUsers, nil
}

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
