// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	CreateProfile(ctx context.Context, profile Profile) (output Profile, err error)
	GetProfile(ctx context.Context, filter map[string]interface{}) (output []Profile, err error)
	UpdateProfile(ctx context.Context, updatedBy map[string]interface{}, updatedData map[string]interface{}) error
	GetLogin(ctx context.Context, filter map[string]interface{}) (output []LoginModel, err error)
	InsertIntoLogin(ctx context.Context, login LoginModel) (output LoginModel, err error)
	UpdateLogin(ctx context.Context, updatedBy map[string]interface{}, updatedData map[string]interface{}) error
}
