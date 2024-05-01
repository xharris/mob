package component

type RoomType int

const (
	Empty RoomType = iota
	Start
	// move to next level
	Finish
	Fight
	// heal all allies
	Rest
	// free training for random % of allies
	Library
	// buy mods
	Shop
	// hire allies
	Recruit
	// fire an ally to pass or lose hp/gold
	Fire
	// fire a random ally to pass or lose hp/gold
	FireRandom
)

var goodRoomTypes []RoomType = []RoomType{Rest, Library, Shop, Recruit}

type Room struct {
	X, Y    int
	Enemies []NPC
	// appears as 'unknown' to player
	Unknown   bool
	Completed bool
	Type      RoomType
}

func (r *RoomType) IsGood() bool {
	if *r == Start {
		return true
	}
	for _, t := range goodRoomTypes {
		if t == *r {
			return true
		}
	}
	return false
}
