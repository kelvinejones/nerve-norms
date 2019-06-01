class StimulusResponse extends Chart {
	constructor(participant, norms) {
		super([1, 20], [0.01, 20], Chart.scaleType.LOG, Chart.scaleType.LOG)
		this.participant = this.calculateParticipant(participant)
		this.norms = norms.SRY.data
		this.sdFunc = Chart.logSD
	}

	calculateParticipant(participant) {
		let peakResponse = 0;
		participant.exVars.forEach(function(exvar) {
			if (exvar.id == 6) {
				peakResponse = exvar.value;
			}
		})

		return participant.sections.SR.data.map((d, i) => {
			return [
				d[1],
				(i + 1) * 2 / 100 * peakResponse,
			]
		})
	}

	get name() { return "Stimulus Response" }
	get xLabel() { return "Stimulus Current (mA)" }
	get yLabel() { return "Peak Response (mV)" }

	updateParticipant(participant) {
		this.participant = this.calculateParticipant(participant)
		this.animateXYLine(this.participant, "sr")
	}

	updateNorms(norms) {
		this.norms = norms.SR.data
		// this.animateNorms(this.norms, "sr")
	}

	drawLines(svg) {
		this.createXYLine(this.participant, "sr")
		// this.createNorms(this.norms, "sr")
		this.animateXYLine(this.participant, "sr")
		// this.animateNorms(this.norms, "sr")
	}
}
