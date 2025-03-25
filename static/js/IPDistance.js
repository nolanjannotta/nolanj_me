
async function getDistance() {

  const clientIPResp = await fetch("https://api.ipify.org?format=json")
  const clientIPJson = await clientIPResp.json()

  const response = await fetch(window.location.href + "distanceToLa/" + clientIPJson.ip)
  const json = await response.json()

  document.getElementById("distance").textContent = json.distance.toFixed(2) + " km " + json.direction


}

document.onload = getDistance()