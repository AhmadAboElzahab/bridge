package seeder

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/AhmadAboElzahab/bridge/internal/constants"
	"github.com/AhmadAboElzahab/bridge/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func SeedMaids(db *gorm.DB, languages []models.Language, skills []models.Skill) {
	const count = 1000
	var maids []models.Maid

	for i := 0; i < count; i++ {
		now := time.Now()

		maid := models.Maid{
			Email:                 fmt.Sprintf("maid%d_%s@example.com", i, randomHash(5)),
			PhoneNumber:           randomPhone(),
			WhatsappNumber:        randomPhone(),
			PasswordHash:          "hashed_password",
			IsEmailVerified:       randBool(),
			IsPhoneVerified:       randBool(),
			VisaStatus:            pickOptionValue(constants.VisaStatusOptions),
			FirstName:             randomFirstName(),
			LastName:              randomLastName(),
			Gender:                pickOptionValue(constants.GenderOptions),
			DateOfBirth:           randomDOB(),
			NationalityID:         uint(rand.Intn(100) + 1),
			CurrentLocation:       pickOptionValue(constants.CitiesOptions),
			ProfilePicture:        randomImageURL(),
			ProfilePictureHash:    randomHash(8),
			YearsOfExperience:     rand.Intn(11),
			Bio:                   randomBio(),
			ExpectedSalary:        float64(rand.Intn(4000) + 1000),
			Availability:          pickOptionValue(constants.AvailabilityOptions),
			PreferredWorkLocation: pickOptionValue(constants.CitiesOptions),
			VideoIntroduction:     randomVideoURL(),
			WorkPreferences:       pickOptionValue(constants.WorkPreferencesOptions),

			SubscriptionStatus:           pickOptionValue(constants.SubscriptionStatusOptions),
			SubscriptionPlan:             pickOptionValue(constants.SubscriptionPlanOptions),
			SubscriptionStartDate:        &now,
			SubscriptionEndDate:          ptrToTime(now.AddDate(0, 1, 0)),
			AutoRenew:                    randBool(),
			SubscriptionPaymentGatewayID: randomHash(12),
			SubscriptionLastPaymentDate:  &now,
			SubscriptionFailedAttempts:   rand.Intn(3),

			IsProfileVisible:            randBool(),
			ProfileCompletionPercentage: rand.Intn(101),
			NotifyOnInterest:            randBool(),
			NotifyOnMessage:             randBool(),
			LastLoginAt:                 &now,
			Status:                      pickOptionValue(constants.StatusOptions),
			Religion:                    pickOptionValue(constants.ReligionOptions),

			FirstAidCertificate:  randBool(),
			CareGiverCertificate: randBool(),
			TeachingCertificate:  randBool(),
		}

		maid.Languages = getRandomSubsetLanguages(languages)
		maid.Skills = getRandomSubsetSkills(skills)

		maids = append(maids, maid)
	}

	if err := db.Create(&maids).Error; err != nil {
		log.Fatalf("❌ Failed to seed maids: %v", err)
	}

	log.Printf("✅ Seeded %d maids successfully", count)
}

func randBool() bool {
	return rand.Intn(2) == 1
}

func randomPhone() string {
	return "+9715" + randomDigits(8)
}

func randomDigits(n int) string {
	digits := ""
	for i := 0; i < n; i++ {
		digits += fmt.Sprintf("%d", rand.Intn(10))
	}
	return digits
}

func randomHash(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randomDOB() time.Time {
	year := rand.Intn(25) + 1975
	month := time.Month(rand.Intn(12) + 1)
	day := rand.Intn(28) + 1
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func ptrToTime(t time.Time) *time.Time {
	return &t
}

func pickOptionValue(jsonData datatypes.JSON) string {
	var options []map[string]string
	if err := json.Unmarshal(jsonData, &options); err != nil || len(options) == 0 {
		return ""
	}
	return options[rand.Intn(len(options))]["value"]
}

func getRandomSubsetLanguages(langs []models.Language) []models.Language {
	if len(langs) == 0 {
		return nil
	}
	n := rand.Intn(3) + 1
	rand.Shuffle(len(langs), func(i, j int) { langs[i], langs[j] = langs[j], langs[i] })
	return langs[:min(n, len(langs))]
}

func getRandomSubsetSkills(skills []models.Skill) []models.Skill {
	if len(skills) == 0 {
		return nil
	}
	n := rand.Intn(5) + 1
	rand.Shuffle(len(skills), func(i, j int) { skills[i], skills[j] = skills[j], skills[i] })
	return skills[:min(n, len(skills))]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// --- Random Text Utilities ---
func randomFirstName() string {
	names := []string{"Fatima", "Aisha", "Sara", "Layla", "Zahra", "Noor", "Huda", "Amira"}
	return names[rand.Intn(len(names))]
}

func randomLastName() string {
	names := []string{"Khan", "Ahmed", "Hussein", "Mohammed", "Ali", "Zayed", "Omar", "Saeed"}
	return names[rand.Intn(len(names))]
}

func randomBio() string {
	bios := []string{
		"Experienced in childcare and housekeeping.",
		"Reliable and hardworking domestic helper.",
		"Fluent in English and Tagalog, excellent cooking skills.",
		"Over 5 years of experience with elderly care.",
		"Skilled in cleaning, laundry, and organizing.",
	}
	return bios[rand.Intn(len(bios))]
}

// --- Media Generators ---
func randomImageURL() string {
	id := rand.Intn(1000) + 1
	return fmt.Sprintf("https://picsum.photos/id/%d/200/200", id)
}

func randomVideoURL() string {
	videos := []string{
		"https://sample-videos.com/video123/mp4/720/big_buck_bunny_720p_1mb.mp4",
		"https://sample-videos.com/video123/mp4/720/big_buck_bunny_720p_2mb.mp4",
		"https://sample-videos.com/video123/mp4/720/big_buck_bunny_720p_5mb.mp4",
		"https://sample-videos.com/video123/mp4/720/big_buck_bunny_720p_10mb.mp4",
	}
	return videos[rand.Intn(len(videos))]
}
