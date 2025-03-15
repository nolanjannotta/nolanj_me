const element = document.getElementById("distance")

async function getDistance() {
  const laLat = 34.052235
  const laLon = -118.243683


  

  const clientIPResp = await fetch("https://api.ipify.org?format=json")
  const clientIPJson = await clientIPResp.json()
  // console.log(clientIPJson)

  const response = await fetch(window.location.href + "ip/" + clientIPJson.ip)
  const json = await response.json()

  element.textContent = json.distance.toFixed(2) + "km " + json.direction


}

document.onload = getDistance()