package seeder

import (
	"log"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
	"gorm.io/gorm"
)

func SeedSkills() {
	skillNames := []string{
		"Cleaning", "Cooking", "Laundry", "Ironing", "Dishwashing",
		"Babysitting", "Driving", "Grocery Shopping", "Pet Care", "Elderly Care",
		"First Aid", "CPR", "Swimming", "Tutoring", "Child Supervision",
		"Child Nutrition", "Bed Making", "Vacuuming", "Mopping", "Dusting",
		"Window Cleaning", "Bathroom Cleaning", "Kitchen Cleaning", "Table Setting", "Wardrobe Organization",
		"Shoe Polishing", "Garbage Disposal", "Basic Sewing", "Gardening", "Plant Watering",
		"Pool Maintenance", "Basic Repairs", "Home Organization", "Time Management", "Multi-tasking",
		"Attention to Detail", "Teamwork", "Independence", "Basic English", "Fluent English",
		"Arabic Speaking", "Basic Computer Skills", "Smartphone Usage", "GPS Navigation", "Map Reading",
		"Defensive Driving", "Manual Car Driving", "Automatic Car Driving", "Vehicle Maintenance", "Car Washing",
		"Errand Running", "Mail Sorting", "Event Setup", "Meal Planning", "Grocery Budgeting",
		"Healthy Cooking", "Special Diet Cooking", "Infant Care", "Toddler Care", "Teen Supervision",
		"Homework Assistance", "Storytelling", "Playtime Activities", "Educational Games", "Safety Procedures",
		"Emergency Handling", "Fire Safety", "Sanitation Knowledge", "Sterilizing Bottles", "Diaper Changing",
		"Feeding Babies", "Sleep Training", "Soothing Techniques", "Elder Mobility Support", "Medication Reminder",
		"Wheelchair Handling", "Companionship", "Mental Health Awareness", "Communication Skills", "Problem Solving",
		"Flexibility", "Punctuality", "Trustworthiness", "Loyalty", "Patience",
		"Respectfulness", "Cultural Sensitivity", "Confidentiality", "Basic Math", "Reading & Writing",
		"Home Appliance Use", "Microwave Use", "Vacuum Cleaner Use", "Washing Machine Use", "Oven Use",
		"Knife Safety", "Food Hygiene", "Inventory Management", "Pet Grooming", "Dog Walking",
	}

	for _, name := range skillNames {
		var existing models.Skill
		err := initializers.DB.
			Where("name = ?", name).
			First(&existing).Error

		if err == nil {
			continue // Already exists
		}

		if err != nil && err != gorm.ErrRecordNotFound {
			log.Printf("❌ DB error on skill %s: %v", name, err)
			continue
		}

		if err := initializers.DB.Create(&models.Skill{Name: name}).Error; err != nil {
			log.Printf("❌ Failed to insert skill %s: %v", name, err)
		} else {
			log.Printf("✅ Inserted skill: %s", name)
		}
	}
}
