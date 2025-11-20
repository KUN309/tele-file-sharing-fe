package api

import "time"

// ---------------- USER ----------------

type MeResponse struct {
	Data User `json:"data"`
}

type User struct {
	ID         int       `json:"id"`
	TelegramID int64     `json:"telegram_id"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"created_at"`
}

// ---------------- FILES ----------------

type FileItem struct {
	ID         int       `json:"id"`
	ObjectKey  string    `json:"object_key"`
	Filename   string    `json:"filename"`
	Size       int64     `json:"size"`
	Mime       string    `json:"mime"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UploadFileRequest struct {
	TelegramFileID string `json:"telegram_file_id"`
	Filename       string `json:"filename"`
}

type UploadFileResponse struct {
	FileID int    `json:"file_id"`
	Status string `json:"status"`
}

// ---------------- SHARES ----------------

type ShareCreateRequest struct {
	FileID          int      `json:"file_id"`
	FromTs          string   `json:"from_ts"`
	ToTs            string   `json:"to_ts"`
	RequirePassword bool     `json:"require_password"`
	Password        string   `json:"password"`
	Recipients      []string `json:"recipients"`
}

type ShareItem struct {
	ID             int        `json:"id"`
	FileID         int        `json:"file_id"`
	OwnerUserID    int        `json:"owner_user_id"`
	Hash           string     `json:"hash"`
	RequirePwd     bool       `json:"require_password"`
	HashPassword   *string    `json:"hash_password"`
	Revoked        bool       `json:"revoked"`
	ExpiresAt      *time.Time `json:"expires_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type RevokeResponse struct {
	Message string `json:"message"`
	ShareID int    `json:"share_id"`
	Status  string `json:"status"`
}

type PasswordAuthorizeRequest struct {
	Password string `json:"password"`
}

type ShareMetadata struct {
	ID        int    `json:"id"`
	FileID    int    `json:"file_id"`
	Filename  string `json:"filename"`
	// tùy backend return thêm gì thì bổ sung
}

// ---------------- PLACEHOLDER ----------------