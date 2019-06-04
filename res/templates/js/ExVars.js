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
		if (score === undefined) {
			return
		}
		const row = document.getElementById(idString);
		if (row === null) {
			// We don't care about this variable
			return
		}
		row.style.background = "linear-gradient(to right, " + ExVars._interpolate(score) + " " + score * 100 + "%, #ffffff 0%)"
	}

	static _setExcitabilityVariableValue(idString, value) {
		const row = document.getElementById(idString);
		if (row === null) {
			// We don't care about this variable
			return
		}
		row.getElementsByClassName("excite-value")[0].innerHTML = value
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

	}

	static update(scores, values) {
		// ExVars._setHeaderScore(".participant-header", scores.outlierScore)
		// const healthLabel = ExVars._labelForScore(scores.outlierScore)
		// ExVars._setExcitabilityVariable("overall-score", healthLabel + " (" + scores.outlierScore.toFixed(2) + ")", 0)
		// const nameSpan = document.getElementById("participant-name");
		// nameSpan.innerHTML = scores.participant + " (" + healthLabel + ")"

		// Object.keys(scores).forEach(function(key) {
		// 	ExVars._setHeaderScore("." + key + "-header", scores[key].outlierScore)
		// })

		const exinds = {}
		scores.ExVars.data.forEach(function(exind) {
			if (exind[1] === 0) {
				return // This means it doesn't have an index
			}
			ExVars._setExcitabilityVariableScore("qtrac-excite-" + exind[1], exind[0])
		})
		values.sections.ExVars.data.forEach(function(exind) {
			const idx = exind[0]
			if (idx === 0) {
				return // This means it doesn't have an idx
			}
			ExVars._setExcitabilityVariableValue("qtrac-excite-" + idx, exind[1])
		})
	}
}
