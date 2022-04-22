package timezone

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func CurrentTimeZone(w http.ResponseWriter, r *http.Request) {
	date_time := time.Now()
	zone, _ := date_time.Zone()
	fmt.Fprint(w, "Current time is : ", date_time.String())
	fmt.Fprint(w, "Current timezone : ", zone)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", CurrentTimeZone)

}
