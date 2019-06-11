class ExVars {
	static get _interpolate() {
		const opacity = 0.8,
			red = d3.hsl("red"),
			green = d3.hsl("lightgreen");

		red.opacity = opacity
		green.opacity = opacity

		const intr = d3.interpolateHsl(green, red);
		return function(score) { return intr(Math.pow(1 - score, 3)) }
	}

	static _setExcitabilityVariableScore(idString, score) {
		const row = document.getElementById(idString)
		if (row === null || !row.classList.contains('display-bar')) {
			// We don't care about this variable
			return
		}
		if (row.classList.contains('was-imp')) {
			row.style.background = "#98AFC7"
			return
		}
		if (score === undefined) {
			score = 1
		}
		row.style.background = "linear-gradient(to right, " + ExVars._interpolate(score) + " " + (1 - score) * 100 + "%, #ffffff 0%)"
	}

	static _setExcitabilityVariableValue(idString, value, wasimp) {
		const row = document.getElementById(idString);
		if (row === null) {
			// We don't care about this variable
			return
		}
		if (wasimp != row.classList.contains("was-imp")) {
			row.classList.toggle("was-imp");
		}
		if (value === undefined) {
			row.getElementsByClassName("excite-value")[0].innerHTML = ""
		} else {
			// If it's anumber and if the length of the number is more than 2 decimals, truncate
			if (Number(value) === value && (value % 1) != 0 && value.toString().split(".")[1].length > 3) {
				value = value.toFixed(3)
			}
			row.getElementsByClassName("excite-value")[0].innerHTML = value
		}
	}

	static _setHeaderScore(str, score) {
		[...document.querySelectorAll(str)].forEach(elm => {
			if (score === undefined) {
				elm.style.background = ''
			} else {
				elm.style.background = ExVars._interpolate(score)
			}
		})
	}

	static _labelForScore(score) {
		if (score > 0.25) {
			return "Healthy"
		} else if (score > .05) {
			return "Atypical"
		} else {
			return "Extremely Atypical"
		}
	}

	static updateScores(scores) {
		ExVars._setHeaderScore(".participant-header", scores.Overall)
		const healthLabel = ExVars._labelForScore(scores.Overall)
		ExVars._setExcitabilityVariableValue("overall-score", healthLabel + " (" + scores.Overall.toFixed(2) + ")", false)

		Object.keys(scores).forEach(function(key) {
			ExVars._setHeaderScore("." + key + "-header", scores[key].Overall)
		})

		scores.ExVars.data.forEach(function(exind) {
			if (exind[1] === 0) {
				return // This means it doesn't have an index
			}
			ExVars._setExcitabilityVariableScore("qtrac-excite-" + exind[1], exind[0])
		})
	}

	static updateValues(values) {
		values.sections.ExVars.data.forEach(function(exind) {
			const idx = exind[0]
			if (idx === 0) {
				return // This means it doesn't have an index
			}
			let val = exind[1]
			if (idx == 18) {
				// This is sex, so treat it differently
				if (exind[1] == 1) {
					val = "Male"
				} else {
					val = "Female"
				}
			}
			ExVars._setExcitabilityVariableValue("qtrac-excite-" + idx, val, exind[2] == 1)
		})
	}

	static clearScores() {
		ExVars._setHeaderScore(".participant-header", undefined)
		ExVars._setHeaderScore(".TE-header", undefined)
		ExVars._setHeaderScore(".RC-header", undefined)
		ExVars._setHeaderScore(".IV-header", undefined)
		ExVars._setHeaderScore(".CD-header", undefined)
		ExVars._setHeaderScore(".SR-header", undefined)
		ExVars._setHeaderScore(".SRel-header", undefined)
		ExVars._setHeaderScore(".ExVars-header", undefined)
		ExVars._setExcitabilityVariableValue("overall-score", "loading...", false)
		const elms = document.getElementsByClassName('qtrac-excite')
		for (let elm of elms) {
			ExVars._setExcitabilityVariableScore(elm.id, undefined)
		}
	}
}
