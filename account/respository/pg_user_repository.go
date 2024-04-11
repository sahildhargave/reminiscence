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

type pGUserRespository struct {
	DB *sqlx.DB
}

// New User Repository is a factory for initializing User Repositories
func NewUserRepository(db *sqlx.DB) model.UserRepository {
	return &pGUserRespository{
		DB: db,
	}
}

// Create implements model.UserRepository.
func (r *pGUserRespository) Create(ctx context.Context, u *model.User) error {

	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *"

	if err := r.DB.GetContext(ctx, u, query, u.Email, u.Password); err != nil {
		//unique constraint

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

func (r *pGUserRespository) FindByID(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	user := &model.User{}

	query := "SELECT *FROM users WHERE uid=$1"

	// check error as it could be somethings other than not found

	if err := r.DB.GetContext(ctx,user, query, uid); err != nil {
		return user, apperrors.NewNotFound("uid", uid.String())

	}

	return user, nil
}

//New NewUser

// FindByEmail retrives user row by email address

func (r *pGUserRespository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}

	query := "SELECT * FROM users WHERE email=$1"

	if err := r.DB.GetContext(ctx, user, query, email); err != nil {
		log.Printf("Unable to get user with email address: %v. Err: %v\n", email, err)
		return user, apperrors.NewNotFound("email", email)
	}

	return user, nil
}

// update updates a user's properties

func (r *pGUserRespository) Update(ctx context.Context, u *model.User) error {
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

// update Image is used to separately update a user's image separate from
// other account details

func (r *pGUserRespository) updateImage(ctx context.Context, uid uuid.UUID, imageURL string) (*model.User, error) {
	query := `
		UPDATE users 
		SET image_url=$2
		WHERE uid=$1
		RETURNING *;
	`

	// instantiated to scan into ref using 'GetContext'

	u := &model.User{}

	err := r.DB.GetContext(ctx, u, query, uid, imageURL)

	if err != nil {
		log.Printf("Error updating image_url in database: %v\n", err)
		return nil, apperrors.NewInternal()
	}

	return u, nil
}
