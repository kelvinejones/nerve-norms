class ParticipantDropDown {
	// data is an object indexed by the drop-down's values
	// The dataProvider is expected to provide a list of objects that implement 'updateParticipant'
	constructor(data, dataProvider) {
		this.dt = data
		const dropDownOptions = [
			"CA-WI20S",
			"CA-AL27H",
			"JP-20-1",
			"JP-70-1",
			"PO-00d97e84",
			"PO-017182a5",
			"CA Mean",
			"JP Mean",
			"PO Mean",
			"Rat Fast Axon",
			"Rat Slow Axon",
			"Rat on Drugs",
		]

		const filter = new Filter(norms => {
			this.normData = norms
			Object.values(dataProvider()).forEach(pl => {
				pl.updateNorms(norms)
			})
		})

		const updateData = (ev) => {
			this.val = ev.srcElement.value
			const participantData = this.dt[this.val]

			ExVars.clearScores()

			Object.values(dataProvider()).forEach(pl => {
				pl.updateParticipant(participantData)
			})

			ExVars.updateValues(participantData)
			filter.update(this.val)
		}

		const dropDown = document.getElementById("select-participant-dropdown")
		dropDown.addEventListener("change", updateData)
		dropDownOptions.forEach(function(opt) {
			dropDown.options[dropDown.options.length] = new Option(opt)
		})
		this.val = dropDown.value

		filter.update(this.val)
	}

	get norms() {
		return this.normData
	}

	get name() {
		return this.val
	}

	get data() {
		return this.dt[this.val]
	}
}
