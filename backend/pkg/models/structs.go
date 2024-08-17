package models

import "time"

const (
	// Privacy levels for posts
	PrivacyPublic        = "public"
	PrivacyPrivate       = "private"
	PrivacyAlmostPrivate = "almost_private"

	// Follow request statuses
	FollowRequestPending  = "pending"
	FollowRequestAccepted = "accepted"
	FollowRequestRejected = "rejected"
)

// RegisterRequest represents the registration form
type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
	Avatar      string `json:"avatar,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	AboutMe     string `json:"about_me,omitempty"`
}

// LoginRequest represents the login form
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
    ID          int       `json:"id"`
    Email       string    `json:"email"`
    FirstName   string    `json:"first_name"`
    LastName    string    `json:"last_name"`
	Password    string    `json:"password"`
    DateOfBirth string    `json:"date_of_birth"`
    Avatar      string    `json:"avatar"`
    Nickname    string    `json:"nickname"`
    AboutMe     string    `json:"about_me"`
    IsPrivate   bool      `json:"is_private"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}


// Post represents a post created by a user
type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	Image     string    `json:"image,omitempty"`
	Privacy   string    `json:"privacy"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Group struct {
	ID          int       `json:"id"`
	CreatorID   int       `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GroupEvent struct {
	ID          int       `json:"id"`
	GroupID     int       `json:"group_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DayTime     time.Time `json:"day_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GroupInvitation struct {
	ID          int       `json:"id"`
	GroupID     int       `json:"group_id"`
	InviterID   int       `json:"inviter_id"`
	InviteeID   int       `json:"invitee_id"`
	Status      string    `json:"status"` // pending, accepted, rejected
	InvitedAt   time.Time `json:"invited_at"`
	RespondedAt time.Time `json:"responded_at,omitempty"`
}

type GroupRequest struct {
	ID          int       `json:"id"`
	GroupID     int       `json:"group_id"`
	RequesterID int       `json:"requester_id"`
	Status      string    `json:"status"` // pending, accepted, rejected
	RequestedAt time.Time `json:"requested_at"`
	RespondedAt time.Time `json:"responded_at,omitempty"`
}

type EventRSVP struct {
	ID          int       `json:"id"`
	EventID     int       `json:"event_id"`
	UserID      int       `json:"user_id"`
	Status      string    `json:"status"` // going, not going
	RespondedAt time.Time `json:"responded_at"`
}

type GroupMembership struct {
	UserID   int        `json:"user_id"`
	GroupID  int        `json:"group_id"`
	JoinedAt time.Time  `json:"joined_at"`
	LeftAt   *time.Time `json:"left_at,omitempty"`
}

type Chat struct {
	ID          int       `json:"id"`
	SenderID    int       `json:"sender_id"`
	RecipientID int       `json:"recipient_id"`
	GroupID     int       `json:"group_id,omitempty"`
	Message     string    `json:"message"`
	IsGroup     bool      `json:"is_group"`
	CreatedAt   time.Time `json:"created_at"`
}

// Notification represents a notification for a user
type Notification struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

// FollowRequest represents a follow request between users
type FollowRequest struct {
	ID          int       `json:"id"`
	SenderID    int       `json:"sender_id"`
	RecipientID int       `json:"recipient_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// Like represents a like on a post
type Like struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PostID    int       `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Dislike represents a dislike on a post
type Dislike struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PostID    int       `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Comment represents a comment on a post
type Comment struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PostID    int       `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
