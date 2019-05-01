class ChargeDuration extends Chart {
	constructor(data) {
		super()
		this.data = data
		this.xscale = this.xscale.domain([0, 1]);
		this.yscale = this.yscale.domain([0, 10]);
	}

	get name() { return "Charge Duration" }
	get xLabel() { return "Stimulus Width (ms)" }
	get yLabel() { return "Threshold Change (mA•ms)" }

	drawLines(svg) {
		let config = Chart.defaultConfig()
		config.xName = 'stimWidth'
		this.animateXYLineWithMean(this.data, config)
	}
}
