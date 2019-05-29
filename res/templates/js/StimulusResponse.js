class StimulusResponse extends Chart {
	constructor(participant, norms) {
		super([1, 20], [0.01, 20], Chart.scaleType.LOG, Chart.scaleType.LOG)
		this.participant = participant.sections.SR.data
		this.norms = norms.sections.SR.data
		this.xName = 'valueX'
		this.yName = 'valueY'
		this.xMeanName = 'meanX'
		this.yMeanName = 'meanY'
		this.xSDName = 'SDlogX'
		this.ySDName = 'SDlogY'
		this.sdFunc = Chart.logSD
	}

	get name() { return "Stimulus Response" }
	get xLabel() { return "Stimulus Current (mA)" }
	get yLabel() { return "Peak Response (mV)" }

	updateParticipant(participant) {
		this.participant = participant.sections.SR.data
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
