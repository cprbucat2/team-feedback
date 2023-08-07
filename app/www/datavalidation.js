
function is_valid(str) {
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
function data_validation_input(event) {
	if (event.target.classList.contains("feedback-data__input")) {
		if (is_valid(event.target.value)) {
			if (event.target.parentElement.classList.contains("feedback-data__cell--invalid")) {event.target.parentElement.classList.remove("feedback-data__cell--invalid");}
			return true;
		}
		//event.target.classList.add("feedback-data__cell--invalid");
		event.target.parentElement.classList.add("feedback-data__cell--invalid");
		//event.querySelector(':invalid');
	}
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
	}
	return true;
}

function data_validation_comment() {
	for (const cell of document.querySelectorAll(".feedback-comments__member-comments")) {
		data = cell.value;
		if (data!="" && typeof data!="undefined" && cell.parentElement.classList.contains("feedback-comments__member-comments--invalid")) {
			cell.parentElement.classList.remove("feedback-comments__member-comments--invalid");
		}
	}
}

function update_average() {

}

window.addEventListener("load", () => {
	document.querySelector(".feedback-data__score-table").addEventListener("keydown", data_validation_down);
	document.querySelector(".feedback-data__score-table").addEventListener("input", data_validation_input);
	var elements = document.querySelectorAll(".feedback-comments__table");
	elements.forEach(function(element) {
		element.addEventListener("input", data_validation_comment);
	});
})
