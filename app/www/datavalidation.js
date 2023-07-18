
function is_valid(str) {
	console.log(str);
	if (str=="") {return true;}
	let hasDot = false;
	for (let i = 0; i < str.length; i++) {
		if (str[i] < '0' || str[i] > '9') {
			if (str[i] == '.' && !hasDot) {
				hasDot = true;
				continue;
			}
			return false;
		}
	}
	if (parseFloat(str) > 4 || parseFloat(str) < 0) {
		return false;
	}
	return true;
}

/**
 *
 * @param {KeyboardEvent} event
 * @returns
 */
function data_validation_up(event) {
	//console.log(event);		//.target .key
	if (event.target.classList.contains("feedback-data__input")) {
		if (is_valid(event.target.value)) {
			if (event.target.parentElement.classList.contains("feedback-data__cell--invalid")) {event.target.parentElement.classList.remove("feedback-data__cell--invalid");}
			return true;
		}
		//event.target.classList.add("feedback-data__cell--invalid");
		event.target.parentElement.classList.add("feedback-data__cell--invalid");
		console.log(event);
		//event.querySelector(':invalid');
		return false;
	}
	return true;
}

function count_decimals(str) {
	if (str=="") {return 0;}
	let count = 0;
	for (let i = 0; i < str.length; i++) {
		if (str[i] == '.') {
			count++;
		}
	}
	return count;
}

function data_validation_down(event) {
	//console.log(event);		//.target .key
	if (event.target.classList.contains("feedback-data__input")) {
		if ((event.key == '.' && count_decimals(event.target.value) == 0) || event.keyCode == '8'
				|| event.keyCode == '9' || event.keyCode == '37' || event.keyCode == '39'
				|| (event.key >= '0' && event.key <= '9')) {
			return true;
		}
		event.preventDefault();
		return false;
	}
	return true;
}

function update_average() {

}

window.addEventListener("load", () => {
	document.querySelector(".feedback-data__score-table").addEventListener("keydown", data_validation_down);
	document.querySelector(".feedback-data__score-table").addEventListener("keyup", data_validation_up);
})
