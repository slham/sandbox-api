package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Workout struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	UserID       string    `json:"user_id"`
	CalendarName string    `json:"calendarName,omitempty"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
	Exercises    Exercises `json:"exercises,omitempty"`
}

type Exercise struct {
	Name      string   `json:"name"`
	Muscles   []Muscle `json:"muscles,omitempty"`
	Sets      []Set    `json:"sets,omitempty"`
	SuperSets []string `json:"superSets,omitempty"`
}

type Exercises []Exercise

type MuscleGroup string

var MuscleGroups = []MuscleGroup{Arms, Back, Chest, Core, Heart, Legs, Shoulders}

type Muscle struct {
	Name        string      `json:"name"`
	MuscleGroup MuscleGroup `json:"muscleGroup,omitempty"`
}

const (
	Arms      MuscleGroup = "arms"
	Back      MuscleGroup = "back"
	Chest     MuscleGroup = "chest"
	Core      MuscleGroup = "core"
	Heart     MuscleGroup = "heart"
	Legs      MuscleGroup = "legs"
	Shoulders MuscleGroup = "shoulders"
)

type Set struct {
	Weight float32 `json:"weight,omitempty"`
	Reps   int8    `json:"reps,omitempty"`
}

func (e Exercises) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e *Exercises) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &e)
}

func (e Exercise) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e *Exercise) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &e)
}

func (m Muscle) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *Muscle) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &m)
}

func (s Set) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *Set) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &s)
}
