package seeder

import (
	"log"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
)

func SeedLanguages() {
	db := initializers.DB

	languages := []models.Language{
		{ISO639_3: "cmn", Name: "Mandarin Chinese"},
		{ISO639_3: "spa", Name: "Spanish"},
		{ISO639_3: "eng", Name: "English"},
		{ISO639_3: "hin", Name: "Hindi"},
		{ISO639_3: "ben", Name: "Bengali"},
		{ISO639_3: "por", Name: "Portuguese"},
		{ISO639_3: "rus", Name: "Russian"},
		{ISO639_3: "jpn", Name: "Japanese"},
		{ISO639_3: "yue", Name: "Cantonese"},
		{ISO639_3: "mar", Name: "Marathi"},
		{ISO639_3: "tel", Name: "Telugu"},
		{ISO639_3: "tur", Name: "Turkish"},
		{ISO639_3: "kor", Name: "Korean"},
		{ISO639_3: "vie", Name: "Vietnamese"},
		{ISO639_3: "tam", Name: "Tamil"},
		{ISO639_3: "urd", Name: "Urdu"},
		{ISO639_3: "jv", Name: "Javanese"},
		{ISO639_3: "ita", Name: "Italian"},
		{ISO639_3: "tha", Name: "Thai"},
		{ISO639_3: "guj", Name: "Gujarati"},
		{ISO639_3: "pol", Name: "Polish"},
		{ISO639_3: "ukr", Name: "Ukrainian"},
		{ISO639_3: "mal", Name: "Malayalam"},
		{ISO639_3: "ori", Name: "Odia"},
		{ISO639_3: "pan", Name: "Punjabi"},
		{ISO639_3: "sun", Name: "Sundanese"},
		{ISO639_3: "hau", Name: "Hausa"},
		{ISO639_3: "bur", Name: "Burmese"},
		{ISO639_3: "amh", Name: "Amharic"},
		{ISO639_3: "yor", Name: "Yoruba"},
		{ISO639_3: "fas", Name: "Persian"},
		{ISO639_3: "mal", Name: "Malay"},
		{ISO639_3: "ibo", Name: "Igbo"},
		{ISO639_3: "uzb", Name: "Uzbek"},
		{ISO639_3: "fil", Name: "Filipino"},
		{ISO639_3: "roh", Name: "Romansh"},
		{ISO639_3: "nld", Name: "Dutch"},
		{ISO639_3: "swe", Name: "Swedish"},
		{ISO639_3: "ron", Name: "Romanian"},
		{ISO639_3: "azb", Name: "South Azerbaijani"},
		{ISO639_3: "ceb", Name: "Cebuano"},
		{ISO639_3: "mya", Name: "Burmese"},
		{ISO639_3: "sin", Name: "Sinhala"},
		{ISO639_3: "khm", Name: "Khmer"},
		{ISO639_3: "nep", Name: "Nepali"},
		{ISO639_3: "zho", Name: "Chinese"},
		{ISO639_3: "ell", Name: "Greek"},
		{ISO639_3: "bul", Name: "Bulgarian"},
		{ISO639_3: "fin", Name: "Finnish"},
		{ISO639_3: "nor", Name: "Norwegian"},
		{ISO639_3: "dan", Name: "Danish"},
		{ISO639_3: "heb", Name: "Hebrew"},
		{ISO639_3: "slk", Name: "Slovak"},
		{ISO639_3: "hun", Name: "Hungarian"},
		{ISO639_3: "ces", Name: "Czech"},
		{ISO639_3: "srp", Name: "Serbian"},
		{ISO639_3: "hrv", Name: "Croatian"},
		{ISO639_3: "lit", Name: "Lithuanian"},
		{ISO639_3: "lav", Name: "Latvian"},
		{ISO639_3: "est", Name: "Estonian"},
		{ISO639_3: "slv", Name: "Slovenian"},
		{ISO639_3: "mkd", Name: "Macedonian"},
		{ISO639_3: "alb", Name: "Albanian"},
		{ISO639_3: "bos", Name: "Bosnian"},
		{ISO639_3: "isl", Name: "Icelandic"},
		{ISO639_3: "gle", Name: "Irish"},
		{ISO639_3: "mlt", Name: "Maltese"},
		{ISO639_3: "glg", Name: "Galician"},
		{ISO639_3: "cat", Name: "Catalan"},
		{ISO639_3: "eus", Name: "Basque"},
		{ISO639_3: "afr", Name: "Afrikaans"},
		{ISO639_3: "zul", Name: "Zulu"},
		{ISO639_3: "xho", Name: "Xhosa"},
		{ISO639_3: "sna", Name: "Shona"},
		{ISO639_3: "som", Name: "Somali"},
		{ISO639_3: "kin", Name: "Kinyarwanda"},
		{ISO639_3: "lug", Name: "Luganda"},
		{ISO639_3: "twi", Name: "Twi"},
		{ISO639_3: "wol", Name: "Wolof"},
		{ISO639_3: "orm", Name: "Oromo"},
		{ISO639_3: "swa", Name: "Swahili"},
		{ISO639_3: "lin", Name: "Lingala"},
		{ISO639_3: "bam", Name: "Bambara"},
		{ISO639_3: "ful", Name: "Fulah"},
		{ISO639_3: "snd", Name: "Sindhi"},
		{ISO639_3: "pus", Name: "Pashto"},
		{ISO639_3: "kur", Name: "Kurdish"},
		{ISO639_3: "kaz", Name: "Kazakh"},
		{ISO639_3: "aze", Name: "Azerbaijani"},
		{ISO639_3: "tgk", Name: "Tajik"},
		{ISO639_3: "kir", Name: "Kyrgyz"},
		{ISO639_3: "tuk", Name: "Turkmen"},
		{ISO639_3: "uig", Name: "Uyghur"},
		{ISO639_3: "mon", Name: "Mongolian"},
		{ISO639_3: "lao", Name: "Lao"},
		{ISO639_3: "hmn", Name: "Hmong"},
		{ISO639_3: "ace", Name: "Acehnese"},
		{ISO639_3: "mad", Name: "Madurese"},
		{ISO639_3: "min", Name: "Minangkabau"},
		{ISO639_3: "bug", Name: "Buginese"},
		{ISO639_3: "ban", Name: "Balinese"},
		{ISO639_3: "ace", Name: "Acehnese"},
		{ISO639_3: "bjn", Name: "Banjar"},
		{ISO639_3: "man", Name: "Mandingo"},
		{ISO639_3: "sus", Name: "Susu"},
		{ISO639_3: "kpe", Name: "Kpelle"},
		{ISO639_3: "dyu", Name: "Dyula"},
		{ISO639_3: "fon", Name: "Fon"},
		{ISO639_3: "ewe", Name: "Ewe"},
		{ISO639_3: "kan", Name: "Kannada"},
		{ISO639_3: "ar", Name: "Arabic"},
	}

	for _, lang := range languages {
		var existing models.Language
		if err := db.Where("iso639_3 = ?", lang.ISO639_3).First(&existing).Error; err == nil {
			log.Printf("ℹ️ Language already exists: %s", lang.Name)
			continue
		}

		if err := db.Create(&lang).Error; err != nil {
			log.Printf("❌ Failed to insert %s: %v", lang.Name, err)
		} else {
			log.Printf("✅ Inserted language: %s", lang.Name)
		}
	}
}
