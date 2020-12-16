package api

import (
	"fmt"
	"net/http"
)

// Handler for serverless function.
func Handler(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()

	bg, _bg := params["bg"]     // Get background image url
	fg, _fg := params["fg"]     // Get foreground image url
	_, _b64 := params["base64"] // Return as raw base64

	// Check is foreground availabe
	if !_fg {
		res.WriteHeader(400)
		fmt.Fprintf(res, "Missing fg (foreground) parameter!")
	} else {
		res.WriteHeader(200)
		if _b64 {
			res.Write([]byte(fg[0] + bg[0]))
		} else {
			fmt.Fprintf(res, "Background: %s\nForeground: %s", bg[0], fg[0])
		}
	}

}
