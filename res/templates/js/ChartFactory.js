class ChartFactory {
	constructor(participants, norms, outlierScores) {
		this.osAccessor = {
			participant: "",
			normative: "",
			getScores: function() {
				return outlierScores[this.participant][this.normative]
			},
		}

		this.partDropDown = new DataDropDown("select-participant-dropdown", participants, (name, currentParticipant) => {
			this.osAccessor.participant = name
			ExVars.update(this.osAccessor.getScores(), currentParticipant);
			plots.forEach(pl => {
				pl.chart.updateParticipant(currentParticipant.plots)
			})
		})

		this.normDropDown = new DataDropDown("select-normative-dropdown", norms, (name, currentNormative) => {
			this.osAccessor.normative = name
			ExVars.update(this.osAccessor.getScores(), this.partDropDown.data);
			plots.forEach(pl => {
				pl.chart.updateNorms(currentNormative.plots)
			})
		})

		this.osAccessor.participant = this.partDropDown.selection
		this.osAccessor.normative = this.normDropDown.selection

		// Create all of the plots
		const plots = [{
			chart: this.build("recoveryCycle"),
			selector: "#recoveryCycle svg",
		}, {
			chart: this.build("thresholdElectrotonus"),
			selector: "#thresholdElectrotonus svg",
		}, {
			chart: this.build("chargeDuration"),
			selector: "#chargeDuration svg",
		}, {
			chart: this.build("thresholdIV"),
			selector: "#thresholdIV svg",
		}, {
			chart: this.build("stimulusResponse"),
			selector: "#stimulusResponse svg",
		}, {
			chart: this.build("stimulusResponseRelative"),
			selector: "#stimulusResponseRelative svg",
		}, ]

		// Draw them all
		plots.forEach(pl => {
			pl.chart.draw(d3.select(pl.selector), true)
			pl.chart.setDelayTime(Chart.fastDelay) // After initial setup, remove the delay
		})

		// Now set all excitability variables
		ExVars.update(this.osAccessor.getScores(), this.partDropDown.data);
	}

	drawModal(typeStr) {
		const chart = this.build(typeStr)
		chart.setDelayTime(Chart.fastDelay).setTransitionTime(Chart.fastTransition)

		document.getElementById('modal-title').innerHTML = chart.name
		d3.selectAll("#modal svg > *").remove()

		chart.draw(d3.select('#modal svg'))
		$('#modal').modal('toggle')
	}


	build(typeStr) {
		switch (typeStr) {
			case "recoveryCycle":
				return new RecoveryCycle(this.partDropDown.data.plots, this.normDropDown.data.plots)
			case "thresholdElectrotonus":
				return new ThresholdElectrotonus(this.partDropDown.data.plots, this.normDropDown.data.plots)
			case "chargeDuration":
				return new ChargeDuration(this.partDropDown.data.plots, this.normDropDown.data.plots)
			case "thresholdIV":
				return new ThresholdIV(this.partDropDown.data.plots, this.normDropDown.data.plots)
			case "stimulusResponse":
				return new StimulusResponse(this.partDropDown.data.plots, this.normDropDown.data.plots)
			case "stimulusResponseRelative":
				return new StimulusRelative(this.partDropDown.data.plots, this.normDropDown.data.plots)
		}
	}
}
