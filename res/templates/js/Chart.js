class Chart {
	constructor(xRange, yRange, xScaleType = Chart.scaleType.LINEAR, yScaleType = Chart.scaleType.LINEAR) {
		this.xRange = xRange
		this.yRange = yRange

		this.margin = { top: 20, right: 20, bottom: 50, left: 50 };
		this.width = 600 - this.margin.left - this.margin.right;
		this.height = 300 - this.margin.top - this.margin.bottom;

		// Set the default scaling
		this.xscale = this.makeScale(xScaleType).range([0, this.width]).domain(xRange);
		this.yscale = this.makeScale(yScaleType).range([this.height, 0]).domain(yRange);

		// Set default config values
		this.xName = 'delay'
		this.yName = 'value'
		this.ySDName = 'SD'
		this.yMeanName = 'mean'

		this.yAnimStart = this.animationStartValue(this.yRange)

		this.delayTime = Chart.slowDelay;
		this.transitionTime = Chart.slowTransition;

		this.group = {}

		this.numSD = 2
		this.sdFunc = Chart.linearSD
	}

	makeScale(name) {
		let scale
		switch (name) {
			case "LINEAR":
				scale = d3.scaleLinear()
				scale.scaleType = Chart.scaleType.LINEAR
				break
			case "LOG":
				scale = d3.scaleLog()
				scale.scaleType = Chart.scaleType.LOG
				break
		}
		return scale
	}

	animationStartValue(range) {
		// The animation start value should be at zero unless that's outside of range
		if (range[0] < 0 && range[1] >= 0) {
			return 0
		} else {
			return range[0]
		}
	}

	updateParticipant(participant) { throw new Error("A Chart must implement updateParticipant(participant)") }
	updateNorms(norms) { throw new Error("A Chart must implement updateNorms(norms)") }
	get name() { throw new Error("A Chart must implement name()") }
	get xLabel() { throw new Error("A Chart must implement xLabel()") }
	get yLabel() { throw new Error("A Chart must implement yLabel()") }
	drawLines(svg) { throw new Error("A Chart must implement drawLines(svg)") }

	draw(svg, hideLabels) {
		// append the svg object to the body of the page
		// appends a 'group' element to 'svg'
		// moves the 'group' element to the top left margin
		svg = svg
			.append("g")
			.attr("transform",
				"translate(" + this.margin.left + "," + this.margin.top + ")");

		// Add layers for various elements
		this.ciLayer = svg.append("g")
		this.meanLayer = svg.append("g")
		this.linesLayer = svg.append("g")
		this.valueLayer = svg.append("g")
		this.circlesLayer = svg.append("g")

		// Add the X Axis
		const xelements = svg.append("g")
			.attr("transform", "translate(0," + this.height + ")")
			.call(d3.axisBottom(this.xscale).ticks(2)
				.tickFormat(d3.format("")));

		// Add the Y Axis
		const yelements = svg.append("g")
			.call(d3.axisLeft(this.yscale));

		if (!hideLabels) {
			svg.attr("transform", "scale(0.9) translate(110, 0)")
			this.labels(svg);

		} else {
			xelements.selectAll("text").remove();
			yelements.selectAll("text").remove();
		}

		this.drawLines(svg)
	}

	labels(svg) {
		// text label for the x axis
		svg.append("text")
			.attr("transform",
				"translate(" + (this.width / 2) + " ," +
				(this.height + this.margin.top + 20) + ")")
			.style("text-anchor", "middle")
			.text(this.xLabel);

		// text label for the y axis
		svg.append("text")
			.attr("transform", "rotate(-90)")
			.attr("y", 0 - this.margin.left)
			.attr("x", 0 - (this.height / 2))
			.attr("dy", "1em")
			.style("text-anchor", "middle")
			.text(this.yLabel);
	}

	static linearSD(val, sign, numSD, sdSize) {
		return val + sign * numSD * (sdSize || 0)
	}

	static logSD(val, sign, numSD, sdSize) {
		return val * Math.pow(10, sign * numSD * (sdSize || 0))
	}

	// sdAtLoc calculates limits based on the standard deviations
	sdAtLoc(dpt, loc) {
		const xmn = this.xMeanName || this.xName
		const ymn = this.yMeanName || this.yName
		const xsd = this.xSDName
		const ysd = this.ySDName
		const numSD = this.numSD
		const sdFunc = this.sdFunc

		function stPoint(xSign, ySign) {
			let scale = numSD
			if (dpt[xsd] !== undefined && xSign != 0 && dpt[ysd] !== undefined && ySign != 0) {
				// Since both are set, scale the edges by sqrt(2) to make a ovoid area
				scale = 0.707 * numSD
			}
			return { x: sdFunc(dpt[xmn], xSign, scale, dpt[xsd]), y: sdFunc(dpt[ymn], ySign, scale, dpt[ysd]) }
		}

		switch (loc) {
			case Chart.limLoc.UPPER_LEFT:
				return stPoint(-1, 1)
			case Chart.limLoc.UPPER_RIGHT:
				return stPoint(1, 1)
			case Chart.limLoc.LOWER_LEFT:
				return stPoint(-1, -1)
			case Chart.limLoc.LOWER_RIGHT:
				return stPoint(1, -1)
			case Chart.limLoc.UPPER:
				return stPoint(0, 1)
			case Chart.limLoc.LOWER:
				return stPoint(0, -11)
			case Chart.limLoc.LEFT:
				return stPoint(-1, 0)
			case Chart.limLoc.RIGHT:
				return stPoint(1, 0)
		}
	}

	normativeLimits(data) {
		if (this.sdNum == 0) {
			return []
		}

		// this.sdAtLoc can be changed if a different calculation is more appropriate.
		return this.scaleArrayWithinRange((Array.from(data)
				.map(d => { return this.sdAtLoc(d, Chart.limLoc.UPPER_LEFT) }))
			.concat(this.sdAtLoc(data[data.length - 1], Chart.limLoc.UPPER))
			.concat(this.sdAtLoc(data[data.length - 1], Chart.limLoc.UPPER_RIGHT))
			.concat(this.sdAtLoc(data[data.length - 1], Chart.limLoc.RIGHT))
			.concat(Array.from(data).reverse().map(d => { return this.sdAtLoc(d, Chart.limLoc.LOWER_RIGHT) }))
			.concat(this.sdAtLoc(data[0], Chart.limLoc.LOWER))
			.concat(this.sdAtLoc(data[0], Chart.limLoc.LOWER_LEFT))
			.concat(this.sdAtLoc(data[0], Chart.limLoc.LEFT)))
	}

	scaleArrayWithinRange(ar) {
		// With a log scale, values can't be plotted at or below zero.
		if (this.xscale.scaleType == Chart.scaleType.LOG) {
			ar = this.raiseZeroValues(ar, 'x', this.xRange[0])
		}
		if (this.yscale.scaleType == Chart.scaleType.LOG) {
			ar = this.raiseZeroValues(ar, 'y', this.yRange[0])
		}
		return ar
	}

	raiseZeroValues(ar, axis, min) {
		return ar.map(function(d) {
			if (d[axis] < min) {
				d[axis] = min
			}
			return d
		})
	}

	dataAsXY(data, xName, yName) {
		if (data[0][xName] === undefined || data[0][yName] === undefined) {
			return []
		}
		return data.map(function(d) { return { x: d[xName], y: d[yName] } })
	}

	xZeroPath() {
		return d3.line()
			.x(d => this.xscale(d.x))
			.y(d => this.yscale(this.yAnimStart))
	}

	xyPath() {
		return d3.line()
			.x(d => this.xscale(d.x))
			.y(d => this.yscale(d.y));
	}

	drawHorizontalLine(svg, yVal) {
		svg.append("path")
			.data([
				[{ x: this.xRange[0], y: yVal }, { x: this.xRange[1], y: yVal }]
			])
			.attr("class", "meanline")
			.attr("d", this.xyPath());
	}

	drawVerticalLine(svg, xVal) {
		svg.append("path")
			.data([
				[{ y: this.yRange[0], x: xVal }, { y: this.yRange[1], x: xVal }]
			])
			.attr("class", "meanline")
			.attr("d", this.xyPath());
	}

	createGroup(svg, typeString, name) {
		svg = svg.append("g")
		this.group[typeString + "-" + name] = svg
		return svg
	}

	animateGroup(typeString, newData, name) {
		return this.group[typeString + "-" + name].selectAll(typeString)
			.data(newData)
			.transition()
			.delay(this.delayTime)
			.duration(this.transitionTime)
	}

	createPath(svg, path, groupName, className) {
		this.createGroup(svg, "path", groupName + "-" + className)
			.append("path")
			.data([path])
			.attr("class", className)
			.attr("d", this.xZeroPath())
	}

	animatePath(path, groupName, className) {
		this.animateGroup("path", [path], groupName + "-" + className)
			.attr("d", this.xyPath());
	}

	createCircles(svg, circleLocations, name) {
		// create circle locations at init for each name, with right amount and position of circles
		// Add circles into a separate SVG group
		this.createGroup(svg, "circle", name)
			.selectAll("circle")
			.data(circleLocations)
			.enter()
			.append("circle")
			.attr("cx", d => this.xscale(d[this.xName]))
			.attr("cy", this.yscale(this.yAnimStart))
			.attr("r", d => 3)
			.style("fill", d => "black");
	}

	fillColor(pt) {
		// If either the x or y value is undefined, then this point should be hidden
		if (pt[this.yName] === undefined || pt[this.yName] === undefined) {
			return "rgba(0, 0, 0, 0)"
		}
		return pt.wasImputed ? "red" : "black"
	}

	animateCircles(circleLocations, name) {
		this.animateGroup("circle", circleLocations, name)
			.attr("r", d => d.wasImputed ? 3 : 5)
			.style("fill", d => this.fillColor(d))
			.attr("cy", d => this.yscale(d[this.yName] || 0))
			.attr("cx", d => this.xscale(d[this.xName] || 0))
	}

	createNorms(lineData, name) {
		this.createPath(this.ciLayer, this.normativeLimits(lineData), name, "confidenceinterval")
		this.createPath(this.meanLayer, this.dataAsXY(lineData, this.xMeanName || this.xName, this.yMeanName || this.yName), name, "meanline")
	}

	createXYLine(lineData, name) {
		this.createPath(this.valueLayer, this.dataAsXY(lineData, this.xName, this.yName), name, "line")
		this.createCircles(this.circlesLayer, lineData, name)
	}

	animateNorms(lineData, name) {
		this.animatePath(this.normativeLimits(lineData), name, "confidenceinterval")
		this.animatePath(this.dataAsXY(lineData, this.xMeanName || this.xName, this.yMeanName || this.yName), name, "meanline")
	}

	animateXYLine(lineData, name) {
		this.animatePath(this.dataAsXY(lineData, this.xName, this.yName), name, "line")
		this.animateCircles(lineData, name)
	}

	setDelayTime(dt) {
		// Return 'this' for chaining
		this.delayTime = dt
		return this
	}

	setTransitionTime(tt) {
		// Return 'this' for chaining
		this.transitionTime = tt
		return this
	}
}

// Set some constants for the class
Object.defineProperty(Chart, 'slowDelay', {
	value: 750,
	enumerable: true,
})
Object.defineProperty(Chart, 'fastDelay', {
	value: 0,
	enumerable: true,
})
Object.defineProperty(Chart, 'slowTransition', {
	value: 2500,
	enumerable: true,
})
Object.defineProperty(Chart, 'fastTransition', {
	value: 1000,
	enumerable: true,
})
Object.defineProperty(Chart, 'scaleType', {
	value: {
		LINEAR: "LINEAR",
		LOG: "LOG",
	},
	enumerable: true,
})
Object.defineProperty(Chart, 'limLoc', {
	value: {
		UPPER_LEFT: "UPPER_LEFT",
		UPPER_RIGHT: "UPPER_RIGHT",
		LOWER_LEFT: "LOWER_LEFT",
		LOWER_RIGHT: "LOWER_RIGHT",
		UPPER: "UPPER",
		LOWER: "LOWER",
		LEFT: "LEFT",
		RIGHT: "RIGHT",
	},
	enumerable: true,
})
