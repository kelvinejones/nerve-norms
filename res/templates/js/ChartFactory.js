class ChartFactory {
	constructor(participants) {
		this.dataManager = new DataManager(participants, () => { return this.plots })

		this.plots = {
			"recoveryCycle": null,
			"thresholdElectrotonus": null,
			"chargeDuration": null,
			"thresholdIV": null,
			"stimulusResponse": null,
			"stimulusResponseRelative": null,
		}

		Object.keys(this.plots).forEach(key => {
			this.plots[key] = this.build(key)
			this.plots[key].draw(d3.select("#" + key + " svg"), true)
			this.plots[key].setDelayTime(Chart.fastDelay) // After initial setup, remove the delay
		})

		// Now set all excitability variables
		ExVars.updateValues(this.dataManager.participantData)
	}

	/**
	 * Detect if the current active responsive breakpoint in Bootstrap is not XS
	 * Modified from code by farside {@link https://stackoverflow.com/users/4354249/farside}
	 */
	static bootstrapSizeIsXS() {
		const $el = $("<div>")
		$el.appendTo($("body"))
		$el.addClass("d-sm-none")
		const isBiggerThanXS = $el.is(":hidden")
		$el.remove()
		return !isBiggerThanXS
	};

	drawModal(typeStr) {
		if (ChartFactory.bootstrapSizeIsXS()) {
			// Screen is too small to display this properly
			return
		}
		const chart = this.build(typeStr)
		chart.setDelayTime(Chart.fastDelay).setTransitionTime(Chart.fastTransition)

		document.getElementById('modal-title').innerHTML = chart.name
		d3.selectAll("#modal svg > *").remove()

		chart.draw(d3.select('#modal svg'))
		$('#modal').modal('toggle')
	}

	build(typeStr) {
		let chart;
		switch (typeStr) {
			case "recoveryCycle":
				chart = new RecoveryCycle(this.dataManager.participantData, this.dataManager.norms)
				break
			case "thresholdElectrotonus":
				chart = new ThresholdElectrotonus(this.dataManager.participantData, this.dataManager.norms)
				break
			case "chargeDuration":
				chart = new ChargeDuration(this.dataManager.participantData, this.dataManager.norms)
				break
			case "thresholdIV":
				chart = new ThresholdIV(this.dataManager.participantData, this.dataManager.norms)
				break
			case "stimulusResponse":
				chart = new StimulusResponse(this.dataManager.participantData, this.dataManager.norms)
				break
			case "stimulusResponseRelative":
				chart = new StimulusRelative(this.dataManager.participantData, this.dataManager.norms)
				break
		}
		chart.updateScore(this.dataManager.outliers)
		return chart
	}
}
