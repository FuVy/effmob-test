package handler

var (
	// cars 2xxx
	errCarAlreadyAdded = &CodeResponse{2001, "car already added"}
	errCarNotFound     = &CodeResponse{2002, "car not found"}
	errCarInvalidID    = &CodeResponse{2003, "car invalid id"}
	errUnknownCars     = &CodeResponse{2004, "cant get data about any car"}
	errNoCarsProvided  = &CodeResponse{2005, "no cars provided"}
)
