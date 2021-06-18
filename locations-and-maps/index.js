const AWS = require("aws-sdk");
const etc = require("./etc");
const deliveryDetails = require("./delivery-details");
const { waypoints } = require("./delivery-details");

AWS.config.update({
  region: process.env.AWS_REGION || 'eu-west-1'
});


const location = new AWS.Location({
  region: process.env.AWS_REGION || 'eu-west-1'
});

async function makeRoute(depot, waypoints) {
  const resp = await location.calculateRoute({
    CalculatorName: "rc1", // calculator resource needs to be created in advance
    DepartureTime: etc.getTomorrowMorning(),
    DeparturePosition: depot,
    DestinationPosition: depot,
    TravelMode: 'Car',
    CarModeOptions: {
      AvoidFerries: false,
      AvoidTolls: true
    },
    WaypointPositions: waypoints
  }).promise();
  return { legs: resp.Legs, summary: resp.Summary };
}

async function geocode(address) {
  const resp = await location.searchPlaceIndexForText({
    IndexName: "explore.place", //should be created by AWS on first run (create, if absent)
    Text: address,
    MaxResults: 1
  }).promise();
  return resp.Results.shift().Place.Geometry.Point;
}

async function buildRouteForTomorrow() {
  const depotCoords = await geocode(deliveryDetails.depot);
  console.log('depot coordinates:\r\n', depotCoords);
  const waypointCoords = [];
  const waypointToAddressMap = new Map();
  await Promise.all(deliveryDetails.waypoints.map((waypoint) => {
    return geocode(waypoint).then((data) => {
      waypointToAddressMap.set(`${data[0]}:${data[1]}`, waypoint);
      waypointCoords.push(data);
    });
  }));
  console.log('waypoints coordinates:\r\n', waypointCoords);
  const waypointCoordsPermutations = etc.getArrayMutations(waypointCoords);

  const routesPermutations = await Promise.all(waypointCoordsPermutations.map((combination) => {
    return makeRoute(depotCoords, combination).then((route) => {
      const distanceKm = Math.round(route.summary.Distance * 100) / 100;
      const durationMin = Math.round(route.summary.DurationSeconds / 60 * 100) / 100;
      console.log(`\tPossible route: km:${distanceKm}, min:${durationMin}`);
      return {
        combination,
        distanceKm,
        durationMin
      };
    });
  }));

  let minDurationRoute = routesPermutations.shift();
  routesPermutations.forEach((route) => {
    if (minDurationRoute.durationMin > route.durationMin) {
      minDurationRoute = route;
    }
  })
  console.log('The quickest route:', minDurationRoute);
  let minDurationWaypoints = minDurationRoute.combination.map((waypoint) => waypointToAddressMap.get(`${waypoint[0]}:${waypoint[1]}`));
  console.log('The quickest route waypoints order:', minDurationWaypoints);

}
buildRouteForTomorrow();
