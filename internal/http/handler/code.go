package handler

var (
	// cars 2xxx
	errCarAlreadyAdded = &CodeResponse{2001, "car already added"}
	errCarNotFound     = &CodeResponse{2002, "car not found"}
)
