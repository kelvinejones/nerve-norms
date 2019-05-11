function initPlots(participants, norms, outlierScores) {
	const osAccessor = {
		participant: "",
		normative: "",
		getScores: function() {
			return outlierScores[this.participant][this.normative]
		},
	}

	const partDropDown = new DataDropDown("select-participant-dropdown", participants, function(name, currentParticipant) {
		osAccessor.participant = name
		ExVars.update(osAccessor.getScores(), currentParticipant);
		plots.forEach(pl => {
			pl.chart.updateParticipant(currentParticipant.plots)
		})
	})

	const normDropDown = new DataDropDown("select-normative-dropdown", norms, function(name, currentNormative) {
		osAccessor.normative = name
		ExVars.update(osAccessor.getScores(), partDropDown.data);
		plots.forEach(pl => {
			pl.chart.updateNorms(currentNormative.plots)
		})
	})

	osAccessor.participant = partDropDown.selection
	osAccessor.normative = normDropDown.selection

	// Create all of the plots
	const plots = [{
		chart: new RecoveryCycle(partDropDown.data.plots, normDropDown.data.plots),
		selector: "#recoveryCycle svg",
	}, {
		chart: new ThresholdElectrotonus(partDropDown.data.plots, normDropDown.data.plots),
		selector: "#thresholdElectrotonus svg",
	}, {
		chart: new ChargeDuration(partDropDown.data.plots, normDropDown.data.plots),
		selector: "#chargeDuration svg",
	}, {
		chart: new ThresholdIV(partDropDown.data.plots, normDropDown.data.plots),
		selector: "#thresholdIV svg",
	}, {
		chart: new StimulusResponse(partDropDown.data.plots, normDropDown.data.plots),
		selector: "#stimulusResponse svg",
	}, {
		chart: new StimulusRelative(partDropDown.data.plots, normDropDown.data.plots),
		selector: "#stimulusResponseRelative svg",
	}, ]

	// Draw them all
	plots.forEach(pl => {
		pl.chart.draw(d3.select(pl.selector), true)
		pl.chart.setDelayTime(Chart.fastDelay) // After initial setup, remove the delay
	})

	// Now set all excitability variables
	ExVars.update(osAccessor.getScores(), partDropDown.data);
}
