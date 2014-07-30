var msgboardApp = angular.module('msgboardApp', [
  'ngRoute',
  'msgboardControllers',
	'msgboardServices'
]);

msgboardApp.config(['$routeProvider',
  function($routeProvider) {
    $routeProvider.
      when('/pages', {
        templateUrl: 'partials/list.html',
        controller: 'PageListCtrl'
      }).
      when('/pages/:pageId', {
        templateUrl: 'partials/show.html',
        controller: 'PageShowCtrl'
      }).
      when('/pages/:pageId/edit', {
        templateUrl: 'partials/edit.html',
        controller: 'PageEditCtrl'
      }).
      otherwise({
        redirectTo: '/pages'
      });
  }
]);