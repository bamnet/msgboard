var msgboardControllers = angular.module('msgboardControllers', []);

msgboardControllers.controller('PageListCtrl', ['$scope', 'Page',
  function ($scope, Page) {
  	$scope.pages = Page.list();
  }
]);

msgboardControllers.controller('PageShowCtrl', ['$scope', '$routeParams', 'Page',
	function($scope, $routeParams, Page) {
		$scope.page =  Page.get({pageId: $routeParams.pageId});
  }
]);