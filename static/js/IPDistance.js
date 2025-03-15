const element = document.getElementById("distance")

async function getDistance() {
  const laLat = 34.052235
  const laLon = -118.243683

  const clientIPResp = await fetch("https://api.ipify.org?format=json")
  const clientIPJson = await clientIPResp.json()
  // console.log(clientIPJson)

  const clientLatLonResp = await fetch("http://ip-api.com/json/" + clientIPJson.ip)
  const clientLatLonJson = await clientLatLonResp.json()

  let lat1 = laLat * Math.PI / 180;
  let lon1 = laLon * Math.PI / 180;

  let lat2 = clientLatLonJson.lat * Math.PI / 180;
  let lon2 = clientLatLonJson.lon * Math.PI / 180;

  let dlon = lon2-lon1
  let y = Math.sin(dlon) * Math.cos(lat2)
  let x = Math.cos(lat1) * Math.sin(lat2) - Math.sin(lat1) * Math.cos(lat2) * Math.cos(dlon)
  let bearing = Math.atan2(y,x)
  bearing *= (180 / Math.PI)
  // bearing = (bearing + 360) % 360

  let directions = ["N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", 
    "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"]
  let index = Math.floor(bearing / 22.5) % 16


  console.log(directions[index])


  const km = Math.acos(Math.sin(lat1) * Math.sin(lat2) + Math.cos(lat1) * Math.cos(lat2) * Math.cos(lon2 - lon1)) * 6371;
  // const miles = Math.floor(km * 0.621371)
  element.textContent = Math.floor(km) + " km away"


}

document.onload = getDistance()