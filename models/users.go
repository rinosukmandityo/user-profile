package models

type User struct {
	ID           string `json:"ID" bson:"_id" msgpack:"_id" db:"ID"`
	Name         string `json:"Name" bson:"Name" msgpack:"Name" db:"Name"`
	Email        string `json:"Email" bson:"Email" msgpack:"Email" db:"Email"`
	Password     string `json:"Password" bson:"Password" msgpack:"Password" db:"Password"`
	Telephone    string `json:"Telephone" bson:"Telephone" msgpack:"Telephone" db:"Telephone"`
	Address      string `json:"Address" bson:"Address" msgpack:"Address" db:"Address"`
	IsActive     bool   `json:"IsActive" bson:"IsActive" msgpack:"IsActive" db:"IsActive"`
	IsGoogleAuth bool   `json:"IsGoogleAuth" bson:"IsGoogleAuth" msgpack:"IsGoogleAuth" db:"IsGoogleAuth"`
}

func (m *User) TableName() string {
	return "users"
}

func (m *User) GetMapFormat() map[string]interface{} {
	return map[string]interface{}{
		"ID":           m.ID,
		"Name":         m.Name,
		"Email":        m.Email,
		"Password":     m.Password,
		"Telephone":    m.Telephone,
		"Address":      m.Address,
		"IsActive":     m.IsActive,
		"IsGoogleAuth": m.IsGoogleAuth,
	}
}

func (m *User) SplitByField() []map[string]interface{} {
	return []map[string]interface{}{
		{"ID": m.ID},
		{"Name": m.Name},
		{"Email": m.Email},
		{"Password": m.Password},
		{"Telephone": m.Telephone},
		{"Address": m.Address},
		{"IsActive": m.IsActive},
		{"IsGoogleAuth": m.IsGoogleAuth},
	}
}

func FromMapToUser(data map[string]interface{}) User {
	return User{
		ID:           data["ID"].(string),
		Name:         data["Name"].(string),
		Email:        data["Email"].(string),
		Password:     data["Password"].(string),
		Telephone:    data["Telephone"].(string),
		Address:      data["Address"].(string),
		IsActive:     data["IsActive"].(bool),
		IsGoogleAuth: data["IsGoogleAuth"].(bool),
	}
}
