package mapper

import (
	"github.com/goxgen/goxgen/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

var mapper, _ = New()

type Phone struct {
	ID     int    `mapto:"PhoneDto.ID"`
	Number string `mapto:"PhoneDto.Value"`
}

type PhoneDto struct {
	ID    int    `mapto:"Phone.ID"`
	Value string `mapto:"Phone.Number"`
}

type User struct {
	ID       int      `mapto:"UserDto.ID,Person.ID"`
	Name     string   `mapto:"UserDto.Name,Person.Name"`
	Phones   []*Phone `mapto:"UserDto.Phones"`
	Language *string  `mapto:"Person.Language"`
}

type UserDto struct {
	ID     int         `mapto:"User.ID"`
	Name   string      `mapto:"User.Name"`
	Phones []*PhoneDto `mapto:"User.Phones"`
}

type NestedPhoneInput struct {
	ID     *int         `mapto:"NestedPhone.ID" `
	Number *string      `mapto:"NestedPhone.Number"`
	Owner  *PersonInput `mapto:"NestedPhone.Owner"`
}

type NestedPhone struct {
	ID     int
	Number string
	Owner  *Person
}

type Person struct {
	ID       int
	Name     string
	Language string
	Role     *string
	Parents  []*Person
	Inviter  *Person
}

type PersonInput struct {
	ID      *int           `mapto:"Person.ID"`
	Name    *string        `mapto:"Person.Name"`
	Role    string         `mapto:"Person.Role"`
	Inviter *PersonInput   `mapto:"Person.Inviter"`
	Parents []*PersonInput `mapto:"Person.Parents"`
}

func TestMapByMapto(t *testing.T) {
	// Test basic struct to struct mapping
	t.Run("Basic mapping", func(t *testing.T) {
		phone := &Phone{ID: 1, Number: "1234567890"}
		phoneDto := &PhoneDto{}
		err := mapper.Map(phone, phoneDto)
		assert.Nil(t, err)
		assert.Equal(t, phone.ID, phoneDto.ID)
		assert.Equal(t, phone.Number, phoneDto.Value)
	})

	// Test slice mapping
	t.Run("Slice mapping", func(t *testing.T) {
		user := &User{
			ID:   1,
			Name: "John",
			Phones: []*Phone{
				{ID: 1, Number: "123"},
				{ID: 2, Number: "456"},
			},
		}
		userDto := &UserDto{}
		err := mapper.Map(user, userDto)
		assert.Nil(t, err)
		assert.Equal(t, user.ID, userDto.ID)
		assert.Equal(t, user.Name, userDto.Name)
		assert.Equal(t, len(user.Phones), len(userDto.Phones))
		for i := range user.Phones {
			assert.Equal(t, user.Phones[i].ID, userDto.Phones[i].ID)
			assert.Equal(t, user.Phones[i].Number, userDto.Phones[i].Value)
		}
	})

	// Test nested struct mapping
	t.Run("Nested struct mapping", func(t *testing.T) {
		user := &User{
			ID:   1,
			Name: "John",
			Phones: []*Phone{
				{ID: 1, Number: "123"},
			},
		}
		userDto := &UserDto{}
		err := mapper.Map(user, userDto)
		assert.Nil(t, err)
		assert.Equal(t, user.ID, userDto.ID)
		assert.Equal(t, user.Name, userDto.Name)
		assert.Equal(t, user.Phones[0].ID, userDto.Phones[0].ID)
		assert.Equal(t, user.Phones[0].Number, userDto.Phones[0].Value)
	})

	// Test pointer and non-pointer fields
	t.Run("Pointer and non-pointer fields", func(t *testing.T) {
		phone := &Phone{ID: 1, Number: "1234567890"}
		phoneDto := PhoneDto{}
		err := mapper.Map(phone, &phoneDto)
		assert.Nil(t, err)
		assert.Equal(t, phone.ID, phoneDto.ID)
		assert.Equal(t, phone.Number, phoneDto.Value)
	})

	t.Run("Non-settable fields", func(t *testing.T) {
		phone := Phone{ID: 1, Number: "1234567890"} // Non-pointer type
		phoneDto := &PhoneDto{}                     // Pointer type
		err := mapper.Map(phone, phoneDto)          // Should be able to map non-pointer to pointer
		assert.Nil(t, err)
		assert.Equal(t, phone.ID, phoneDto.ID)
		assert.Equal(t, phone.Number, phoneDto.Value)
	})

	t.Run("Nested pointers to structs", func(t *testing.T) {
		nestedInput := &NestedPhoneInput{
			ID:     utils.IntP(1),
			Number: utils.StringP("1234567890"),
			Owner: &PersonInput{
				ID:   utils.IntP(2),
				Name: utils.StringP("John"),
			},
		}
		nested := &NestedPhone{}

		err := mapper.Map(nestedInput, nested)
		assert.Nil(t, err)
		assert.Equal(t, *nestedInput.ID, nested.ID)
		assert.Equal(t, *nestedInput.Number, nested.Number)
		assert.NotNil(t, nested.Owner)
		assert.Equal(t, *nestedInput.Owner.ID, nested.Owner.ID)
		assert.Equal(t, *nestedInput.Owner.Name, nested.Owner.Name)
	})

	t.Run("Nil pointer and non-pointer fields", func(t *testing.T) {
		user := User{
			ID:       1,
			Name:     "User 1",
			Language: nil,
		}
		person := &Person{}             // Pointer type
		err := mapper.Map(user, person) // Should be able to map non-pointer to pointer
		assert.Nil(t, err)
		assert.Equal(t, user.ID, person.ID)
		assert.Equal(t, user.Name, person.Name)
	})

	t.Run("Non pointer destination", func(t *testing.T) {
		user := User{
			ID:       1,
			Name:     "User 1",
			Language: utils.StringP("EN"),
		}
		person := Person{}
		err := mapper.Map(user, person)
		assert.IsType(t, err, &DestinationNotPointerError{})
	})

	t.Run("Non struct source", func(t *testing.T) {
		person := Person{}
		err := mapper.Map(10, &person)
		assert.IsType(t, err, &SourceNotStructError{})
	})

	t.Run("Non struct dest", func(t *testing.T) {
		userDto := UserDto{}
		err := mapper.Map(userDto, utils.StringP("bad_dest"))
		assert.IsType(t, err, &DestinationNotStructError{})
	})

	t.Run("Non pointer to pointer", func(t *testing.T) {
		personInput := PersonInput{
			ID:   utils.IntP(1),
			Name: utils.StringP("John"),
			Role: "admin",
		}
		person := &Person{}
		err := mapper.Map(personInput, person)
		assert.Nil(t, err)
		assert.Equal(t, *personInput.ID, person.ID)
		assert.Equal(t, *personInput.Name, person.Name)
		assert.Equal(t, personInput.Role, *person.Role)
	})

	t.Run("Self ref struct", func(t *testing.T) {
		personInput := PersonInput{
			ID:   utils.IntP(1),
			Name: utils.StringP("John"),
			Role: "admin",
			Inviter: &PersonInput{
				ID:   utils.IntP(2),
				Name: utils.StringP("Jane"),
				Role: "admin",
			},
		}
		person := &Person{}
		err := mapper.Map(personInput, person)
		assert.Nil(t, err)
		assert.Equal(t, *personInput.ID, person.ID)
		assert.Equal(t, *personInput.Name, person.Name)
		assert.Equal(t, personInput.Role, *person.Role)
		assert.Equal(t, *personInput.Inviter.ID, person.Inviter.ID)
		assert.Equal(t, *personInput.Inviter.Name, person.Inviter.Name)
		assert.Equal(t, personInput.Inviter.Role, *person.Inviter.Role)
	})

	t.Run("Self ref slice", func(t *testing.T) {
		personInput := PersonInput{
			ID:   utils.IntP(1),
			Name: utils.StringP("John"),
			Role: "admin",
			Parents: []*PersonInput{
				{
					ID:   utils.IntP(2),
					Name: utils.StringP("Jane"),
					Role: "admin",
				},
			},
		}
		person := &Person{}
		err := mapper.Map(personInput, person)
		assert.Nil(t, err)
		assert.Equal(t, *personInput.ID, person.ID)
		assert.Equal(t, *personInput.Name, person.Name)
		assert.Equal(t, personInput.Role, *person.Role)
		assert.Equal(t, *personInput.Parents[0].ID, person.Parents[0].ID)
		assert.Equal(t, *personInput.Parents[0].Name, person.Parents[0].Name)
		assert.Equal(t, personInput.Parents[0].Role, *person.Parents[0].Role)
	})
}
