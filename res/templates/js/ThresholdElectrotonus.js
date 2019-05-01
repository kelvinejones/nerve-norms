class ThresholdElectrotonus extends Chart {
	constructor(hy40, de40, hy20, de20) {
		super()

		this.hy40 = hy40
		this.de40 = de40
		this.hy20 = hy20
		this.de20 = de20

		this.xscale = this.xscale.domain([1, 200]);
		this.yscale = this.yscale.domain([-150, 100]);
	}

	get name() { return "Threshold Electrotonus" }
	get xLabel() { return "Threshold Reduction (%)" }
	get yLabel() { return "Delay (ms)" }

	drawLines(svg) {
		this.animateXYLineWithMean(this.hy40, Chart.defaultConfig())
		this.animateXYLineWithMean(this.de40, Chart.defaultConfig())
		this.animateXYLineWithMean(this.hy20, Chart.defaultConfig())
		this.animateXYLineWithMean(this.de20, Chart.defaultConfig())
	}
}
