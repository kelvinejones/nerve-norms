class Filter {
	static asQueryString() {
		var data = new FormData(document.querySelector("form"))
		const filter = { "species": "human", "nerve": "median" }
		for (const entry of data) {
			switch (entry[0]) {
				case "sex-options":
					filter.sex = entry[1]
				case "country-options":
					filter.country = entry[1]
				case "age-options":
					Filter.setAgeOptions(filter, entry[1])
			}
		}

		return "?" +
			Object.keys(filter).map(function(key) {
				return encodeURIComponent(key) + "=" +
					encodeURIComponent(filter[key]);
			}).join("&");
	}

	static setAgeOptions(filter, opts) {
		switch (opts) {
			case "any":
				filter.minAge = 0
				filter.maxAge = 200
			case "-30":
				filter.minAge = 0
				filter.maxAge = 30
			case "31-40":
				filter.minAge = 31
				filter.maxAge = 40
			case "41-50":
				filter.minAge = 41
				filter.maxAge = 50
			case "51-60":
				filter.minAge = 51
				filter.maxAge = 60
			case "61-":
				filter.minAge = 61
				filter.maxAge = 200
		}
	}
}
