/**
 * @file Handles submissions from index.html.
 * @author Aidan Hoover, Aiden Woodruff
 * @copyright Aidan Hoover and Aiden Woodruff 2023
 * @license BSD-3-Clause
 */

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
	// Collect feedback score table data.
	const submissions = collect_scores(document.querySelector(".feedback-data__score-table"));
	const improvement = document.querySelector("#self_improvement").value;

	const comments = [];
	document.querySelectorAll(".feedback-comments__member-comments").forEach(e => {
		if (e.id !== "self_improvement") {
			comments.push(e.value);
		}
	});
	if (submissions.length !== comments.length) {
		console.error("Score table and comment rows do not match.");
		/** @todo Make this a fatal error once pages are generated. */
	}
	for (let i = 0; i < submissions.length; ++i) {
		submissions[i].comment = comments[i];
	}

	// Create Submission object.
	const membersubmission = {
		author: scores[0].name,
		submissions,
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
