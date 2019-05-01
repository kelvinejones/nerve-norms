class StimulusResponse extends Chart {
	constructor(data) {
		super([0, 10], [0, 10])
		this.data = data
		this.xName = 'valueX'
		this.yName = 'valueY'
		this.xMeanName = 'meanX'
		this.yMeanName = 'meanY'
		this.ySDName = 'SDY'
	}

	get name() { return "Stimulus Response" }
	get xLabel() { return "Stimulus Current (mA)" }
	get yLabel() { return "Peak Response (mV)" }

	drawLines(svg) {
		this.animateXYLineWithMean(this.data)
	}
}
