class Fetch {
	static get url() { return "https://us-central1-nervenorms.cloudfunctions.net/" }

	static MEM(queryString, data, callback) {
		ExVars.clearScores()

		const query = Fetch.url + "convert" + queryString
		fetch(query, { method: 'POST', body: data })
			.then(response => {
				return response.json()
			})
			.then(convertedMem => {
				callback(convertedMem)
			})
	}

	static Outliers(queryString, name, data, callback) {
		let query = Fetch.url + "outliers" + queryString
		let body = undefined
		if (name != null) {
			query = query + "&name=" + name
		} else if (data != null) {
			body = { method: 'POST', body: JSON.stringify(data) }
		} else {
			console.log("Error in fetch", name, data)
		}
		fetch(query, body)
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
}
