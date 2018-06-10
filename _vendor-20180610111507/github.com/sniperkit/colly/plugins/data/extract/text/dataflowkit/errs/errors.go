// Dataflow kit - Errs
//
// Copyright © 2017-2018 Slotix s.r.o. <dm@slotix.sk>
//
//
// All rights reserved. Use of this source code is governed
// by the BSD 3-Clause License license.

package errs

//	Network errors
//
// BadRequest 400
//
// The server cannot or will not process the request due to an apparent client error (e.g., malformed request syntax, size too large, invalid request message framing, or deceptive request routing).
type BadRequest struct {
	Err error
}

func (e *BadRequest) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return "400 Bad Request"
}

// Unauthorized 401
//
// Client does not have access rights to the content.
type Unauthorized struct {
}

func (e *Unauthorized) Error() string { return "401 Unauthorized" }

// ForbiddenByRobots 403
//
// Client does not have access rights to the content caused by robots.txt restrictions.
type ForbiddenByRobots struct {
	URL string
}

func (e *ForbiddenByRobots) Error() string { return "403 Forbidden by robots.txt: "+ e.URL}

// Forbidden 403
//
// Client does not have access rights to the content so server is rejecting to give proper response.
type Forbidden struct {
	URL string
}

func (e *Forbidden) Error() string { return "403 Forbidden: "+ e.URL}

// NotFound 404
//
// Server can not find requested resource. This response code probably is most famous one due to its frequency to occur in web.
type NotFound struct {
	URL string
}

func (e *NotFound) Error() string {
	return "404 Not found: " + e.URL
}

//500 Internal Server Error
//A generic error message, given when an unexpected condition was encountered and no more specific message is suitable
type InternalServerError struct {
}

func (*InternalServerError) Error() string {
	return "500 Internal Server Error"
}

// BadGateway 502
//
// This error response means that the server, while working as a gateway to get a response needed to handle the request, got an invalid response.
type BadGateway struct {
	What string
}

func (e *BadGateway) Error() string {
	return "502 Invalid "+ e.What + " from server"
}

// GatewayTimeout Gateway Time-out 504
//
// This error response is given when the server is acting as a gateway and cannot get a response in time.
type GatewayTimeout struct {
}

func (e *GatewayTimeout) Error() string {
	return "504 Timeout exceeded rendering page"
}

//Parser Errors generated by Dataflow kit Parser service

//ParserError returned if web page cannot be parsed correctly due to wrong payload structure
type ParserError string

const (
	ErrNoParts          ParserError = "no parts found"
	ErrNoSelectors                  = "no selectors found"
	ErrEmptyResults                 = "empty results"
	ErrNoCommonAncestor             = "no common ancestor for selectors found"
)

//BadPayload error is returned if Payload is invalid 400
type BadPayload struct {
	ParserError
}

func (e *BadPayload) Error() string {
	return "400: " + string(e.ParserError)
}

// Error represents all the rest (unspecified errors).
type Error struct {
	Err string
}

func (e *Error) Error() string {
	return e.Err
}
