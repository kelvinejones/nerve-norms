class Filter {
	constructor(url, action) {
		var form = document.querySelector("form")

		this.applyFilter = (event) => {
			var data = new FormData(form)
			this.filter = { "species": "human", "nerve": "median" }
			for (const entry of data) {
				switch (entry[0]) {
					case "sex-options":
						this.filter.sex = entry[1]
					case "country-options":
						this.filter.country = entry[1]
					case "age-options":
						this.setAgeOptions(entry[1])
				}
			}

			fetch(url + this.filterAsQueryString())
				.then(function(response) {
					return response.json()
				})
				.then(function(myJson) {
					action(myJson)
				})
			event.preventDefault()
		}

		form.addEventListener("submit", this.applyFilter)
	}

	setAgeOptions(opts) {
		switch (opts) {
			case "any":
				this.filter.minAge = 0
				this.filter.maxAge = 200
			case "-30":
				this.filter.minAge = 0
				this.filter.maxAge = 30
			case "31-40":
				this.filter.minAge = 31
				this.filter.maxAge = 40
			case "41-50":
				this.filter.minAge = 41
				this.filter.maxAge = 50
			case "51-60":
				this.filter.minAge = 51
				this.filter.maxAge = 60
			case "61-":
				this.filter.minAge = 61
				this.filter.maxAge = 200
		}
	}

	filterAsQueryString() {
		const filt = this.filter
		return "?" +
			Object.keys(filt).map(function(key) {
				return encodeURIComponent(key) + "=" +
					encodeURIComponent(filt[key]);
			}).join("&");
	}
}
