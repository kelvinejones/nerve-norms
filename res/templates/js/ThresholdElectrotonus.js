class ThresholdElectrotonus extends Chart {
	constructor(hy40, de40, hy20, de20) {
		super([0, 200], [-150, 100])

		this.hy40 = hy40
		this.de40 = de40
		this.hy20 = hy20
		this.de20 = de20
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
