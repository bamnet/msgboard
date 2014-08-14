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
		$scope.page = Page.get({pageId: pageId});

		$scope.update = function(page) {
			$scope.page = angular.copy(page);
			$scope.page.$update({pageId: pageId}, function(){
				$location.path('/pages/' + pageId);
			}, function(err){
				$scope.error = err.data;
			});
		};
	}
]);

msgboardControllers.controller('PageCreateCtrl', ['$scope', '$location', 'Page',
	function($scope, $location, Page) {
		$scope.page =  new Page();
		$scope.create = function(page) {
			$scope.page = angular.copy(page);
			$scope.page.$create(function(){
				$location.path('/pages/');
			}, function(err){
				$scope.error = err.data;
			});
		};
	}
]);