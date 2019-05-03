class StimulusResponse extends Chart {
	constructor(plots) {
		super([0.01, 20], [0.01, 20], Chart.scaleType.LOG, Chart.scaleType.LOG)
		this.data = plots.sr.data
		this.xName = 'valueX'
		this.yName = 'valueY'
		this.xMeanName = 'meanX'
		this.yMeanName = 'meanY'
		this.xSDName = 'SDX'
		this.ySDName = 'SDY'
	}

	get name() { return "Stimulus Response" }
	get xLabel() { return "Stimulus Current (mA)" }
	get yLabel() { return "Peak Response (mV)" }

	drawLines(svg) {
		this.animateXYLineWithMean(this.data, "sr")
	}
}
