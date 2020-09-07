package movements

// MovementType represents the different types of data in the movement table rows.
type MovementType int

const (
	Move   MovementType = iota // Move represents a ship movement.
	Notice                     // Notice represents a notice from the Harbour Master.
)

// Location holds the names for a location.
type Location struct {
	Abbreviation string // The abbreviation of the location.
	Name         string // The full name of the location.
}

// Movement represents one row in the table of movements.
// Although this is called Movement it can also hold data for a Move or a Notice.
type Movement struct {
	Type     MovementType // Move or Notice.
	Position string       // The position in the table, starts at 1.
	Time     string       // The time of the movement.
	Name     string       // The name of the ship or the text of the notice.
	From     Location     // Where the ship is moving from.
	To       Location     // Where the ship is moving to.
	Method   string       // Any special movement methods (e.g. TA for Tug Assisted)
	Remarks  string       // Remarks.
}
