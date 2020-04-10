package litterboxuser

// LitterboxUser = defines the attributes of a cat using the litterbox
type LitterboxUser struct {
	ID                   string
	Name                 string
	NameProbability      float64
	Direction            string
	DirectionProbability float64
	Photo                string
}
