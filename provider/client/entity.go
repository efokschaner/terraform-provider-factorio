package client

import "strconv"

type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type UnitNumber uint32

func (u UnitNumber) String() string {
	return strconv.FormatUint(uint64(u), 10)
}

func ParseUnitNumber(unitNumberStr string) (UnitNumber, error) {
	unitNumber, err := strconv.ParseUint(unitNumberStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return UnitNumber(unitNumber), nil
}

type Entity struct {
	UnitNumber UnitNumber `json:"unit_number"`
	Surface    string     `json:"surface"`
	Name       string     `json:"name"`
	Position   Position   `json:"position"`
	Direction  Direction  `json:"direction"`
	Force      string     `json:"force"`
}

// Looks up an entity by its unit_number
func (client *FactorioClient) EntityGet(unitNumber UnitNumber) (*Entity, error) {
	var entity Entity
	pEntity := &entity
	return pEntity, client.Read(
		"entity",
		map[string]interface{}{"unit_number": unitNumber},
		&pEntity)
}

// Corresponding to https://lua-api.factorio.com/latest/LuaSurface.html#LuaSurface.create_entity
type EntityCreateOptions struct {
	Surface                  string                 `json:"surface"` // eg. "nauvis"
	Name                     string                 `json:"name"`
	Position                 Position               `json:"position"`
	Direction                Direction              `json:"direction"`
	Force                    string                 `json:"force"` // eg. "player", "enemy", "neutral"
	EntitySpecificParameters map[string]interface{} `json:"entity_specific_parameters"`

	// Unimplemented
	/*
		target;
		source;
		fast_replace;
		player;
		spill;
		raise_built;
		create_build_effect_smoke;
		spawn_decorations;
		move_stuck_players;
		item;
	*/
}

func (client *FactorioClient) EntityCreate(opts *EntityCreateOptions) (*Entity, error) {
	var result Entity
	return &result, client.Create("entity", opts, &result)
}

// All params are optional
type EntityUpdateOptions struct {
	Direction *Direction `json:"direction,omitempty"`
	Force     *string    `json:"force,omitempty"`
}

func (client *FactorioClient) EntityUpdate(unitNumber UnitNumber, opts *EntityUpdateOptions) (*Entity, error) {
	var result Entity
	return &result, client.Update(
		"entity",
		strconv.FormatUint(uint64(unitNumber), 10),
		opts,
		&result)
}

func (client *FactorioClient) EntityDelete(unitNumber UnitNumber) error {
	return client.Delete("entity", strconv.FormatUint(uint64(unitNumber), 10))
}
