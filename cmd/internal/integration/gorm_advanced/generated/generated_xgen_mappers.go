package generated
import (
	"context"
	"github.com/goxgen/goxgen/plugins/cli/server"
)

// ToCarModel Map CarInput to Car model
func (ra *CarInput) ToCarModel(ctx context.Context) (*Car, error){
	mapper := server.GetMapper(ctx)
	target := &Car{}
	err := mapper.Map(ra, target)
	return target, err
}

// ToCarModel Map CarBrowseInput to Car model
func (ra *CarBrowseInput) ToCarModel(ctx context.Context) (*Car, error){
	mapper := server.GetMapper(ctx)
	target := &Car{}
	err := mapper.Map(ra, target)
	return target, err
}

// ToPhoneModel Map PhoneNumberBrowseInput to Phone model
func (ra *PhoneNumberBrowseInput) ToPhoneModel(ctx context.Context) (*Phone, error){
	mapper := server.GetMapper(ctx)
	target := &Phone{}
	err := mapper.Map(ra, target)
	return target, err
}

// ToPhoneModel Map PhoneNumberInput to Phone model
func (ra *PhoneNumberInput) ToPhoneModel(ctx context.Context) (*Phone, error){
	mapper := server.GetMapper(ctx)
	target := &Phone{}
	err := mapper.Map(ra, target)
	return target, err
}

// ToUserModel Map ListUser to User model
func (ra *ListUser) ToUserModel(ctx context.Context) (*User, error){
	mapper := server.GetMapper(ctx)
	target := &User{}
	err := mapper.Map(ra, target)
	return target, err
}

// ToUserModel Map UserInput to User model
func (ra *UserInput) ToUserModel(ctx context.Context) (*User, error){
	mapper := server.GetMapper(ctx)
	target := &User{}
	err := mapper.Map(ra, target)
	return target, err
}

// ToUserModel Map DeleteUsers to User model
func (ra *DeleteUsers) ToUserModel(ctx context.Context) (*User, error){
	mapper := server.GetMapper(ctx)
	target := &User{}
	err := mapper.Map(ra, target)
	return target, err
}