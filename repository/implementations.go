package repository

import (
	"context"
	"fmt"
)

func (r *Repository) CreateProfile(ctx context.Context, profile Profile) (output Profile, err error) {
	tx := r.Db.WithContext(ctx).Create(&profile)
	if tx.Error != nil {
		err = tx.Error
	}

	output = profile
	output.UserId = profile.UserId
	return
}

func (r *Repository) GetProfile(ctx context.Context, filter map[string]interface{}) (output []Profile, err error) {
	tx := r.Db.WithContext(ctx).Select("user_id, full_name, password, phone, status, created_at, updated_at")

	for k, v := range filter {
		tx = tx.Where(fmt.Sprintf("%s = ?", k), v)
	}

	find := tx.Find(&output)
	err = find.Error
	return
}

func (r *Repository) UpdateProfile(ctx context.Context, updatedBy map[string]interface{}, updatedData map[string]interface{}) error {
	tx := r.Db.Table("users")

	for k, v := range updatedBy {
		tx = tx.Where(fmt.Sprintf("%s = ?", k), v)
	}

	res := tx.WithContext(ctx).Updates(updatedData)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *Repository) GetLogin(ctx context.Context, filter map[string]interface{}) (output []LoginModel, err error) {
	tx := r.Db.WithContext(ctx).Select("login_id, user_id, ip, token, expires, requests, created_at, updated_at")

	for k, v := range filter {
		tx = tx.Where(fmt.Sprintf("%s = ?", k), v)
	}

	find := tx.Find(&output)
	err = find.Error
	return
}

func (r *Repository) InsertIntoLogin(ctx context.Context, login LoginModel) (output LoginModel, err error) {
	tx := r.Db.WithContext(ctx).Create(&login)
	if tx.Error != nil {
		err = tx.Error
	}

	output = login
	return
}

func (r *Repository) UpdateLogin(ctx context.Context, updatedBy map[string]interface{}, updatedData map[string]interface{}) error {
	tx := r.Db.Table("login")

	for k, v := range updatedBy {
		tx = tx.Where(fmt.Sprintf("%s = ?", k), v)
	}

	res := tx.WithContext(ctx).Updates(updatedData)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
