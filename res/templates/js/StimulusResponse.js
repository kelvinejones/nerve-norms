class StimulusResponse extends Chart {
	constructor(plots) {
		super([1, 20], [0.01, 20], Chart.scaleType.LOG, Chart.scaleType.LOG)
		this.participant = plots.sr.data
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

	updatePlots(plots) {
		this.participant = plots.sr.data
		this.animateXYLine(this.participant, "sr")
	}

	updateNorms(plots) {
		this.animateNorms(this.participant, "sr")
	}

	drawLines(svg) {
		this.createXYLine(this.participant, "sr")
		this.createNorms(this.participant, "sr")
		this.animateXYLine(this.participant, "sr")
		this.animateNorms(this.participant, "sr")
	}
}
