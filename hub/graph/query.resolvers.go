package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"

	"github.com/kiel-live/kiel-live/hub/graph/model"
)

// Stop is the resolver for the stop field.
func (r *queryResolver) Stop(ctx context.Context, id string) (*model.Stop, error) {
	panic(fmt.Errorf("not implemented: Stop - stop"))
}

// Vehicle is the resolver for the vehicle field.
func (r *queryResolver) Vehicle(ctx context.Context, id string) (*model.Vehicle, error) {
	panic(fmt.Errorf("not implemented: Vehicle - vehicle"))
}

// Trip is the resolver for the trip field.
func (r *queryResolver) Trip(ctx context.Context, id string) (*model.Trip, error) {
	panic(fmt.Errorf("not implemented: Trip - trip"))
}

// Route is the resolver for the route field.
func (r *queryResolver) Route(ctx context.Context, id string) (*model.Route, error) {
	panic(fmt.Errorf("not implemented: Route - route"))
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryResolver) Stops(ctx context.Context) ([]*model.Stop, error) {
	panic(fmt.Errorf("not implemented: Stops - stops"))
}
func (r *queryResolver) Vehicles(ctx context.Context) ([]*model.Vehicle, error) {
	panic(fmt.Errorf("not implemented: Vehicles - vehicles"))
}
func (r *queryResolver) Trips(ctx context.Context) ([]*model.Trip, error) {
	panic(fmt.Errorf("not implemented: Trips - trips"))
}
func (r *queryResolver) Routes(ctx context.Context) ([]*model.Route, error) {
	panic(fmt.Errorf("not implemented: Routes - routes"))
}
