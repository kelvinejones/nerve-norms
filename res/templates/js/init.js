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
	let opacity = 0.8,
		red = d3.hsl("red"),
		green = d3.hsl("lightgreen");

	red.opacity = opacity
	green.opacity = opacity

	let interpolate = d3.interpolateHsl(green, red);

	function setExcitabilityVariable(idString, value, score) {
		var row = document.getElementById(idString);
		row.getElementsByClassName("excite-value")[0].innerHTML = value

		percent = score * 100
		color = interpolate(score)

		row.style.background = "linear-gradient(to right, " + color + " " + percent + "%, #ffffff 0%)"
	}

	setExcitabilityVariable("overall-score", data.outlierScore, data.outlierScore)

	Object.keys(data.plots).map(function(key) { return data.plots[key].discreteMeasures; }).flat()
		.concat(data.discreteMeasures)
		.forEach(function(exind) {
			setExcitabilityVariable("qtrac-excite-" + exind.qtracExciteID, exind.value, exind.outlierScore)
		})
}
