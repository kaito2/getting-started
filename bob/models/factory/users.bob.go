// Code generated by BobGen mysql v0.30.0. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package factory

import (
	"context"
	"testing"

	"github.com/aarondl/opt/null"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/jaswdr/faker/v2"
	models "github.com/kaito2/getting-started/bob/models"
	"github.com/stephenafamo/bob"
)

type UserMod interface {
	Apply(*UserTemplate)
}

type UserModFunc func(*UserTemplate)

func (f UserModFunc) Apply(n *UserTemplate) {
	f(n)
}

type UserModSlice []UserMod

func (mods UserModSlice) Apply(n *UserTemplate) {
	for _, f := range mods {
		f.Apply(n)
	}
}

// UserTemplate is an object representing the database table.
// all columns are optional and should be set by mods
type UserTemplate struct {
	ID   func() int32
	Name func() null.Val[string]

	f *Factory
}

// Apply mods to the UserTemplate
func (o *UserTemplate) Apply(mods ...UserMod) {
	for _, mod := range mods {
		mod.Apply(o)
	}
}

// toModel returns an *models.User
// this does nothing with the relationship templates
func (o UserTemplate) toModel() *models.User {
	m := &models.User{}

	if o.ID != nil {
		m.ID = o.ID()
	}
	if o.Name != nil {
		m.Name = o.Name()
	}

	return m
}

// toModels returns an models.UserSlice
// this does nothing with the relationship templates
func (o UserTemplate) toModels(number int) models.UserSlice {
	m := make(models.UserSlice, number)

	for i := range m {
		m[i] = o.toModel()
	}

	return m
}

// setModelRels creates and sets the relationships on *models.User
// according to the relationships in the template. Nothing is inserted into the db
func (t UserTemplate) setModelRels(o *models.User) {}

// BuildSetter returns an *models.UserSetter
// this does nothing with the relationship templates
func (o UserTemplate) BuildSetter() *models.UserSetter {
	m := &models.UserSetter{}

	if o.ID != nil {
		m.ID = omit.From(o.ID())
	}
	if o.Name != nil {
		m.Name = omitnull.FromNull(o.Name())
	}

	return m
}

// BuildManySetter returns an []*models.UserSetter
// this does nothing with the relationship templates
func (o UserTemplate) BuildManySetter(number int) []*models.UserSetter {
	m := make([]*models.UserSetter, number)

	for i := range m {
		m[i] = o.BuildSetter()
	}

	return m
}

// Build returns an *models.User
// Related objects are also created and placed in the .R field
// NOTE: Objects are not inserted into the database. Use UserTemplate.Create
func (o UserTemplate) Build() *models.User {
	m := o.toModel()
	o.setModelRels(m)

	return m
}

// BuildMany returns an models.UserSlice
// Related objects are also created and placed in the .R field
// NOTE: Objects are not inserted into the database. Use UserTemplate.CreateMany
func (o UserTemplate) BuildMany(number int) models.UserSlice {
	m := make(models.UserSlice, number)

	for i := range m {
		m[i] = o.Build()
	}

	return m
}

func ensureCreatableUser(m *models.UserSetter) {
	if m.ID.IsUnset() {
		m.ID = omit.From(random_int32(nil))
	}
}

// insertOptRels creates and inserts any optional the relationships on *models.User
// according to the relationships in the template.
// any required relationship should have already exist on the model
func (o *UserTemplate) insertOptRels(ctx context.Context, exec bob.Executor, m *models.User) (context.Context, error) {
	var err error

	return ctx, err
}

// Create builds a user and inserts it into the database
// Relations objects are also inserted and placed in the .R field
func (o *UserTemplate) Create(ctx context.Context, exec bob.Executor) (*models.User, error) {
	_, m, err := o.create(ctx, exec)
	return m, err
}

// MustCreate builds a user and inserts it into the database
// Relations objects are also inserted and placed in the .R field
// panics if an error occurs
func (o *UserTemplate) MustCreate(ctx context.Context, exec bob.Executor) *models.User {
	_, m, err := o.create(ctx, exec)
	if err != nil {
		panic(err)
	}
	return m
}

// CreateOrFail builds a user and inserts it into the database
// Relations objects are also inserted and placed in the .R field
// It calls `tb.Fatal(err)` on the test/benchmark if an error occurs
func (o *UserTemplate) CreateOrFail(ctx context.Context, tb testing.TB, exec bob.Executor) *models.User {
	tb.Helper()
	_, m, err := o.create(ctx, exec)
	if err != nil {
		tb.Fatal(err)
		return nil
	}
	return m
}

// create builds a user and inserts it into the database
// Relations objects are also inserted and placed in the .R field
// this returns a context that includes the newly inserted model
func (o *UserTemplate) create(ctx context.Context, exec bob.Executor) (context.Context, *models.User, error) {
	var err error
	opt := o.BuildSetter()
	ensureCreatableUser(opt)

	m, err := models.Users.Insert(opt).One(ctx, exec)
	if err != nil {
		return ctx, nil, err
	}
	ctx = userCtx.WithValue(ctx, m)

	ctx, err = o.insertOptRels(ctx, exec, m)
	return ctx, m, err
}

// CreateMany builds multiple users and inserts them into the database
// Relations objects are also inserted and placed in the .R field
func (o UserTemplate) CreateMany(ctx context.Context, exec bob.Executor, number int) (models.UserSlice, error) {
	_, m, err := o.createMany(ctx, exec, number)
	return m, err
}

// MustCreateMany builds multiple users and inserts them into the database
// Relations objects are also inserted and placed in the .R field
// panics if an error occurs
func (o UserTemplate) MustCreateMany(ctx context.Context, exec bob.Executor, number int) models.UserSlice {
	_, m, err := o.createMany(ctx, exec, number)
	if err != nil {
		panic(err)
	}
	return m
}

// CreateManyOrFail builds multiple users and inserts them into the database
// Relations objects are also inserted and placed in the .R field
// It calls `tb.Fatal(err)` on the test/benchmark if an error occurs
func (o UserTemplate) CreateManyOrFail(ctx context.Context, tb testing.TB, exec bob.Executor, number int) models.UserSlice {
	tb.Helper()
	_, m, err := o.createMany(ctx, exec, number)
	if err != nil {
		tb.Fatal(err)
		return nil
	}
	return m
}

// createMany builds multiple users and inserts them into the database
// Relations objects are also inserted and placed in the .R field
// this returns a context that includes the newly inserted models
func (o UserTemplate) createMany(ctx context.Context, exec bob.Executor, number int) (context.Context, models.UserSlice, error) {
	var err error
	m := make(models.UserSlice, number)

	for i := range m {
		ctx, m[i], err = o.create(ctx, exec)
		if err != nil {
			return ctx, nil, err
		}
	}

	return ctx, m, nil
}

// User has methods that act as mods for the UserTemplate
var UserMods userMods

type userMods struct{}

func (m userMods) RandomizeAllColumns(f *faker.Faker) UserMod {
	return UserModSlice{
		UserMods.RandomID(f),
		UserMods.RandomName(f),
	}
}

// Set the model columns to this value
func (m userMods) ID(val int32) UserMod {
	return UserModFunc(func(o *UserTemplate) {
		o.ID = func() int32 { return val }
	})
}

// Set the Column from the function
func (m userMods) IDFunc(f func() int32) UserMod {
	return UserModFunc(func(o *UserTemplate) {
		o.ID = f
	})
}

// Clear any values for the column
func (m userMods) UnsetID() UserMod {
	return UserModFunc(func(o *UserTemplate) {
		o.ID = nil
	})
}

// Generates a random value for the column using the given faker
// if faker is nil, a default faker is used
func (m userMods) RandomID(f *faker.Faker) UserMod {
	return UserModFunc(func(o *UserTemplate) {
		o.ID = func() int32 {
			return random_int32(f)
		}
	})
}

// Set the model columns to this value
func (m userMods) Name(val null.Val[string]) UserMod {
	return UserModFunc(func(o *UserTemplate) {
		o.Name = func() null.Val[string] { return val }
	})
}

// Set the Column from the function
func (m userMods) NameFunc(f func() null.Val[string]) UserMod {
	return UserModFunc(func(o *UserTemplate) {
		o.Name = f
	})
}

// Clear any values for the column
func (m userMods) UnsetName() UserMod {
	return UserModFunc(func(o *UserTemplate) {
		o.Name = nil
	})
}

// Generates a random value for the column using the given faker
// if faker is nil, a default faker is used
func (m userMods) RandomName(f *faker.Faker) UserMod {
	return UserModFunc(func(o *UserTemplate) {
		o.Name = func() null.Val[string] {
			if f == nil {
				f = &defaultFaker
			}

			if f.Bool() {
				return null.FromPtr[string](nil)
			}

			return null.From(random_string(f))
		}
	})
}
