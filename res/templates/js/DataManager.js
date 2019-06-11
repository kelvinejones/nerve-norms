class DataManager {
	// data is an object indexed by the drop-down's values
	// The dataUsers is expected to provide a list of objects that implement 'updateParticipant' and 'updateNorms'
	constructor(data, dataUsers) {
		this.dt = data
		this.dataUsers = dataUsers
		this.uploadCount = 0

		this.participants = [
			Participant.load("CA-WI20S", data),
			Participant.load("CA-AL27H", data),
			Participant.load("JP-20-1", data),
			Participant.load("JP-70-1", data),
			Participant.load("PO-00d97e84", data),
			Participant.load("PO-017182a5", data),
			Participant.load("CA Mean", data),
			Participant.load("JP Mean", data),
			Participant.load("PO Mean", data),
			Participant.load("Rat Fast Axon", data),
			Participant.load("Rat Slow Axon", data),
			Participant.load("Rat on Drugs", data),
		]

		Filter.setCallback(() => this._fetchUpdates(this.filtername))

		this.dropDown = document.getElementById("select-participant-dropdown")
		this.dropDown.addEventListener("change", (ev) => {
			this.val = ev.srcElement.value
			ExVars.clearScores()
			if (this.val == DataManager.uploadOption) {
				this._uploadMEM()
			} else {
				this.uploadData = undefined
				this._updateParticipant(this.dt[this.val])
				this._fetchUpdates(this.val)
			}
		})
		this._updateDropDownOptions()

		this.val = this.dropDown.value

		this._fetchUpdates(this.val)
	}

	_fetchUpdates(name) {
		const lastQuery = this.queryString
		this.queryString = Filter.queryString
		const normChanged = (lastQuery != this.queryString)
		if (normChanged) {
			Fetch.Norms(this.queryString, norms => {
				this.normData = norms
				Object.values(this.dataUsers()).forEach(pl => {
					pl.updateNorms(norms)
				})
			})
		}

		if (name != undefined) {
			this.filterdata = undefined
		}
		const lastParticipant = this.filtername
		this.filtername = name
		const nameChanged = (lastParticipant != this.filtername)
		if (normChanged || nameChanged) {
			ExVars.clearScores()
			Fetch.Outliers(this.queryString, this.filtername, this.filterdata, scores => {
				ExVars.updateScores(this.filtername, scores)
			})
		}
	}

	static get uploadOption() { return "Upload MEM..." }

	_updateDropDownOptions() {
		const selection = this.dropDown.selectedIndex

		let index = 0;
		this.participants.forEach(opt => {
			this.dropDown.options[index++] = new Option(opt.name)
		})
		this.dropDown.options[index] = new Option(DataManager.uploadOption)

		if (selection >= 0) {
			this.dropDown.selectedIndex = selection
		}
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
				this.filterlastParticipant = undefined
				Fetch.MEM(this.queryString, content, convertedMem => {
					this.uploadData = convertedMem.participant
					this._updateParticipant(this.uploadData)
					this.filtername = undefined
					this.filterdata = this.uploadData

					this.uploadCount = this.uploadCount + 1
					this.participants[this.participants.length] = new Participant(this.uploadData, "Upload " + this.uploadCount + ": " + this.uploadData.header.name)
					this._updateDropDownOptions()

					ExVars.updateScores(this.filtername, convertedMem.outlierScores)
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
