package models

import (
	"time"

	"gorm.io/gorm"
)

type Maid struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	Email          string `json:"email" gorm:"unique;not null"`
	PhoneNumber    string `json:"phone_number"`
	WhatsappNumber string `json:"whatsapp_number" gorm:"not null"`
	PasswordHash   string `json:"-" gorm:"not null"`

	IsEmailVerified bool `json:"is_email_verified"`
	IsPhoneVerified bool `json:"is_phone_verified"`

	VisaStatus  string    `json:"visa_status" gorm:"type:varchar(20)"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Gender      string    `json:"gender"`
	DateOfBirth time.Time `json:"date_of_birth"`

	NationalityID uint    `json:"nationality_id" gorm:"index"`
	Nationality   Country `json:"nationality" gorm:"foreignKey:NationalityID"`

	CurrentLocation    string `json:"current_location" gorm:"index"`
	ProfilePicture     string `json:"profile_picture"`
	ProfilePictureHash string `json:"profile_picture_hash"`

	YearsOfExperience     int     `json:"years_of_experience"`
	Bio                   string  `json:"bio" gorm:"type:varchar(450)"`
	ExpectedSalary        float64 `json:"expected_salary"`
	Availability          string  `json:"availability"`
	PreferredWorkLocation string  `json:"preferred_work_location"`
	VideoIntroduction     string  `json:"video_introduction"`

	WorkPreferences string `json:"work_preferences"`

	SubscriptionStatus           string     `json:"subscription_status"`
	SubscriptionPlan             string     `json:"subscription_plan"`
	SubscriptionStartDate        *time.Time `json:"subscription_start_date"`
	SubscriptionEndDate          *time.Time `json:"subscription_end_date"`
	AutoRenew                    bool       `json:"auto_renew"`
	SubscriptionPaymentGatewayID string     `json:"subscription_payment_gateway_id"`
	SubscriptionLastPaymentDate  *time.Time `json:"subscription_last_payment_date"`
	SubscriptionFailedAttempts   int        `json:"subscription_failed_attempts"`

	IsProfileVisible            bool           `json:"is_profile_visible"`
	ProfileCompletionPercentage int            `json:"profile_completion_percentage" gorm:"default:0"`
	NotifyOnInterest            bool           `json:"notify_on_interest"`
	NotifyOnMessage             bool           `json:"notify_on_message"`
	LastLoginAt                 *time.Time     `json:"last_login_at"`
	CreatedAt                   time.Time      `json:"created_at"`
	UpdatedAt                   time.Time      `json:"updated_at"`
	DeletedAt                   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`

	FirstAidCertificate  bool `json:"first_aid_certificate"`
	CareGiverCertificate bool `json:"care_giver_certificate"`
	TeachingCertificate  bool `json:"teaching_certificate"`

	Status   string `json:"status"`
	Religion string `json:"religion"`

	Languages []Language `json:"languages" gorm:"many2many:maid_languages;"`
	Skills    []Skill    `json:"skills" gorm:"many2many:maid_skills;"`
}
