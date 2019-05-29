class StimulusResponse extends Chart {
	constructor(participant, norms) {
		super([1, 20], [0.01, 20], Chart.scaleType.LOG, Chart.scaleType.LOG)
		this.participant = this.calculateParticipant(participant)
		this.norms = norms.sections.SR.data
		this.xName = 'x'
		this.yName = 'y'
		this.xMeanName = 'meanX'
		this.yMeanName = 'meanY'
		this.xSDName = 'SDlogX'
		this.ySDName = 'SDlogY'
		this.sdFunc = Chart.logSD
	}

	calculateParticipant(participant) {
		let peakResponse = 0;
		let stimFor50PercentMax = 0;
		participant.exVars.forEach(function(exvar) {
			if (exvar.id == 1) {
				stimFor50PercentMax = exvar.value;
			} else if (exvar.id == 6) {
				peakResponse = exvar.value;
			}
		})

		return participant.sections.SR.data.map((d, i) => {
			return {
				'y': (i + 1) * 2 / 100 * peakResponse,
				'x': d.valueX * stimFor50PercentMax / 100,
			}
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
		this.norms = norms.sections.SR.data
		this.animateNorms(this.norms, "sr")
	}

	drawLines(svg) {
		this.createXYLine(this.participant, "sr")
		this.createNorms(this.norms, "sr")
		this.animateXYLine(this.participant, "sr")
		this.animateNorms(this.norms, "sr")
	}
}
