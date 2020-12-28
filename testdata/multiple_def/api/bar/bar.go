package bar

import (
	"net/http"

	"github.com/swaggo/swag/testdata/multiple_def/models/mbar"
)

// @Description get Bar
// @ID get-bar
// @Accept json
// @Produce json
// @Category bar
// @Success 200 {object} mbar.Bar
// @Router /testapi/get-bar [get]
func GetBar(w http.ResponseWriter, r *http.Request) {
	//write your code
	var _ = mbar.Bar{}
}
