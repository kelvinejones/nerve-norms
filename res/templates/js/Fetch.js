class Fetch {
	static get url() { return "https://us-central1-nervenorms-294404.cloudfunctions.net/" }

	static MEM(queryString, data, callback) {
		const query = Fetch.url + "convert" + queryString
		fetch(query, { method: 'POST', body: data })
			.then(response => {
				return response.json()
			})
			.then(convertedMem => {
				callback(convertedMem)
			})
	}

	static OutliersFromName(queryString, name, callback) {
		fetch(Fetch.url + "outliers" + queryString + "&name=" + name)
			.then(response => {
				return response.json()
			})
			.then(scores => {
				callback(scores)
			})
	}

	static OutliersFromJSON(queryString, data, callback) {
		fetch(Fetch.url + "outliers" + queryString, {
				method: 'POST',
				body: JSON.stringify(data)
			})
			.then(response => {
				return response.json()
			})
			.then(scores => {
				callback(scores)
			})
	}

	static Norms(queryString, callback) {
		fetch(Fetch.url + "norms" + queryString)
			.then(response => {
				return response.json()
			})
			.then(norms => {
				callback(norms)
			})
	}

	//API for email is in SendEmail.js, due to it being on multiple pages
}
