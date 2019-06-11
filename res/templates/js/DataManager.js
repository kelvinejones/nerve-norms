class DataManager {
	// data is an object indexed by the drop-down's values
	// The dataUsers is expected to provide a list of objects that implement 'updateParticipant' and 'updateNorms'
	constructor(data, dataUsers) {
		this.dt = data
		this.dataUsers = dataUsers

		const uploadOption = "Upload MEM..."
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
			uploadOption,
		]

		this.filter = new Filter(norms => {
			this.normData = norms
			Object.values(this.dataUsers()).forEach(pl => {
				pl.updateNorms(norms)
			})
		})

		const updateData = (ev) => {
			this.val = ev.srcElement.value
			ExVars.clearScores()
			if (this.val == uploadOption) {
				this._uploadMEM()
			} else {
				this.uploadData = undefined
				this._updateParticipant(this.dt[this.val])
				this.filter.update(this.val)
			}
		}

		const dropDown = document.getElementById("select-participant-dropdown")
		dropDown.addEventListener("change", updateData)
		dropDownOptions.forEach(function(opt) {
			dropDown.options[dropDown.options.length] = new Option(opt)
		})
		this.val = dropDown.value

		this.filter.update(this.val)
	}

	_uploadMEM() {
		// This code is modified from https://stackoverflow.com/a/40971885
		var input = document.createElement('input');
		input.type = 'file';

		input.onchange = e => {
			var file = e.target.files[0];

			var reader = new FileReader();
			reader.readAsText(file, 'UTF-8');

			reader.onload = readerEvent => {
				var content = readerEvent.target.result; // this is the content!
				this.filter.fetchMEM(content, convertedMem => {
					this.uploadData = convertedMem.participant
					this._updateParticipant(this.uploadData)
					this.filter.setParticipantData(this.uploadData)
				})
			}
		}

		input.click();
	}

	_updateParticipant(participantData) {
		Object.values(this.dataUsers()).forEach(pl => {
			pl.updateParticipant(participantData)
		})

		ExVars.updateValues(participantData)
	}

	get norms() {
		return this.normData
	}

	get participantName() {
		return this.val
	}

	get participantData() {
		if (this.uploadData != null) {
			return this.uploadData
		} else {
			return this.dt[this.val]
		}
	}
}
