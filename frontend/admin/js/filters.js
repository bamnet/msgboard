angular.module('msgboardFilters', []).filter('capitalize', function () {
	return function (input) {
		if (!input || input.length <= 1) {
			return input
		}
		return input.charAt(0).toUpperCase() + input.slice(1);
	};
});