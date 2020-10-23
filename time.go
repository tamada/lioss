package lioss

import (
	"fmt"
	"time"
)

/*Time represents the time for marshaling a specific format.*/
type Time struct {
	time time.Time
}

/*Now returns now time.*/
func Now() *Time {
	return &Time{time.Now()}
}

func (t *Time) format() string {
	return t.time.Format("2006-01-02T15:04:05-07:00")
}

/*UnmarshalJSON is called on unmarshaling JSON.*/
func (t *Time) UnmarshalJSON(data []byte) error {
	var time, err = time.Parse("\"2006-01-02T15:04:05-07:00\"", string(data))
	t.time = time
	return err
}

/*MarshalJSON is called on marshaling JSON.*/
func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.format())), nil
}
