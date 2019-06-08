class ChartFactory {
	constructor(participants) {
		this.url = "https://us-central1-nervenorms.cloudfunctions.net/"

		this.partDropDown = new DataDropDown("select-participant-dropdown", participants, (name, currentParticipant) => {
			this.participant = name

			ExVars.clearScores()

			Object.values(this.plots).forEach(pl => {
				pl.updateParticipant(currentParticipant)
			})
			ExVars.updateValues(currentParticipant)
			this.updateOutliers(this.participant)
		}, ["CA-CR21S", "CA-WI20S", "JP-10-1", "JP-70-1", "PO-00d97e84", "PO-017182a5", "Rat on Drugs"])
		this.participant = this.partDropDown.selection

		const queryString = Filter.asQueryString()
		this.updateNorms(queryString)
		this.updateOutliers(this.participant, queryString)

		document.querySelector("form").addEventListener("submit", (event) => {
			ExVars.clearScores()

			const queryString = Filter.asQueryString()
			this.updateNorms(queryString)
			this.updateOutliers(this.participant, queryString)

			event.preventDefault()
		})

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
		ExVars.updateValues(this.partDropDown.data)
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
		switch (typeStr) {
			case "recoveryCycle":
				return new RecoveryCycle(this.partDropDown.data, this.norms)
			case "thresholdElectrotonus":
				return new ThresholdElectrotonus(this.partDropDown.data, this.norms)
			case "chargeDuration":
				return new ChargeDuration(this.partDropDown.data, this.norms)
			case "thresholdIV":
				return new ThresholdIV(this.partDropDown.data, this.norms)
			case "stimulusResponse":
				return new StimulusResponse(this.partDropDown.data, this.norms)
			case "stimulusResponseRelative":
				return new StimulusRelative(this.partDropDown.data, this.norms)
		}
	}

	updateOutliers(name, queryString) {
		if (queryString === undefined) {
			queryString = Filter.asQueryString()
		}

		fetch(this.url + "outliers" + queryString + "&name=" + name)
			.then(response => {
				return response.json()
			})
			.then(scores => {
				this.scores = scores
				ExVars.updateScores(this.participant, scores)
			})
	}

	updateNorms(queryString) {
		if (queryString === undefined) {
			queryString = Filter.asQueryString()
		}

		fetch(this.url + "norms" + queryString)
			.then(response => {
				return response.json()
			})
			.then(norms => {
				this.norms = norms
				Object.values(this.plots).forEach(pl => {
					pl.updateNorms(norms)
				})
			})
	}
}
