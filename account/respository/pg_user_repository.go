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

type pGUserRepository struct {
	DB *sqlx.DB
}

// New User Repository is a factory for initializing User Repositories
func NewUserRepository(db *sqlx.DB) model.UserRepository {
	return &pGUserRepository{
		DB: db,
	}
}

// Create implements model.UserRepository.
func (r *pGUserRepository) Create(ctx context.Context, u *model.User) error {
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *"

	if err := r.DB.GetContext(ctx, u, query, u.Email, u.Password); err != nil {
		// check unique constraint
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, err.Code.Name())
			return apperrors.NewConflict("email", u.Email)
		}

		log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, err)
		return apperrors.NewInternal()
	}
	return nil
}

// Find Id feteched user by ID

func (r *pGUserRepository) FindByID(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	user := &model.User{}

	query := "SELECT * FROM users WHERE uid=$1"

	// we need to actually check errors as it could be something other than not found
	if err := r.DB.GetContext(ctx, user, query, uid); err != nil {
		return user, apperrors.NewNotFound("uid", uid.String())
	}

	return user, nil
}

// Find By EMAILID
func (r *pGUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}

	query := "SELECT * FROM users WHERE email=$1"

	if err := r.DB.GetContext(ctx, user, query, email); err != nil {
		log.Printf("Unable to get user with email address: %v. Err: %v\n", email, err)
		return user, apperrors.NewNotFound("email", email)
	}

	return user, nil
}

// update a user's properties
func (r *pGUserRepository) Update(ctx context.Context, u *model.User) error {
	query := `
	    UPDATE users
		SET name=:name, email=:email, website=:website
		WHERE uid=:uid
		RETURNING *;
	`

	nstmt, err := r.DB.PrepareNamedContext(ctx, query)

	if err != nil {
		log.Printf("Unable to prepare user update query: %v\n", err)
		return apperrors.NewInternal()
	}

	if err := nstmt.GetContext(ctx, u, u); err != nil {
		log.Printf("Failed to update details for user: %v\n", u)
		return apperrors.NewInternal()
	}

	return nil
}
