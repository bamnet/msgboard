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

msgboardControllers.controller('PageEditCtrl', ['$scope', '$routeParams', '$location', 'Page',
  function($scope, $routeParams, $location, Page) {
    var pageId = $routeParams.pageId;
    $scope.page =  Page.get({pageId: pageId});

    $scope.update = function(page) {
      $scope.page = angular.copy(page);
      $scope.page.$update({pageId: pageId});
      $location.path('/pages/' + pageId);
    };
  }
]);