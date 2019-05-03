class ThresholdElectrotonus extends Chart {
	constructor(plots) {
		super([0, 200], [-150, 100])
		this.data = plots.te.data
	}

	get name() { return "Threshold Electrotonus" }
	get xLabel() { return "Threshold Reduction (%)" }
	get yLabel() { return "Delay (ms)" }

	updatePlots(plots) {
		this.data = plots.te.data
		this.animateLines()
	}

	drawLines(svg) {
		Object.keys(this.data).forEach(key => {
			this.createXYLineWithMean(this.data[key], key)
		})
		this.drawHorizontalLine(this.linesLayer, 0)
		this.animateLines()
	}

	animateLines() {
		Object.keys(this.data).forEach(key => {
			this.animateXYLineWithMean(this.data[key], key)
		})
	}
}
