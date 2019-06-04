class ExVars {
	static get _interpolate() {
		const opacity = 0.8,
			red = d3.hsl("red"),
			green = d3.hsl("lightgreen");

		red.opacity = opacity
		green.opacity = opacity

		const intr = d3.interpolateHsl(green, red);
		return function(score) { return intr(Math.pow(score, 3)) }
	}

	static _setExcitabilityVariableScore(idString, score) {
		const row = document.getElementById(idString);
		if (row === null) {
			// We don't care about this variable
			return
		}
		if (score === undefined) {
			score = 0
		}
		row.style.background = "linear-gradient(to right, " + ExVars._interpolate(score) + " " + score * 100 + "%, #ffffff 0%)"
	}

	static _setExcitabilityVariableValue(idString, value) {
		const row = document.getElementById(idString);
		if (row === null) {
			// We don't care about this variable
			return
		}
		if (value === undefined) {
			row.getElementsByClassName("excite-value")[0].innerHTML = ""
		} else {
			row.getElementsByClassName("excite-value")[0].innerHTML = value
		}
	}

	static _setHeaderScore(str, score) {
		[...document.querySelectorAll(str)].forEach(elm => {
			elm.style.background = ExVars._interpolate(score);
		})
	}

	static _labelForScore(score) {
		if (score < 0.75) {
			return "Healthy"
		} else if (score < .95) {
			return "Atypical"
		} else {
			return "Extremely Atypical"
		}
	}

	static updateScores(scores) {
		scores.ExVars.data.forEach(function(exind) {
			if (exind[1] === 0) {
				return // This means it doesn't have an index
			}
			ExVars._setExcitabilityVariableScore("qtrac-excite-" + exind[1], exind[0])
		})
	}

	static updateValues(values) {
		values.sections.ExVars.data.forEach(function(exind) {
			if (exind[0] === 0) {
				return // This means it doesn't have an index
			}
			ExVars._setExcitabilityVariableValue("qtrac-excite-" + exind[0], exind[1])
		})
	}

	static setToZero() {
		const elms = document.getElementsByClassName('qtrac-excite')
		for (let elm of elms) {
			ExVars._setExcitabilityVariableScore(elm.id, undefined)
			ExVars._setExcitabilityVariableValue(elm.id, undefined)
		}
	}
}
