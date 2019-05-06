function initPlots(data) {
	// Create all of the plots
	const plots = [{
		chart: new RecoveryCycle(data.plots),
		selector: "#recoveryCycle svg",
	}, {
		chart: new ThresholdElectrotonus(data.plots),
		selector: "#thresholdElectrotonus svg",
	}, {
		chart: new ChargeDuration(data.plots),
		selector: "#chargeDuration svg",
	}, {
		chart: new ThresholdIV(data.plots),
		selector: "#thresholdIV svg",
	}, {
		chart: new StimulusResponse(data.plots),
		selector: "#stimulusResponse svg",
	}, {
		chart: new StimulusRelative(data.plots),
		selector: "#stimulusResponseRelative svg",
	}, ]

	// Draw them all
	plots.forEach(pl => {
		pl.chart.draw(d3.select(pl.selector), true)
		pl.chart.setDelayTime(Chart.fastDelay) // After initial setup, remove the delay
	})

	// Now set all excitability variables
	updateIndices(data);

	function changeParticipant(ev) {
		plots.forEach(pl => {
			currentParticipant = participants[ev.srcElement.value] // This is a global
			pl.chart.updatePlots(currentParticipant.plots)
			updateIndices(currentParticipant)
		})
	}

	document.getElementById("select-participant-dropdown")
		.addEventListener("change", changeParticipant);
}

function updateIndices(data) {
	const interpolate = function() {
		let opacity = 0.8,
			red = d3.hsl("red"),
			green = d3.hsl("lightgreen");

		red.opacity = opacity
		green.opacity = opacity

		const intr = d3.interpolateHsl(green, red);
		return function(score) { return intr(Math.pow(score, 3)) }
	}()

	function setExcitabilityVariable(idString, value, score) {
		const row = document.getElementById(idString);
		row.getElementsByClassName("excite-value")[0].innerHTML = value

		if (score !== undefined) {
			percent = score * 100
			color = interpolate(score)
			row.style.background = "linear-gradient(to right, " + color + " " + percent + "%, #ffffff 0%)"
		}
	}

	function setHeaderScore(str, score) {
		[...document.querySelectorAll(str)].forEach(elm => {
			elm.style.background = interpolate(score);
		})
	}

	function labelForScore(score) {
		if (score < 0.75) {
			return "Healthy"
		} else if (score < .95) {
			return "Abnormal"
		} else {
			return "Extremely Abnormal"
		}
	}

	setHeaderScore(".participant-header", data.outlierScore)
	const healthLabel = labelForScore(data.outlierScore)
	setExcitabilityVariable("overall-score", healthLabel + " (" + data.outlierScore.toFixed(2) + ")", 0)
	var nameSpan = document.getElementById("participant-name");
	nameSpan.innerHTML = data.participant + " (" + healthLabel + ")"


	Object.keys(data.plots).map(function(key) {
			setHeaderScore("." + key + "-header", data.plots[key].outlierScore)
			return data.plots[key].discreteMeasures;
		}).flat()
		.concat(data.discreteMeasures)
		.forEach(function(exind) {
			setExcitabilityVariable("qtrac-excite-" + exind.qtracExciteID, exind.value, exind.outlierScore)
		})
}
