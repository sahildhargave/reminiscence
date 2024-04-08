package repository

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/lib/pq"

	"github.com/google/uuid"
	"github.com/sahildhargave/memories/account/model"
	"github.com/sahildhargave/memories/account/model/apperrors"
)

//PGUserRespository is data.respository implementation
//of service layer UserRepository

type PGUserRespository struct {
	DB *sqlx.DB
}

// Create implements model.UserRepository.
func (r *PGUserRespository) Create(ctx context.Context, u *model.User) error {

	query := "INSERT INTO users (email,password) VALUES ($1,$2) RETURNING *"

	if err := r.DB.Get(u, query, u.Email, u.Password); err != nil {
		//unique constraint

		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			log.Printf("Could not create a user with email: %v.Reson: %v\n", u.Email, err.Code.Name())
			return apperrors.NewConflict("email", u.Email)
		}

		log.Printf("Could not create auser with email: %v. Reason: %v\n", u.Email, err)
		return apperrors.NewInternal()
	}

	return nil
}

// Find Id feteched user by ID

func (r *PGUserRespository) FindByID(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	user := &model.User{}

	query := "SELECT *FROM users WHERE uid=$1"

	// check error as it could be somethings other than not found

	if err := r.DB.Get(user, query, uid); err != nil {
		return user, apperrors.NewNotFound("uid", uid.String())

	}

	return user, nil
}

//New NewUser
