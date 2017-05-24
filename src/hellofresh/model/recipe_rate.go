package model

import (
	"time"
)

// RecipeRate recipe rate entity
type RecipeRate struct {
	// ID
	ID interface{} `json:"_id,omitempty" bson:"_id,omitempty"`
	// RecipeID Recipe ID
	RecipeID string
	// Rate Rated score from 1-5
	Rate int
	// User Who rated it => we can definitely do recommandation based on rate score and user name
	User string
	// Modified Last modified time
	Modified time.Time
}
