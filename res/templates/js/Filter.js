class Filter {
	static get url() { return "https://us-central1-nervenorms.cloudfunctions.net/" }

	constructor(callback) {
		this.callback = callback

		document.querySelector("form").addEventListener("submit", (event) => {
			this.update(this.name)
			event.preventDefault()
		})
	}

	setParticipantData(data) {
		this.name = undefined
		this.data = data
	}

	update(name) {
		const lastQuery = this.queryString
		this.queryString = Filter._queryString
		const normChanged = (lastQuery != this.queryString)
		if (normChanged) {
			this._fetchNorms()
		}

		if (name != undefined) {
			this.data = undefined
		}
		const lastParticipant = this.name
		this.name = name
		const nameChanged = (lastParticipant != this.name)
		if (normChanged || nameChanged) {
			ExVars.clearScores()
			this._fetchOutliers(this.data)
		}
	}

	static get _queryString() {
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

	fetchMEM(data, callback) {
		this.lastParticipant = undefined
		ExVars.clearScores()

		const query = Filter.url + "convert" + this.queryString
		fetch(query, { method: 'POST', body: data })
			.then(response => {
				return response.json()
			})
			.then(convertedMem => {
				const name = callback(convertedMem)
				ExVars.updateScores(this.name, convertedMem.outlierScores)
			})
		return this
	}

	_fetchOutliers(data) {
		let query = Filter.url + "outliers" + this.queryString
		let body = undefined
		if (this.name != null) {
			query = query + "&name=" + this.name
		} else if (data != null) {
			body = { method: 'POST', body: JSON.stringify(data) }
		} else {
			console.log("Error in fetch", this.name, data)
		}
		fetch(query, body)
			.then(response => {
				return response.json()
			})
			.then(scores => {
				ExVars.updateScores(this.name, scores)
			})
		return this
	}

	_fetchNorms() {
		fetch(Filter.url + "norms" + this.queryString)
			.then(response => {
				return response.json()
			})
			.then(norms => {
				this.callback(norms)
			})
		return this
	}
}
