// Package middleware contains middleware for handling error and logging propagation
package middleware

import "github.com/aaronsky/applereleaser/pkg/context"

// Action is a function that takes a context and returns an error.
// It is is used on Pipers, Defaulters and Publishers, although they are not
// aware of this generalization.
type Action func(ctx *context.Context) error
