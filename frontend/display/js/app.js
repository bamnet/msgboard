var msgboardApp = angular.module('msgboardApp', [
	'ngRoute',
	'ngSanitize',
	'msgboardControllers',
	'msgboardServices',
]);

msgboardApp.config(['$routeProvider',
	function($routeProvider) {
		$routeProvider.
			when('/', {
				templateUrl: 'partials/displays/show.html',
				controller: 'DisplayShowCtrl'
			}).
			otherwise({
				redirectTo: '/'
			});
	}
]);
