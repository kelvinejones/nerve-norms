function initPlots(participants, norms, outlierScores) {
	const osAccessor = {
		participant: "",
		getScores: function() {
			return outlierScores[this.participant]
		},
	}

	const partDropDown = new DataDropDown("select-participant-dropdown", participants, function(name, currentParticipant) {
		osAccessor.participant = name
		ExVars.update(osAccessor.getScores());
		plots.forEach(pl => {
			pl.chart.updatePlots(currentParticipant.plots)
			pl.chart.updateNorms(norms[name].plots)
		})
	})

	osAccessor.participant = partDropDown.selection

	// Create all of the plots
	const plots = [{
		chart: new RecoveryCycle(partDropDown.data.plots, norms[partDropDown.selection].plots),
		selector: "#recoveryCycle svg",
	}, {
		chart: new ThresholdElectrotonus(partDropDown.data.plots, norms[partDropDown.selection].plots),
		selector: "#thresholdElectrotonus svg",
	}, {
		chart: new ChargeDuration(partDropDown.data.plots, norms[partDropDown.selection].plots),
		selector: "#chargeDuration svg",
	}, {
		chart: new ThresholdIV(partDropDown.data.plots, norms[partDropDown.selection].plots),
		selector: "#thresholdIV svg",
	}, {
		chart: new StimulusResponse(partDropDown.data.plots, norms[partDropDown.selection].plots),
		selector: "#stimulusResponse svg",
	}, {
		chart: new StimulusRelative(partDropDown.data.plots, norms[partDropDown.selection].plots),
		selector: "#stimulusResponseRelative svg",
	}, ]

	// Draw them all
	plots.forEach(pl => {
		pl.chart.draw(d3.select(pl.selector), true)
		pl.chart.setDelayTime(Chart.fastDelay) // After initial setup, remove the delay
	})

	// Now set all excitability variables
	ExVars.update(osAccessor.getScores());
}
