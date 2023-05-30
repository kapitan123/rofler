package domain

import "time"

type Factory struct {
}

type MessageProcessor struct {
}

func MustNewFactory() Factory {
	return Factory{}
}

func (f Factory) NewAvailableHour(hour time.Time) (*Hour, error) {
	if err := f.validateTime(hour); err != nil {
		return nil, err
	}

	return &Hour{
		hour:         hour,
		availability: Available,
	}, nil
}
