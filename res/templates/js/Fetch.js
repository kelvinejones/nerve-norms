class Fetch {
	static get url() { return "https://us-central1-nervenorms.cloudfunctions.net/" }

	static MEM(queryString, name, data, callback) {
		ExVars.clearScores()

		const query = Fetch.url + "convert" + queryString
		fetch(query, { method: 'POST', body: data })
			.then(response => {
				return response.json()
			})
			.then(convertedMem => {
				const name = callback(convertedMem)
				ExVars.updateScores(name, convertedMem.outlierScores)
			})
	}

	static Outliers(queryString, name, data) {
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
				ExVars.updateScores(name, scores)
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
