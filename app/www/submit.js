/**
 * @file Handles submissions from index.html.
 * @author Aidan Hoover, Aiden Woodruff
 * @copyright 2023 Aidan Hoover and Aiden Woodruff
 * @license BSD-3-Clause
 */

/**
 * Check valid input in all cells
 * @param {str}
 * @returns if given str is a string of numeric floating point value between 0 and 4 inclusive
 */
function is_valid(str) {
	if (str=="" || typeof str=="undefined") {return false;}
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
 * Check valid input in all cells
 * @param {table}
 * @requires table Is a well-formed .feedback-data__score-table.
 * @returns true or false table is full and valid
 */
function validate_scores_table(table) {
	let validated = true;
	for (const row of table.rows) {
		if (!row.classList.contains("feedback-data__categories") &&
		!row.classList.contains("feedback-data__colavg")) {
			let data;
			for (const cell of row.cells) {
				if (!cell.classList.contains("feedback-data__memavg") &&
					!cell.classList.contains("feedback-data__row-name") &&
					!cell.classList.contains("feedback-data__colavg")) {
					data = cell.firstChild.value;
					if (!is_valid(data)) {
						if (!cell.classList.contains("feedback-data__cell--invalid")) {
							cell.classList.add("feedback-data__cell--invalid");
						}
						validated = false;
					}
				}
			}
		}
	}
	if (!validated) {
		console.error("Invalid score table data");
		document.getElementById("successful_submit").innerText = "Form submission ERROR.";
	}
	return validated;
}



/**
 * Check there is input in all comment boxes
 * @returns true if all comments are filled in
 */
function validate_comments() {
	//document.querySelector(".feedback-data__score-table")
	let validated = true;
	let data;
	for (const cell of document.querySelectorAll(".feedback-comments__member-comments")) {
		data = cell.value;
		if (data=="" || typeof data=="undefined") {
			validated = false;
			cell.parentElement.classList.add("feedback-comments__member-comments--invalid");
		}
	}
	if (!validated) {
		console.error("Incomplete comments");
		document.getElementById("successful_submit").innerText = "Form submission ERROR.";
	}
	return validated;
}

/**
 * Collect scores from a .feedback-data__score-table.
 * @param {HTMLTableElement} table
 * @requires table Is a well-formed .feedback-data__score-table.
 * @returns A list of member names and score lists.
 */
function collect_scores(table) {
	const data = [];
	for (const row of table.rows) {
		if (!row.classList.contains("feedback-data__categories") &&
		!row.classList.contains("feedback-data__colavg")) {
			let name;
			const datarow = [];
			for (const cell of row.cells) {
				if (cell.classList.contains("feedback-data__row-name")) {
					name = cell.innerText;
				} else if (!cell.classList.contains("feedback-data__row-avg")) {
					datarow.push(parseFloat(cell.firstChild.value));
				}
			}
			data.push({name, scores: datarow});
		}
	}
	return data;
}

/**
 * Submit form and print status beneath button or reject.
 * @listens MouseEvent
 */
function submit_form() {
	let check_comments = validate_comments();
	let check_scores = validate_scores_table(document.querySelector(".feedback-data__score-table"));
	if (!check_scores || !check_comments) {
		return;
	}
	// Collect feedback score table data.
	const entries = collect_scores(document.querySelector(".feedback-data__score-table"));
	/** @type {string} */
	const improvement = document.querySelector("#self_improvement").value;

	const comments = [];
	document.querySelectorAll(".feedback-comments__member-comments").forEach(e => {
		if (e.id !== "self_improvement") {
			comments.push(e.value);
		}
	});
	if (entries.length !== comments.length) {
		console.error("Score table and comment rows do not match.");
		/** @todo Make this a fatal error once pages are generated. */
	}
	for (let i = 0; i < entries.length; ++i) {
		entries[i].comment = comments[i];
	}

	// Create Submission object.
	const membersubmission = {
		author: entries[0].name,
		entries,
		improvement
	};

	// Submit and report status.
	fetch("/api/submit", {
		method: "POST",
		body: JSON.stringify(membersubmission),
		headers: {
			"Content-type": "application/json; charset=UTF-8"
		}
	}).then(res => {
		if (res.ok) {
			document.getElementById("successful_submit").innerText = "Form submitted successfully.";
		} else {
			document.getElementById("successful_submit").innerText = "Form submission error.";
		}
	}).catch(err => {
		document.getElementById("successful_submit").innerText = "Form submission error.";
	});
}
