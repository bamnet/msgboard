var msgboardApp = angular.module('msgboardApp', [
	'ngRoute',
	'ngSanitize',
	'msgboardControllers',
	'msgboardServices',
	'msgboardFilters',
	'msgboardDirectives'
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
			when('/blurbs', {
				templateUrl: 'partials/blurbs/show.html',
				controller: 'BlurbsShowCtrl'
			}).
			when('/blurbs/add', {
				templateUrl: 'partials/blurbs/add.html',
				controller: 'BlurbsAddCtrl'
			}).
			when('/blurbs/edit', {
				templateUrl: 'partials/blurbs/edit.html',
				controller: 'BlurbsEditCtrl'
			}).
			otherwise({
				redirectTo: '/pages'
			});
	}
]);
