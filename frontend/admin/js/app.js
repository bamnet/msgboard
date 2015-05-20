var msgboardApp = angular.module('msgboardApp', [
	'ngRoute',
	'ngSanitize',
	'msgboardControllers',
	'msgboardServices',
	'msgboardFilters'
]);

msgboardApp.config(['$routeProvider',
	function($routeProvider) {
		$routeProvider.
			when('/pages', {
				templateUrl: 'partials/pages/list.html',
				controller: 'PageListCtrl'
			}).
			when('/pages/new', {
				templateUrl: 'partials/pages/new.html',
				controller: 'PageCreateCtrl'
			}).
			when('/pages/:pageId', {
				templateUrl: 'partials/pages/show.html',
				controller: 'PageShowCtrl'
			}).
			when('/pages/:pageId/edit', {
				templateUrl: 'partials/pages/edit.html',
				controller: 'PageEditCtrl'
			}).
			otherwise({
				redirectTo: '/pages'
			});
	}
]);