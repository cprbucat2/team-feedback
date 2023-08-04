/**
 * @file Handles updating averages for .feedback-data__score-table
 * @author Aiden Woodruff
 * @copyright 2023 Aidan Hoover and Aiden Woodruff
 * @license BSD-3-Clause
 */

(function () {
	/**
	 * Update member averages.
	 * @param {HTMLTableElement} table A table.feedback-data__score-table
	 */
	function update_member_averages(table) {
		for (const row of table.rows) {
			if (row.classList.contains("feedback-data__categories") ||
			row.classList.contains("feedback-data__colavg-row")) {
				continue;
			}
			let sum = 0, count = 0;
			for (const cell of row.cells) {
				if (cell.classList.contains("feedback-data__cell")) {
					const val = parseFloat(cell.querySelector(".feedback-data__input").value);
					if (val) sum += val;
					++count;
				} else if (cell.classList.contains("feedback-data__memavg")) {
					cell.innerText = (sum / count).toFixed(2);
				}
			}
		}
	}

	/**
	 * Update column averages.
	 * @param {HTMLTableElement} table A table.feedback-data__score-table
	 */
	function update_column_averages(table) {
		const sums = [];
		for (const row of table.rows) {
			if (row.classList.contains("feedback-data__categories") ||
			row.classList.contains("feedback-data__colavg-row")) {
				continue;
			}
			for (let i = 0; i < row.cells.length; ++i) {
				if (row.cells[i].classList.contains("feedback-data__cell")) {
					sums[i] = sums[i] ? sums[i] : 0;
					if (row.cells[i].querySelector(".feedback-data__input").value) {
						sums[i] += parseFloat(row.cells[i].querySelector(".feedback-data__input").value);
					}
				} else if (row.cells[i].classList.contains("feedback-data__memavg")) {
					sums[i] = sums[i] ? sums[i] : 0;
					if (row.cells[i].innerText) {
						sums[i] += parseFloat(row.cells[i].innerText);
					}
				}
			}
		}
		const team_size = table.rows.length - 2;
		const avg_row = table.rows[table.rows.length - 1];
		for (let i = 1; i < avg_row.cells.length; ++i) {
			if (!avg_row.cells[i].classList.contains("feedback-data__row-name")) {
				avg_row.cells[i].innerText = (sums[i] / team_size).toFixed(2);
			}
		}
	}

	/**
	 * Update average row.
	 * @param {KeyboardEvent} event
	 * @listens KeyboardEvent
	 */
	function update_averages(event) {
		/** The parent table. @type {HTMLTableElement} */
		const table = event.target.parentElement.parentElement.parentElement.parentElement;
		update_member_averages(table);
		update_column_averages(table);
	}

	window.addEventListener("load", () => {
		document.querySelectorAll(".feedback-data__score-table .feedback-data__input").forEach(el => {
			el.addEventListener("input", update_averages);
		});

		// Initial averages (0s or saved form values)
		document.querySelectorAll(".feedback-data__score-table").forEach(el => {
			update_member_averages(el);
			update_column_averages(el);
		})
	});
})();
