package main #blah

import (
	"fmt"
    "github.com/nowaitapp/guestApi"
    "github.com/nowaitapp/humanReadableSchedule"
)

func main() {
    request := guestApi.GetRestaruantsRequest{Lat:40,Lon:-79}
    response,err := request.Send();
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

    fmt.Printf("%+v\n", response)

    for _,restaurant := range response.Restaurants {
    	schedule := humanReadableSchedule.New(restaurant.Schedule)
    	fmt.Printf("%s\n", restaurant.Name)
    	fmt.Printf(schedule.String())
    }

}
