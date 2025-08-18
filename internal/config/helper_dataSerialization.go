package config

import (
	"time"

	"github.com/google/uuid"
	"github.com/rickNoise/chirpy/internal/database"
)

/*
Data Serialization

This project uses sqlc to generate Go models and functions for handling database interactions.
Because these models and functions are generated, we cannot have easy control over every aspect of their form.

We need control over structs in order to manage json struct tags, which are necessary for marshaling JSON when sending http responses.
Additionally, we do not always want to all DB fields (e.g. the "hashed_password" field from users), and managing our own data serialization will allow control over which fields are publicly available.
*/

/* USERS */

// A user struct with appropriate struct tags for public API responses.
// Note that the hashed password is purposely excluded.
type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

// Returns a user struct appropriate for public API responses (e.g. no hashed password included) (including json struct tags)
func DatabaseUserToAPIUser(u database.User) User {
	return User{
		ID:          u.ID,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		Email:       u.Email,
		IsChirpyRed: u.IsChirpyRed,
	}
}

/* CHIRPS */

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

// Returns a chirp struct appropriate for public API responses (including json struct tags)
func DatabaseChirpToAPIChirp(c database.Chirp) Chirp {
	return Chirp{
		Id:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserId:    c.UserID,
	}
}
