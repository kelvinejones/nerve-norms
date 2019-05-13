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

	static _setExcitabilityVariable(idString, value, score) {
		const row = document.getElementById(idString);
		row.getElementsByClassName("excite-value")[0].innerHTML = value
		if (score !== undefined) {
			row.style.background = "linear-gradient(to right, " + ExVars._interpolate(score) + " " + score * 100 + "%, #ffffff 0%)"
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
			return "Abnormal"
		} else {
			return "Extremely Abnormal"
		}
	}

	static update(scores, values) {
		ExVars._setHeaderScore(".participant-header", scores.outlierScore)
		const healthLabel = ExVars._labelForScore(scores.outlierScore)
		ExVars._setExcitabilityVariable("overall-score", healthLabel + " (" + scores.outlierScore.toFixed(2) + ")", 0)
		const nameSpan = document.getElementById("participant-name");
		nameSpan.innerHTML = scores.participant + " (" + healthLabel + ")"

		const exinds = {}
		Object.keys(scores.plots).map(function(key) {
				ExVars._setHeaderScore("." + key + "-header", scores.plots[key].outlierScore)
				return scores.plots[key].discreteMeasures;
			}).flat()
			.concat(scores.discreteMeasures)
			.forEach(function(exind) {
				exinds[exind.id] = { score: exind.outlierScore }
			})
		Object.keys(values.plots).map(function(key) {
				return values.plots[key].discreteMeasures;
			}).flat()
			.concat(values.discreteMeasures)
			.forEach(function(exind) {
				exinds[exind.id].value = exind.value
			})

		Object.keys(exinds).forEach(function(id) {
			ExVars._setExcitabilityVariable("qtrac-excite-" + id, exinds[id].value, exinds[id].score)
		})
	}
}
