package doc

// swagger:route POST /api/v1/harbors harbors UpsertHarbor
// Upsert harbors by either creating new entries or updating existing ones.
//
// This endpoint receives a JSON object containing multiple harbor entries, where each key is the UN/LOCODE
// and the value is the harbor's details. It processes the input **streaming-style** for efficiency, ensuring
// that even large files do not exhaust memory. The service **validates** the structure of each harbor before
// attempting to insert or update it in the database.
//
// ---
// produces:
// - application/json
// consumes:
// - application/json
// responses:
//   200: upsertHarborResponse
//   400: badRequestResponse
//   500: internalServerErrorResponse

// swagger:parameters UpsertHarbor
type UpsertHarborParamsWrapper struct {
	// in:body
	Body map[string]struct {
		// Harbor name
		// example: Los Angeles
		Name string `json:"name"`

		// City where the harbor is located
		// example: Los Angeles
		City string `json:"city"`

		// Country where the harbor is located
		// example: United States
		Country string `json:"country"`

		// List of alternative names for the harbor
		// example: ["LA Harbor", "Port of LA"]
		Alias []string `json:"alias"`

		// List of associated regions
		// example: ["West Coast"]
		Regions []string `json:"regions"`

		// Geographic coordinates of the harbor
		// example: [-118.2437, 34.0522]
		Coordinates []float64 `json:"coordinates"`

		// Province or state where the harbor is located
		// example: California
		Province string `json:"province"`

		// Timezone of the harbor
		// example: America/Los_Angeles
		Timezone string `json:"timezone"`

		// UN/LOCODE(s) associated with the harbor
		// example: ["USLAX"]
		UNLocs []string `json:"unlocs"`

		// Port code
		// example: 53001
		Code string `json:"code"`
	}
}

// swagger:response upsertHarborResponse
type UpsertHarborResponseWrapper struct {
	// in:body
	Body struct {
		// Success message
		// example: harbors upserted
		Message string `json:"message"`
	}
}

// swagger:response badRequestResponse
type BadRequestResponseWrapper struct {
	// in:body
	Body struct {
		// Error message
		// example: invalid JSON structure
		Error string `json:"error"`
	}
}

// swagger:response internalServerErrorResponse
type InternalServerErrorResponseWrapper struct {
	// in:body
	Body struct {
		// Error message
		// example: internal server error
		Error string `json:"error"`
	}
}
