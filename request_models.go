package main

import (
	"gorm.io/gorm"
	"strings"
)

// RegisterRequest is the model for requests sent to /auth/register.
type RegisterRequest struct {
	AcceptedTOSVersion int    `json:"acceptedTOSVersion"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	Email              string `json:"email"`
	Day                string `json:"day"`
	Month              string `json:"month"`
	Year               string `json:"year"`
	RecaptchaCode      string `json:"recaptchaCode"`
}

type UpdateUserRequest struct {
	AcceptedTOSVersion int      `json:"acceptedTOSVersion"`
	Bio                string   `json:"bio"`
	BioLinks           []string `json:"bioLinks"`
	Birthday           string   `json:"birthday"`
	CurrentPassword    string   `json:"currentPassword"`
	DisplayName        string   `json:"displayName"`
	Email              string   `json:"email"`
	Password           string   `json:"password"`
	Status             string   `json:"status"`
	StatusDescription  string   `json:"statusDescription"`
	Tags               []string `json:"tags"`
	Unsubscribe        bool     `json:"unsubscribe"`
	UserIcon           string   `json:"userIcon"`
}

func (r *UpdateUserRequest) EmailChecks(u *User) (bool, error) {
	if r.Email == "" {
		return false, nil
	}

	pwdMatch, err := u.CheckPassword(r.CurrentPassword)
	if !pwdMatch || err != nil {
		return false, invalidCredentialsErrorInUserUpdate
	}

	if DB.Model(&User{}).Where("email = ?", r.Email).Or("pending_email = ?", r.Email).Error != gorm.ErrRecordNotFound {
		return false, userWithEmailAlreadyExistsErrorInUserUpdate
	}

	u.PendingEmail = r.Email
	// TODO: Queue up verification email send
	return true, nil
}

func (r *UpdateUserRequest) StatusChecks(u *User) (bool, error) {
	var status UserStatus
	if r.Status == "" {
		return false, nil
	}

	switch strings.ToLower(r.Status) {
	case "join me":
		status = UserStatus(r.Status)
	case "active":
		status = UserStatus(r.Status)
	case "ask me":
		status = UserStatus(r.Status)
	case "busy":
		status = UserStatus(r.Status)
	case "offline":
		if !u.IsStaff() {
			return false, invalidStatusDescriptionErrorInUserUpdate
		}
		status = UserStatus(r.Status)
	default:
		return false, invalidUserStatusErrorInUserUpdate
	}

	u.Status = status
	return true, nil
}

func (r *UpdateUserRequest) StatusDescriptionChecks(u *User) (bool, error) {
	if r.StatusDescription == "" {
		return false, nil
	}

	if len(r.StatusDescription) > 32 {
		return false, invalidStatusDescriptionErrorInUserUpdate
	}

	u.StatusDescription = r.StatusDescription
	return true, nil
}

func (r *UpdateUserRequest) BioChecks(u *User) (bool, error) {
	if r.Bio == "" {
		return false, nil
	}

	if len(r.Bio) > 512 {
		return false, invalidBioErrorInUserUpdate
	}

	u.Bio = r.Bio
	return true, nil
}
