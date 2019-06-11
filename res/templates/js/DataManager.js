class DataManager {
	// data is an object indexed by the drop-down's values
	// The dataUsers is expected to provide a list of objects that implement 'updateParticipant' and 'updateNorms'
	constructor(data, dataUsers) {
		this.dt = data
		this.dataUsers = dataUsers
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

		this.filter = new Filter(norms => {
			this.normData = norms
			Object.values(this.dataUsers()).forEach(pl => {
				pl.updateNorms(norms)
			})
		})

		const updateData = (ev) => {
			this.val = ev.srcElement.value
			this._updateParticipant()
		}

		const dropDown = document.getElementById("select-participant-dropdown")
		dropDown.addEventListener("change", updateData)
		dropDownOptions.forEach(function(opt) {
			dropDown.options[dropDown.options.length] = new Option(opt)
		})
		this.val = dropDown.value

		this.filter.update(this.val)
	}

	_updateParticipant() {
		const participantData = this.dt[this.val]

		ExVars.clearScores()

		Object.values(this.dataUsers()).forEach(pl => {
			pl.updateParticipant(participantData)
		})

		ExVars.updateValues(participantData)
		this.filter.update(this.val)
	}

	get norms() {
		return this.normData
	}

	get participantName() {
		return this.val
	}

	get participantData() {
		return this.dt[this.val]
	}
}
