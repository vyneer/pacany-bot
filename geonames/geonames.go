package geonames

import (
	"errors"
	"strings"
	"time"

	gn "github.com/mkrou/geonames"
	"github.com/mkrou/geonames/models"
)

var CityToTimezone = cityToTimezone{
	associationMap: map[string]string{},
}

var ErrNoSuchTimezone = errors.New("no timezone associated with this city")

type cityToTimezone struct {
	associationMap map[string]string
}

func New() error {
	c := cityToTimezone{
		associationMap: map[string]string{},
	}

	p := gn.NewParser()

	if err := p.GetGeonames(gn.Cities15000, func(g *models.Geoname) error {
		if _, exists := c.associationMap[g.AsciiName]; !exists {
			c.associationMap[strings.ToLower(g.AsciiName)] = g.Timezone
		}

		return nil
	}); err != nil {
		return err
	}

	CityToTimezone = c

	return nil
}

func (c *cityToTimezone) Get(cityName string) (*time.Location, error) {
	timezone, ok := c.associationMap[strings.ToLower(cityName)]
	if !ok {
		return nil, ErrNoSuchTimezone
	}

	return time.LoadLocation(timezone)
}
