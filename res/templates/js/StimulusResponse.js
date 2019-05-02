class StimulusResponse extends Chart {
	constructor(data) {
		super([0.01, 20], [0.01, 20], d3.scaleLog(), d3.scaleLog())
		this.data = data
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
		this.animateXYLineWithMean(this.data)
	}
}
