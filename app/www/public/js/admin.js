/**
 * @file Handle UI elements on admin pages.
 * @author Aiden Woodruff
 * @copyright 2023 Aidan Hoover and Aiden Woodruff
 * @license BSD-3-Clause
 */

jQuery(function ($) {
	/**
	 * Listen to checkbox clicks and update remove button and head checkbox.
	 * @listens MouseEvent
	 */
	$(".admin-list__checkbox").on("click", () => {
		const boxes = $(".admin-list__checkbox");
		const count = boxes.map((i, e) => e.checked).get().reduce((x, sum) => x + sum);
		if (count === 0) {
			$(".admin-remove-btn").prop("disabled", true);
			$(".admin-list__head-checkbox").prop("checked", false);
			$(".admin-list__head-checkbox").prop("indeterminate", false);
		} else {
			$(".admin-remove-btn").prop("disabled", false);
			$(".admin-list__head-checkbox").prop("checked", count === boxes.length);
			$(".admin-list__head-checkbox").prop("indeterminate", count !== boxes.length);
		}
		$(".admin-remove-btn").prop("value", `Remove (${count} selected)`);
	});

	/**
	 * Update checkboxes based off the header checkbox.
	 * @listens MouseEvent
	 */
	$(".admin-list__head-checkbox").on("click", event => {
		if (!event.target.checked) {
			$(".admin-list__checkbox").prop("checked", false);
			$(".admin-remove-btn").prop("disabled", true);
			$(".admin-remove-btn").prop("value", `Remove (0 selected)`);
		} else {
			const boxes = $(".admin-list__checkbox");
			boxes.prop("checked", true);
			$(".admin-remove-btn").prop("disabled", false);
			$(".admin-remove-btn").prop("value", `Remove (${boxes.length} selected)`);
		}
	});

	/**
	 * Control the add button.
	 * @listens MouseEvent
	 */
	$("team-add-btn").on("click", event => {
		if (event.target.disabled) {
			event.stopPropagation();
		} else {
			event.target.disabled = true;
			$(".team-list__add-info").removeClass("team-list__add-info--hidden");
		}
	});
});
