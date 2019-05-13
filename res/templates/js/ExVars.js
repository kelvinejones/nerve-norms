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

	static update(data) {
		ExVars._setHeaderScore(".participant-header", data.outlierScore)
		const healthLabel = ExVars._labelForScore(data.outlierScore)
		ExVars._setExcitabilityVariable("overall-score", healthLabel + " (" + data.outlierScore.toFixed(2) + ")", 0)
		const nameSpan = document.getElementById("participant-name");
		nameSpan.innerHTML = data.participant + " (" + healthLabel + ")"

		Object.keys(data.plots).map(function(key) {
				ExVars._setHeaderScore("." + key + "-header", data.plots[key].outlierScore)
				return data.plots[key].discreteMeasures;
			}).flat()
			.concat(data.discreteMeasures)
			.forEach(function(exind) {
				ExVars._setExcitabilityVariable("qtrac-excite-" + exind.qtracExciteID, exind.value, exind.outlierScore)
			})
	}
}
