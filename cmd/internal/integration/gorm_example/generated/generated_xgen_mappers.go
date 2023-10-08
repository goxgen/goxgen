package generated
import (
	"context"
	"github.com/goxgen/goxgen/plugins/cli/server"
)

// ToPhoneModel Map PhoneNumberInput to Phone model
func (ra *PhoneNumberInput) ToPhoneModel(ctx context.Context) (*Phone, error){
	mapper := server.GetMapper(ctx)
	target := &Phone{}
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

// ToUserModel Map BrowseUserInput to User model
func (ra *BrowseUserInput) ToUserModel(ctx context.Context) (*User, error){
	mapper := server.GetMapper(ctx)
	target := &User{}
	err := mapper.Map(ra, target)
	return target, err
}