package qhm

// MovementType is a type that represents the different types of data in the movement table rows.
type MovementType int

// Set constants for actual movement types.
const (
	Move   MovementType = iota // Move represents a ship movement.
	Notice                     // Notice represents a notice from the Harbour Master.
)

// Movement represents one row in the table of movements.
// Although this is called Movement it can hold data for a Move or a Notice.
type Movement struct {
	Type            MovementType `json:"type"`     // Move or Notice.
	Position        string       `json:"position"` // The position in the table, starts at 1.
	Time            string       `json:"time"`     // The time of the movement.
	Name            string       `json:"name"`     // The name of the ship or the text of the notice.
	From            Location     `json:"from"`     // Where the ship is moving from.
	To              Location     `json:"to"`       // Where the ship is moving to.
	Method          string       `json:"method"`   // Any special movement methods (e.g. TA for Tug Assisted)
	Remarks         string       `json:"remarks"`  // Remarks.
	ImageUrl        string       `json:"imageUrl"` // A URL to an image of the ship (if the Movement is of type Move).
	VesselFinderUrl string       `json:"infoUrl"`  // A URL to a vessel finder search for the ship.
}
