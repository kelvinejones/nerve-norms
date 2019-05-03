class ChargeDuration extends Chart {
	constructor(plots) {
		super([0, 1], [0, 10])
		this.data = plots.cd.data
		this.xName = 'stimWidth'
	}

	get name() { return "Charge Duration" }
	get xLabel() { return "Stimulus Width (ms)" }
	get yLabel() { return "Threshold Change (mA•ms)" }

	updatePlots(plots) {
		this.data = plots.cd.data
		this.animateXYLineWithMean(this.data, "cd", 0)
	}

	drawLines(svg) {
		this.createXYLineWithMean(this.data, "cd")
		this.animateXYLineWithMean(this.data, "cd")
	}
}
