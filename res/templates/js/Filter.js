class Filter {
	static get url() { return "https://us-central1-nervenorms.cloudfunctions.net/" }

	static setCallback(submitCallback) {
		document.querySelector("form").addEventListener("submit", (event) => {
			submitCallback()
			event.preventDefault()
		})
	}

	static get queryString() {
		var data = new FormData(document.querySelector("form"))
		const filter = { "species": "human", "nerve": "median" }
		for (const entry of data) {
			switch (entry[0]) {
				case "sex-options":
					filter.sex = entry[1]
					break
				case "country-options":
					filter.country = entry[1]
					break
				case "age-options":
					Filter._setAgeOptions(filter, entry[1])
					break
			}
		}

		return "?" +
			Object.keys(filter).map(function(key) {
				return encodeURIComponent(key) + "=" +
					encodeURIComponent(filter[key]);
			}).join("&");
	}

	static _setAgeOptions(filter, opts) {
		switch (opts) {
			case "any":
				filter.minAge = 0
				filter.maxAge = 200
				break
			case "-30":
				filter.minAge = 0
				filter.maxAge = 30
				break
			case "31-40":
				filter.minAge = 31
				filter.maxAge = 40
				break
			case "41-50":
				filter.minAge = 41
				filter.maxAge = 50
				break
			case "51-60":
				filter.minAge = 51
				filter.maxAge = 60
				break
			case "61-":
				filter.minAge = 61
				filter.maxAge = 200
				break
		}
	}

	static fetchMEM(queryString, name, data, callback) {
		ExVars.clearScores()

		const query = Filter.url + "convert" + queryString
		fetch(query, { method: 'POST', body: data })
			.then(response => {
				return response.json()
			})
			.then(convertedMem => {
				const name = callback(convertedMem)
				ExVars.updateScores(name, convertedMem.outlierScores)
			})
	}

	static fetchOutliers(queryString, name, data) {
		let query = Filter.url + "outliers" + queryString
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

	static fetchNorms(queryString, callback) {
		fetch(Filter.url + "norms" + queryString)
			.then(response => {
				return response.json()
			})
			.then(norms => {
				callback(norms)
			})
	}
}
