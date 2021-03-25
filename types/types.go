package types

type Location struct {
	Zip                  string
	URL                  string
	SlotDetails          []SlotDetail
	OpenTimeSlots        int
	OpenAppointmentSlots int
	Name                 string
	Longitude            float64
	Latitude             float64
	City                 string
}

type SlotDetail struct {
	OpenTimeslots        int
	OpenAppointmentSlots int
	Manufacturer         string
}

type RedisCredentials struct {
	URL      string
	User     string
	Password string
	Host     string
	Port     string
}
