package models

import (
	"encoding/json"
	"fmt"
)

type Status string

const (
	ONGOING   Status = "ongoing"
	COMPLETED Status = "completed"
)

// cobra Custom Value method
func (s *Status) String() string {
	return string(*s)
}

// cobra Custom Value method
func (s *Status) Type() string {
	return "string"
}

// cobra Custom Value method
func (s *Status) Set(val string) error {
	status := Status(val)
	if !status.IsValid() {
		return fmt.Errorf("Invalid status type (%s)", val)
	}
	*s = status
	return nil
}

func (s Status) IsValid() bool {
	switch s {
	case ONGOING, COMPLETED:
		return true
	}
	return false
}

func (s Status) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("%v", s)
	return json.Marshal(str)
}

func (s *Status) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	status := Status(v)
	if !status.IsValid() {
		return fmt.Errorf("Invalid status type (%v)", status)
	}

	*s = status

	return nil
}

type Series struct {
	Id     int64    `json:"id"`
	Title  string   `json:"title"`
	Status Status   `json:"status"`
	Isbns  []string `json:"isbns"`
}

type NewSeries struct {
	Title  string
	Status Status
	Isbns  []string
}

func (s NewSeries) CanRegister() bool {
	return s.Title != "" && s.Status.IsValid()
}
