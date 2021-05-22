package client

import "fmt"

type Direction int

const (
	North Direction = iota + 1
	Northeast
	East
	Southeast
	South
	Southwest
	West
	Northwest
)

var directionToString = map[Direction]string{
	North:     "north",
	Northeast: "northeast",
	East:      "east",
	Southeast: "southeast",
	South:     "south",
	Southwest: "southwest",
	West:      "west",
	Northwest: "northwest",
}

func invertMap(m *map[Direction]string) *map[string]Direction {
	result := make(map[string]Direction)
	for k, v := range *m {
		result[v] = k
	}
	return &result
}

var stringToDirection = invertMap(&directionToString)

func (d Direction) String() string {
	return directionToString[d]
}

func ParseDirection(direction string) (Direction, error) {
	d, found := (*stringToDirection)[direction]
	if found {
		return d, nil
	} else {
		return 0, fmt.Errorf("\"%s\" is not a Direction", direction)
	}
}

// Implement TextMarshaler
func (d Direction) MarshalText() (text []byte, err error) {
	text = []byte(d.String())
	return
}

// Implement TextUnmarshaler
func (d *Direction) UnmarshalText(text []byte) error {
	var err error
	*d, err = ParseDirection(string(text))
	return err
}
