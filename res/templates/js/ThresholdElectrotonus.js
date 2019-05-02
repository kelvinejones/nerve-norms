class ThresholdElectrotonus extends Chart {
	constructor(plots) {
		super([0, 200], [-150, 100])

		this.hy40 = plots.teh40.data
		this.de40 = plots.ted40.data
		this.hy20 = plots.teh20.data
		this.de20 = plots.ted20.data
	}

	get name() { return "Threshold Electrotonus" }
	get xLabel() { return "Threshold Reduction (%)" }
	get yLabel() { return "Delay (ms)" }

	drawLines(svg) {
		this.animateXYLineWithMean(this.hy40)
		this.animateXYLineWithMean(this.de40)
		this.animateXYLineWithMean(this.hy20)
		this.animateXYLineWithMean(this.de20)
		this.drawHorizontalLine(this.linesLayer, 0)
	}
}
