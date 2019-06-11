class ParticipantDropDown {
	// elementID is the HTML ID.
	// data is an object indexed by the drop-down's values
	// The dataProvider is expected to provide a list of objects that implement 'updateParticipant'
	constructor(elementID, data, filter, dataProvider, initialOptions) {
		this.dt = data

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

		const dropDown = document.getElementById(elementID)
		dropDown.addEventListener("change", updateData)
		initialOptions.forEach(function(opt) {
			dropDown.options[dropDown.options.length] = new Option(opt)
		})
		this.val = dropDown.value

		filter.update(this.val)
	}

	get name() {
		return this.val
	}

	get data() {
		return this.dt[this.val]
	}
}
