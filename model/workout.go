package model

import "time"

type Workout struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	UserID       int       `json:"user_id"`
	CalendarName string    `json:"calendarName,omitempty"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
	Exercises    Exercises `db:"exercises" json:"exercises"`
}

type Exercise struct {
	Name      string   `db:"name" json:"name"`
	Muscles   []Muscle `db:"muscles" json:"muscles"`
	Sets      []Set    `db:"sets" json:"sets"`
	SuperSets []string `db:"superSets" json:"superSets"`
}

type Exercises []Exercise

type MuscleGroup string

var MuscleGroups = []MuscleGroup{Arms, Back, Chest, Core, Heart, Legs, Shoulders}

type Muscle struct {
	Name        string      `db:"name" json:"name"`
	MuscleGroup MuscleGroup `db:"muscleGroup" json:"muscleGroup"`
}

const (
	Arms      MuscleGroup = "Arms"
	Back      MuscleGroup = "Back"
	Chest     MuscleGroup = "Chest"
	Core      MuscleGroup = "Core"
	Heart     MuscleGroup = "Heart"
	Legs      MuscleGroup = "Legs"
	Shoulders MuscleGroup = "Shoulders"
)

type Set struct {
	Weight float32 `db:"weight" json:"weight"`
	Reps   int8    `db:"reps" json:"reps"`
}
