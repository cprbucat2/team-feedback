/**
 * @file Handle UI elements on useradmin.html.
 * @author Aiden Woodruff
 * @copyright 2023 Aidan Hoover and Aiden Woodruff
 * @license BSD-3-Clause
 */

jQuery(function ($) {
	/**
	 * Listen to checkbox clicks and update remove button and head checkbox.
	 * @listens MouseEvent
	 */
	$(".user-list__checkbox").on("click", () => {
		const boxes = $(".user-list__checkbox");
		const count = boxes.map((i, e) => e.checked).get().reduce((x, sum) => x + sum);
		if (count === 0) {
			$(".remove-users-btn").prop("disabled", true);
			$(".user-list__head-checkbox").prop("checked", false);
			$(".user-list__head-checkbox").prop("indeterminate", false);
		} else {
			$(".remove-users-btn").prop("disabled", false);
			$(".user-list__head-checkbox").prop("checked", count === boxes.length);
			$(".user-list__head-checkbox").prop("indeterminate", count !== boxes.length);
		}
		$(".remove-users-btn").prop("value", `Remove (${count} selected)`);
	});

	/**
	 * Update checkboxes based off the header checkbox.
	 * @listens MouseEvent
	 */
	$(".user-list__head-checkbox").on("click", event => {
		if (!event.target.checked) {
			$(".user-list__checkbox").prop("checked", false);
			$(".remove-users-btn").prop("disabled", true);
			$(".remove-users-btn").prop("value", `Remove (0 selected)`);
		} else {
			const boxes = $(".user-list__checkbox");
			boxes.prop("checked", true);
			$(".remove-users-btn").prop("disabled", false);
			$(".remove-users-btn").prop("value", `Remove (${boxes.length} selected)`);
		}
	});
});
