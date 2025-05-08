package models

import (
	"time"

	"gorm.io/gorm"
)

type Maid struct {
	ID             uint   `gorm:"primaryKey"`
	Email          string `gorm:"unique;not null"`
	PhoneNumber    string
	WhatsappNumber string `gorm:"not null"`
	PasswordHash   string `gorm:"not null"`

	IsEmailVerified bool
	IsPhoneVerified bool

	VisaStatus  string `gorm:"type:varchar(20)"` // e.g. "canceled", "other"
	FirstName   string
	LastName    string
	Gender      string
	DateOfBirth time.Time

	NationalityID uint    `gorm:"index"`
	Nationality   Country `gorm:"foreignKey:NationalityID"`

	CurrentLocation    string `gorm:"index"`
	ProfilePicture     string
	ProfilePictureHash string

	YearsOfExperience     int
	Bio                   string `gorm:"type:varchar(450)"`
	ExpectedSalary        float64
	Availability          string
	PreferredWorkLocation string
	VideoIntroduction     string // URL

	WorkPreferences string // e.g. "live-in,full-time"

	SubscriptionStatus           string
	SubscriptionPlan             string
	SubscriptionStartDate        *time.Time
	SubscriptionEndDate          *time.Time
	AutoRenew                    bool
	SubscriptionPaymentGatewayID string
	SubscriptionLastPaymentDate  *time.Time
	SubscriptionFailedAttempts   int

	IsProfileVisible            bool
	ProfileCompletionPercentage int `gorm:"default:0"`
	NotifyOnInterest            bool
	NotifyOnMessage             bool

	LastLoginAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// Document flags
	FirstAidCertificate  bool
	CareGiverCertificate bool
	TeachingCertificate  bool

	Status   string
	Religion string

	// Relations
	Languages []Language `gorm:"many2many:maid_languages;"`
	Skills    []Skill    `gorm:"many2many:maid_skills;"`
}
