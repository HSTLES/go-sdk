package core_models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/hstles/go-sdk/core_datastore"
	"github.com/hstles/go-sdk/shared_helpers"
)

type User struct {
	HstlesUserID string
	Name         string
	Email        string
}

var ErrUserNotFound = errors.New("user not found")

// GetUserByID retrieves a user from CoreDB based on their hstles_user_id.
func GetUserByID(mgr *core_datastore.Manager, userID string) (*User, error) {
	log.Printf("Entering GetUserByID with userID: %s", userID)

	row := mgr.CoreDB.QueryRow(
		`SELECT hstles_user_id, name, email
           FROM users
          WHERE hstles_user_id = ?`,
		userID,
	)

	user := &User{}
	if err := row.Scan(&user.HstlesUserID, &user.Name, &user.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("query GetUserByID: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user from CoreDB based on their email address.
func GetUserByEmail(mgr *core_datastore.Manager, email string) (*User, error) {
	log.Printf("Entering GetUserByEmail with email: %s", email)

	row := mgr.CoreDB.QueryRow(
		`SELECT hstles_user_id, name, email
           FROM users
          WHERE email = ?`,
		email,
	)

	user := &User{}
	if err := row.Scan(&user.HstlesUserID, &user.Name, &user.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("query GetUserByEmail: %w", err)
	}

	return user, nil
}

// CreateUser creates a new user in CoreDB and sends a welcome email.
func CreateUser(mgr *core_datastore.Manager, name, email string) (*User, error) {
	log.Printf("Entering CreateUser with name: %s, email: %s", name, email)

	// 1) Generate a new condensed UUID for the user
	hstlesUserID, err := shared_helpers.GenerateCondensedUUID()
	if err != nil {
		return nil, fmt.Errorf("GenerateCondensedUUID: %w", err)
	}

	// 2) Insert the new user into the core users table
	_, err = mgr.CoreDB.Exec(
		`INSERT INTO users (hstles_user_id, name, email)
                 VALUES (?,             ?,    ?)`,
		hstlesUserID, name, email,
	)
	if err != nil {
		return nil, fmt.Errorf("insert new user: %w", err)
	}

	// 3) Build the User object to return
	newUser := &User{
		HstlesUserID: hstlesUserID,
		Name:         name,
		Email:        email,
	}

	// 4) Send a welcome email (failure to send is non‐fatal)
	// TODO: Implement email sending functionality
	// smtpCfg := config.LoadSMTPConfig()
	// mailer := services_email.NewMailer(smtpCfg)
	// if mailErr := mailer.SendWelcomeHTMLEmail(email, name); mailErr != nil {
	//     log.Printf("⚠️  warning: failed to send welcome email to %s: %v", email, mailErr)
	// }
	log.Printf("Welcome email sending not yet implemented for user: %s", email)

	return newUser, nil
}

// UpdateUser updates an existing user's information in the core database.
func UpdateUser(mgr *core_datastore.Manager, user User) error {
	log.Printf("Updating user: %s", user.HstlesUserID)

	query := `
		UPDATE users 
		SET name = ?, email = ?, updated_at = CURRENT_TIMESTAMP
		WHERE hstles_user_id = ?
	`
	_, err := mgr.CoreDB.Exec(query, user.Name, user.Email, user.HstlesUserID)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return fmt.Errorf("failed to update user: %w", err)
	}

	log.Printf("User successfully updated: %s", user.HstlesUserID)
	return nil
}

// DeleteUser deletes a user from the core database.
func DeleteUser(mgr *core_datastore.Manager, userID string) error {
	log.Printf("Deleting user: %s", userID)

	query := `
		DELETE FROM users 
		WHERE hstles_user_id = ?
	`
	_, err := mgr.CoreDB.Exec(query, userID)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	log.Printf("User successfully deleted: %s", userID)
	return nil
}
