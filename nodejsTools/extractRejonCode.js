

async function a(){
	const url = "https://api.windy.com/webcams/api/v3/regions?lang=en";

	const headers = {
  		"accept": "application/json",
  		"x-windy-api-key": "4tpguJklGSjb3f0nVny1wwR9bqHquToz",
	};

	const response = await fetch(url, {
 		method: "GET",
  		headers,
	});

	jsonArray = await response.json();
	
	console.log(jsonArray);
	const codes = jsonArray.filter((item) => item.hasOwnProperty("code")).map((item) => item.code);
	console.log(codes);
}


a();
